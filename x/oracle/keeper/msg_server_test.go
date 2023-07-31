package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"

	"github.com/comdex-official/comdex/x/oracle/types"
	oracletypes "github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GenerateSalt generates a random salt, size length/2,  as a HEX encoded string.
func GenerateSalt(length int) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("failed to generate salt: zero length")
	}

	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *IntegrationTestSuite) TestMsgServer_AggregateExchangeRatePrevote() {
	ctx := s.ctx

	exchangeRatesStr := "123.2:OJO"
	salt, err := GenerateSalt(32)
	s.Require().NoError(err)
	hash := oracletypes.GetAggregateVoteHash(salt, exchangeRatesStr, valAddr)

	invalidHash := &types.MsgAggregateExchangeRatePrevote{
		Hash:      "invalid_hash",
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}
	invalidFeeder := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    "invalid_feeder",
		Validator: valAddr.String(),
	}
	invalidValidator := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: "invalid_val",
	}
	validMsg := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}

	_, err = s.msgServer.AggregateExchangeRatePrevote(sdk.WrapSDKContext(ctx), invalidHash)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(sdk.WrapSDKContext(ctx), invalidFeeder)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(sdk.WrapSDKContext(ctx), invalidValidator)
	s.Require().Error(err)
	_, err = s.msgServer.AggregateExchangeRatePrevote(sdk.WrapSDKContext(ctx), validMsg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_AggregateExchangeRateVote() {
	ctx := s.ctx

	ratesStr := "ojo:123.2"
	ratesStrInvalidCoin := "ojo:123.2,badcoin:234.5"
	salt, err := GenerateSalt(32)
	s.Require().NoError(err)
	hash := oracletypes.GetAggregateVoteHash(salt, ratesStr, valAddr)
	hashInvalidRate := oracletypes.GetAggregateVoteHash(salt, ratesStrInvalidCoin, valAddr)

	prevoteMsg := &types.MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    addr.String(),
		Validator: valAddr.String(),
	}
	voteMsg := &types.MsgAggregateExchangeRateVote{
		Feeder:        addr.String(),
		Validator:     valAddr.String(),
		Salt:          salt,
		ExchangeRates: ratesStr,
	}
	voteMsgInvalidRate := &types.MsgAggregateExchangeRateVote{
		Feeder:        addr.String(),
		Validator:     valAddr.String(),
		Salt:          salt,
		ExchangeRates: ratesStrInvalidCoin,
	}

	// Flattened acceptList symbols to make checks easier
	acceptList := s.app.OracleKeeper.GetParams(ctx).AcceptList
	var acceptListFlat []string
	for _, v := range acceptList {
		acceptListFlat = append(acceptListFlat, v.SymbolDenom)
	}

	// No existing prevote
	_, err = s.msgServer.AggregateExchangeRateVote(sdk.WrapSDKContext(ctx), voteMsg)
	s.Require().EqualError(err, sdkerrors.Wrap(types.ErrNoAggregatePrevote, valAddr.String()).Error())
	_, err = s.msgServer.AggregateExchangeRatePrevote(sdk.WrapSDKContext(ctx), prevoteMsg)
	s.Require().NoError(err)
	// Reveal period mismatch
	_, err = s.msgServer.AggregateExchangeRateVote(sdk.WrapSDKContext(ctx), voteMsg)
	s.Require().EqualError(err, types.ErrRevealPeriodMissMatch.Error())

	// Valid
	s.app.OracleKeeper.SetAggregateExchangeRatePrevote(
		ctx,
		valAddr,
		types.NewAggregateExchangeRatePrevote(
			hash, valAddr, 7,
		))
	_, err = s.msgServer.AggregateExchangeRateVote(sdk.WrapSDKContext(ctx), voteMsg)
	s.Require().NoError(err)
	vote, err := s.app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Nil(err)
	for _, v := range vote.ExchangeRates {
		s.Require().Contains(acceptListFlat, strings.ToLower(v.Denom))
	}

	// Valid, but with an exchange rate which isn't in AcceptList
	s.app.OracleKeeper.SetAggregateExchangeRatePrevote(
		ctx,
		valAddr,
		types.NewAggregateExchangeRatePrevote(
			hashInvalidRate, valAddr, 7,
		))
	_, err = s.msgServer.AggregateExchangeRateVote(sdk.WrapSDKContext(ctx), voteMsgInvalidRate)
	s.Require().NoError(err)
	vote, err = s.app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().NoError(err)
	for _, v := range vote.ExchangeRates {
		s.Require().Contains(acceptListFlat, strings.ToLower(v.Denom))
	}
}

