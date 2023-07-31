package oracle_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	ojoapp "github.com/comdex-official/comdex/app"
	appparams "github.com/comdex-official/comdex/app/params"
	"github.com/comdex-official/comdex/types/decmath"
	"github.com/comdex-official/comdex/x/oracle"
	"github.com/comdex-official/comdex/x/oracle/types"
)

const (
	displayDenom string = appparams.DisplayDenom
	bondDenom    string = appparams.BondDenom
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *ojoapp.App
}

const (
	initialPower = int64(1000)
)

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	isCheckTx := false
	app := ojoapp.Setup(s.T())
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
	})

	oracle.InitGenesis(ctx, app.OracleKeeper, *types.DefaultGenesisState())

	// validate setup... app.Setup creates one validator, with 1uumee self delegation
	setupVals := app.StakingKeeper.GetBondedValidatorsByPower(ctx)
	s.Require().Len(setupVals, 1)
	s.Require().Equal(int64(1), setupVals[0].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))

	sh := teststaking.NewHelper(s.T(), ctx, *app.StakingKeeper)
	sh.Denom = bondDenom

	// mint and send coins to validators
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins.MulInt(sdk.NewIntFromUint64(3))))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr3, initCoins))

	// mint and send coins to oracle module to fill up reward pool
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, initCoins))

	// self delegate 999 in total ... 1 val with 1uumee is already created in app.Setup
	sh.CreateValidatorWithValPower(valAddr1, valPubKey1, 599, true)
	sh.CreateValidatorWithValPower(valAddr2, valPubKey2, 398, true)
	sh.CreateValidatorWithValPower(valAddr3, valPubKey3, 2, true)

	staking.EndBlocker(ctx, *app.StakingKeeper)

	app.OracleKeeper.SetVoteThreshold(ctx, sdk.MustNewDecFromStr("0.4"))

	s.app = app
	s.ctx = ctx
}

// Test addresses
var (
	valPubKeys = simapp.CreateTestPubKeys(3)

	valPubKey1 = valPubKeys[0]
	pubKey1    = secp256k1.GenPrivKey().PubKey()
	addr1      = sdk.AccAddress(pubKey1.Address())
	valAddr1   = sdk.ValAddress(pubKey1.Address())

	valPubKey2 = valPubKeys[1]
	pubKey2    = secp256k1.GenPrivKey().PubKey()
	addr2      = sdk.AccAddress(pubKey2.Address())
	valAddr2   = sdk.ValAddress(pubKey2.Address())

	valPubKey3 = valPubKeys[2]
	pubKey3    = secp256k1.GenPrivKey().PubKey()
	addr3      = sdk.AccAddress(pubKey3.Address())
	valAddr3   = sdk.ValAddress(pubKey3.Address())

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(bondDenom, initTokens))
)

func createVotes(hash string, val sdk.ValAddress, rates sdk.DecCoins, blockHeight uint64) (types.AggregateExchangeRatePrevote, types.AggregateExchangeRateVote) {
	preVote := types.AggregateExchangeRatePrevote{
		Hash:        hash,
		Voter:       val.String(),
		SubmitBlock: uint64(blockHeight),
	}
	vote := types.AggregateExchangeRateVote{
		ExchangeRates: rates,
		Voter:         val.String(),
	}
	return preVote, vote
}

func (s *IntegrationTestSuite) TestEndBlockerVoteThreshold() {
	app, ctx := s.app, s.ctx
	ctx = ctx.WithBlockHeight(0)
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val3DecCoins sdk.DecCoins
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1DecCoins = append(val1DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("1.0"),
		})
		val2DecCoins = append(val2DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.5"),
		})
		val3DecCoins = append(val3DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
	}

	// add junk coin and ensure ballot still is counted
	junkCoin := sdk.DecCoin{
		Denom:  "JUNK",
		Amount: sdk.MustNewDecFromStr("0.05"),
	}
	val1DecCoins = append(val1DecCoins, junkCoin)
	val2DecCoins = append(val2DecCoins, junkCoin)
	val3DecCoins = append(val3DecCoins, junkCoin)

	h := uint64(ctx.BlockHeight())
	val1PreVotes, val1Votes := createVotes("hash1", valAddr1, val1DecCoins, h)
	val2PreVotes, val2Votes := createVotes("hash2", valAddr2, val2DecCoins, h)
	val3PreVotes, val3Votes := createVotes("hash3", valAddr3, val3DecCoins, h)

	// total voting power per denom is 100%
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	err := oracle.EndBlocker(ctx, app.OracleKeeper)
	s.Require().NoError(err)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("1.0"), rate)
	}

	// Test: only val2 votes (has 39% vote power).
	// Total voting power per denom must be bigger or equal than 40% (see SetupTest).
	// So if only val2 votes, we won't have any prices next block.
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = h

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().ErrorIs(err, types.ErrUnknownDenom.Wrap(denom.SymbolDenom))
		s.Require().Equal(sdk.ZeroDec(), rate)
	}

	// Test: val2 and val3 votes.
	// now we will have 40% of the power, so now we should have prices
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = h
	val3PreVotes.SubmitBlock = h

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		rate, err := app.OracleKeeper.GetExchangeRate(ctx, denom.SymbolDenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.MustNewDecFromStr("0.5"), rate)
	}

	// Test: val1 and val2 vote again
	// umee has 69.9% power, and atom has 30%, so we should have price for umee, but not for atom
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	h = uint64(ctx.BlockHeight())
	val1PreVotes.SubmitBlock = h
	val2PreVotes.SubmitBlock = h

	val1Votes.ExchangeRates = sdk.DecCoins{
		sdk.NewDecCoinFromDec("ojo", sdk.MustNewDecFromStr("1.0")),
	}
	val2Votes.ExchangeRates = sdk.DecCoins{
		sdk.NewDecCoinFromDec("atom", sdk.MustNewDecFromStr("0.5")),
	}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	rate, err := app.OracleKeeper.GetExchangeRate(ctx, "ojo")
	s.Require().NoError(err)
	s.Require().Equal(sdk.MustNewDecFromStr("1.0"), rate)
	rate, err = app.OracleKeeper.GetExchangeRate(ctx, "atom")
	s.Require().ErrorIs(err, types.ErrUnknownDenom.Wrap("atom"))
	s.Require().Equal(sdk.ZeroDec(), rate)
}

func (s *IntegrationTestSuite) TestEndBlockerValidatorRewards() {
	app, ctx := s.app, s.ctx
	originalBlockHeight := ctx.BlockHeight()
	preVoteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx) / 2)
	voteBlockDiff := int64(app.OracleKeeper.VotePeriod(ctx)/2 + 1)

	// start test in new slash window
	ctx = ctx.WithBlockHeight(int64(app.OracleKeeper.SlashWindow(ctx)))
	oracle.EndBlocker(ctx, app.OracleKeeper)

	app.OracleKeeper.SetMandatoryList(ctx, types.DenomList{
		{
			BaseDenom:   bondDenom,
			SymbolDenom: appparams.DisplayDenom,
			Exponent:    uint32(6),
		},
		{
			BaseDenom:   "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9",
			SymbolDenom: "atom",
			Exponent:    uint32(6),
		},
	})

	var (
		val1DecCoins sdk.DecCoins
		val2DecCoins sdk.DecCoins
		val3DecCoins sdk.DecCoins
	)
	for _, denom := range app.OracleKeeper.AcceptList(ctx) {
		val1DecCoins = append(val1DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
		val2DecCoins = append(val2DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
		val3DecCoins = append(val3DecCoins, sdk.DecCoin{
			Denom:  denom.SymbolDenom,
			Amount: sdk.MustNewDecFromStr("0.6"),
		})
	}

	h := uint64(ctx.BlockHeight())
	val1PreVotes, val1Votes := createVotes("hash1", valAddr1, val1DecCoins, h)
	val2PreVotes, val2Votes := createVotes("hash2", valAddr2, val2DecCoins, h)
	val3PreVotes, val3Votes := createVotes("hash3", valAddr3, val3DecCoins, h)
	// validator 1, 2, and 3 vote on both currencies so all have 0 misses
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 142), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val3PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1 and 3 votes on both currencies to end up with 0 misses
	// validator 2 votes on 1 currency to end up with 1 misses
	val1DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val2DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val3DecCoins = sdk.DecCoins{
		sdk.DecCoin{
			Denom:  "ojo",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
		sdk.DecCoin{
			Denom:  "atom",
			Amount: sdk.MustNewDecFromStr("0.6"),
		},
	}
	val1Votes.ExchangeRates = val1DecCoins
	val2Votes.ExchangeRates = val2DecCoins
	val3Votes.ExchangeRates = val3DecCoins

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 284), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 275), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 284), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])

	// update prevotes' block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + preVoteBlockDiff)
	val1PreVotes.SubmitBlock = uint64(ctx.BlockHeight())
	val2PreVotes.SubmitBlock = uint64(ctx.BlockHeight())

	// validator 1, 2, and 3 miss both currencies so validator 1 and 3 has 2 misses and
	// validator 2 has 3 misses
	val1Votes.ExchangeRates = sdk.DecCoins{}
	val2Votes.ExchangeRates = sdk.DecCoins{}
	val3Votes.ExchangeRates = sdk.DecCoins{}

	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr1, val1PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr2, val2PreVotes)
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr3, val3PreVotes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + voteBlockDiff)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, val1Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr2, val2Votes)
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr3, val3Votes)
	oracle.EndBlocker(ctx, app.OracleKeeper)

	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 426), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr1).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 408), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr2).Rewards[0])
	s.Require().Equal(sdk.NewInt64DecCoin("uojo", 426), app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr3).Rewards[0])

	ctx = ctx.WithBlockHeight(originalBlockHeight)
}