func (s *IntegrationTestSuite) TestMsgServer_DelegateFeedConsent() {
	app, ctx := s.app, s.ctx

	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	_, err := s.msgServer.DelegateFeedConsent(sdk.WrapSDKContext(ctx), &types.MsgDelegateFeedConsent{
		Operator: valAddr.String(),
		Delegate: feederAddr.String(),
	})
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestMsgServer_UpdateGovParams() {
	govAccAddr := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	testCases := []struct {
		name      string
		req       *types.MsgGovUpdateParams
		expectErr bool
		errMsg    string
	}{
		{
			"valid accept list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"AcceptList"},
				Changes: types.Params{
					AcceptList: types.DenomList{
						{
							BaseDenom:   oracletypes.OjoDenom,
							SymbolDenom: oracletypes.OjoSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   oracletypes.AtomDenom,
							SymbolDenom: oracletypes.AtomSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   "base",
							SymbolDenom: "symbol",
							Exponent:    6,
						},
					},
				},
			},
			false,
			"",
		},
		{
			"valid mandatory list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"MandatoryList"},
				Changes: types.Params{
					MandatoryList: types.DenomList{
						{
							BaseDenom:   oracletypes.OjoDenom,
							SymbolDenom: oracletypes.OjoSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   oracletypes.AtomDenom,
							SymbolDenom: oracletypes.AtomSymbol,
							Exponent:    6,
						},
					},
				},
			},
			false,
			"",
		},
		{
			"invalid mandatory list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"MandatoryList"},
				Changes: types.Params{
					MandatoryList: types.DenomList{
						{
							BaseDenom:   "test",
							SymbolDenom: "test",
							Exponent:    6,
						},
					},
				},
			},
			true,
			"denom in MandatoryList not present in AcceptList",
		},
		{
			"valid reward band list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"RewardBands"},
				Changes: types.Params{
					RewardBands: types.RewardBandList{
						{
							SymbolDenom: types.OjoSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
						{
							SymbolDenom: types.AtomSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
						{
							SymbolDenom: "symbol",
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
					},
				},
			},
			false,
			"",
		},
		{
			"invalid reward band list",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"RewardBands"},
				Changes: types.Params{
					RewardBands: types.RewardBandList{
						{
							SymbolDenom: types.OjoSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 0),
						},
						{
							SymbolDenom: types.AtomSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
					},
				},
			},
			true,
			"oracle parameter RewardBand must be between [0, 1]",
		},
		{
			"multiple valid params",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys: []string{
					"VotePeriod",
					"VoteThreshold",
					"RewardDistributionWindow",
					"SlashFraction",
					"SlashWindow",
					"MinValidPerWindow",
					"HistoricStampPeriod",
					"MedianStampPeriod",
					"MaximumPriceStamps",
					"MaximumMedianStamps",
				},
				Changes: types.Params{
					VotePeriod:               10,
					VoteThreshold:            sdk.NewDecWithPrec(40, 2),
					RewardDistributionWindow: types.BlocksPerWeek,
					SlashFraction:            sdk.NewDecWithPrec(2, 4),
					SlashWindow:              types.BlocksPerDay,
					MinValidPerWindow:        sdk.NewDecWithPrec(4, 2),
					HistoricStampPeriod:      10 * types.BlocksPerMinute,
					MedianStampPeriod:        5 * types.BlocksPerHour,
					MaximumPriceStamps:       40,
					MaximumMedianStamps:      30,
				},
			},
			false,
			"",
		},
		{
			"invalid vote threshold",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"VoteThreshold"},
				Changes: types.Params{
					VoteThreshold: sdk.NewDecWithPrec(10, 2),
				},
			},
			true,
			"threshold must be bigger than 0.330000000000000000 and <= 1",
		},
		{
			"invalid slash window",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"VotePeriod", "SlashWindow"},
				Changes: types.Params{
					VotePeriod:  5,
					SlashWindow: 4,
				},
			},
			true,
			"oracle parameter SlashWindow must be greater than or equal with VotePeriod",
		},
		{
			"invalid key",
			&types.MsgGovUpdateParams{
				Authority:   govAccAddr,
				Title:       "test",
				Description: "test",
				Keys:        []string{"test"},
				Changes:     types.Params{},
			},
			true,
			"test is not an existing oracle param key",
		},

		{
			"bad authority",
			&types.MsgGovUpdateParams{
				Authority:   "ojo1zypqa76je7pxsdwkfah6mu9a583sju6xzthge3",
				Title:       "test",
				Description: "test",
				Keys:        []string{"RewardBands"},
				Changes: types.Params{
					RewardBands: types.RewardBandList{
						{
							SymbolDenom: types.OjoSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
					},
				},
			},
			true,
			"invalid gov authority to perform these changes",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.req.ValidateBasic()
			if err == nil {
				_, err = s.msgServer.GovUpdateParams(s.ctx, tc.req)
			}
			if tc.expectErr {
				s.Require().ErrorContains(err, tc.errMsg)
			} else {
				s.Require().NoError(err)

				switch tc.name {
				case "valid accept list":
					acceptList := s.app.OracleKeeper.AcceptList(s.ctx)
					s.Require().Equal(acceptList, types.DenomList{
						{
							BaseDenom:   oracletypes.OjoDenom,
							SymbolDenom: oracletypes.OjoSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   oracletypes.AtomDenom,
							SymbolDenom: oracletypes.AtomSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   "base",
							SymbolDenom: "symbol",
							Exponent:    6,
						},
					}.Normalize())

				case "valid mandatory list":
					mandatoryList := s.app.OracleKeeper.MandatoryList(s.ctx)
					s.Require().Equal(mandatoryList, types.DenomList{
						{
							BaseDenom:   oracletypes.OjoDenom,
							SymbolDenom: oracletypes.OjoSymbol,
							Exponent:    6,
						},
						{
							BaseDenom:   oracletypes.AtomDenom,
							SymbolDenom: oracletypes.AtomSymbol,
							Exponent:    6,
						},
					}.Normalize())

				case "valid reward band list":
					rewardBand := s.app.OracleKeeper.RewardBands(s.ctx)
					s.Require().Equal(rewardBand, types.RewardBandList{
						{
							SymbolDenom: types.OjoSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
						{
							SymbolDenom: types.AtomSymbol,
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
						{
							SymbolDenom: "symbol",
							RewardBand:  sdk.NewDecWithPrec(2, 2),
						},
					})

				case "multiple valid params":
					votePeriod := s.app.OracleKeeper.VotePeriod(s.ctx)
					voteThreshold := s.app.OracleKeeper.VoteThreshold(s.ctx)
					rewardDistributionWindow := s.app.OracleKeeper.RewardDistributionWindow(s.ctx)
					slashFraction := s.app.OracleKeeper.SlashFraction(s.ctx)
					slashWindow := s.app.OracleKeeper.SlashWindow(s.ctx)
					minValidPerWindow := s.app.OracleKeeper.MinValidPerWindow(s.ctx)
					historicStampPeriod := s.app.OracleKeeper.HistoricStampPeriod(s.ctx)
					medianStampPeriod := s.app.OracleKeeper.MedianStampPeriod(s.ctx)
					maximumPriceStamps := s.app.OracleKeeper.MaximumPriceStamps(s.ctx)
					maximumMedianStamps := s.app.OracleKeeper.MaximumMedianStamps(s.ctx)
					s.Require().Equal(votePeriod, uint64(10))
					s.Require().Equal(voteThreshold, sdk.NewDecWithPrec(40, 2))
					s.Require().Equal(rewardDistributionWindow, types.BlocksPerWeek)
					s.Require().Equal(slashFraction, sdk.NewDecWithPrec(2, 4))
					s.Require().Equal(slashWindow, types.BlocksPerDay)
					s.Require().Equal(minValidPerWindow, sdk.NewDecWithPrec(4, 2))
					s.Require().Equal(historicStampPeriod, 10*types.BlocksPerMinute)
					s.Require().Equal(medianStampPeriod, 5*types.BlocksPerHour)
					s.Require().Equal(maximumPriceStamps, uint64(40))
					s.Require().Equal(maximumMedianStamps, uint64(30))
				}
			}
		})
	}
}