var exchangeRates = map[string][]sdk.Dec{
	"ATOM": {
		sdk.MustNewDecFromStr("12.99"),
		sdk.MustNewDecFromStr("12.22"),
		sdk.MustNewDecFromStr("13.1"),
		sdk.MustNewDecFromStr("11.6"),
	},
	"OJO": {
		sdk.MustNewDecFromStr("1.89"),
		sdk.MustNewDecFromStr("2.05"),
		sdk.MustNewDecFromStr("2.34"),
		sdk.MustNewDecFromStr("1.71"),
	},
}

func (s *IntegrationTestSuite) TestEndblockerHistoracle() {
	app, ctx := s.app, s.ctx
	blockHeight := ctx.BlockHeight()

	var historicStampPeriod int64 = 3
	var medianStampPeriod int64 = 12
	var maximumPriceStamps int64 = 4
	var maximumMedianStamps int64 = 3

	app.OracleKeeper.SetHistoricStampPeriod(ctx, uint64(historicStampPeriod))
	app.OracleKeeper.SetMedianStampPeriod(ctx, uint64(medianStampPeriod))
	app.OracleKeeper.SetMaximumPriceStamps(ctx, uint64(maximumPriceStamps))
	app.OracleKeeper.SetMaximumMedianStamps(ctx, uint64(maximumMedianStamps))

	// Start at the last block of the first stamp period
	blockHeight += medianStampPeriod
	blockHeight += -1
	ctx = ctx.WithBlockHeight(blockHeight)

	for i := int64(0); i <= maximumMedianStamps; i++ {
		for j := int64(0); j < maximumPriceStamps; j++ {

			blockHeight += historicStampPeriod
			ctx = ctx.WithBlockHeight(blockHeight)

			decCoins := sdk.DecCoins{}
			for denom, prices := range exchangeRates {
				decCoins = append(decCoins, sdk.DecCoin{
					Denom:  denom,
					Amount: prices[j],
				})
			}

			vote := types.AggregateExchangeRateVote{
				ExchangeRates: decCoins,
				Voter:         valAddr1.String(),
			}
			app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr1, vote)
			oracle.EndBlocker(ctx, app.OracleKeeper)
		}

		for denom, denomRates := range exchangeRates {
			// check median
			expectedMedian, err := decmath.Median(denomRates)
			s.Require().NoError(err)

			medians := app.OracleKeeper.AllMedianPrices(ctx)
			medians = *medians.FilterByBlock(uint64(blockHeight)).FilterByDenom(denom)
			actualMedian := medians[0].ExchangeRate.Amount
			s.Require().Equal(expectedMedian, actualMedian)

			// check median deviation
			expectedMedianDeviation, err := decmath.MedianDeviation(actualMedian, denomRates)
			s.Require().NoError(err)

			medianDeviations := app.OracleKeeper.AllMedianDeviationPrices(ctx)
			medianDeviations = *medianDeviations.FilterByBlock(uint64(blockHeight)).FilterByDenom(denom)
			actualMedianDeviation := medianDeviations[0].ExchangeRate.Amount
			s.Require().Equal(expectedMedianDeviation, actualMedianDeviation)
		}
	}
	numberOfAssets := int64(len(exchangeRates))

	historicPrices := app.OracleKeeper.AllHistoricPrices(ctx)
	s.Require().Equal(maximumPriceStamps*numberOfAssets, int64(len(historicPrices)))

	for i := int64(0); i < maximumPriceStamps; i++ {
		expectedBlockNum := blockHeight - (historicStampPeriod * (maximumPriceStamps - int64(i+1)))
		actualBlockNum := historicPrices[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}

	medians := app.OracleKeeper.AllMedianPrices(ctx)
	s.Require().Equal(maximumMedianStamps*numberOfAssets, int64(len(medians)))

	for i := int64(0); i < maximumMedianStamps; i++ {
		expectedBlockNum := blockHeight - (medianStampPeriod * (maximumMedianStamps - int64(i+1)))
		actualBlockNum := medians[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}

	medianDeviations := app.OracleKeeper.AllMedianPrices(ctx)
	s.Require().Equal(maximumMedianStamps*numberOfAssets, int64(len(medianDeviations)))

	for i := int64(0); i < maximumMedianStamps; i++ {
		expectedBlockNum := blockHeight - (medianStampPeriod * (maximumMedianStamps - int64(i+1)))
		actualBlockNum := medianDeviations[i].BlockNum
		s.Require().Equal(expectedBlockNum, int64(actualBlockNum))
	}
}

func TestOracleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
