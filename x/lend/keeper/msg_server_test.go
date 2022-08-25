package keeper_test

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
)

func (s *KeeperTestSuite) TestMsgLend() {

	assetOneID := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	assetTwoID := s.CreateNewAsset("ASSET2", "uasset2", 2000000)
	assetThreeID := s.CreateNewAsset("ASSET3", "uasset3", 2000000)
	assetFourID := s.CreateNewAsset("ASSET4", "uasset4", 2000000)
	cAssetOneID := s.CreateNewAsset("CASSET1", "ucasset1", 1000000)
	cAssetTwoID := s.CreateNewAsset("CASSET2", "ucasset2", 2000000)
	cAssetThreeID := s.CreateNewAsset("CASSET3", "ucasset3", 2000000)
	//cAssetFourID := s.CreateNewAsset("CASSET4", "ucasset4", 2000000)

	var (
		assetDataPoolOne []types.AssetDataPoolMapping
		assetDataPoolTwo []types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := types.AssetDataPoolMapping{
		AssetID:   assetOneID,
		IsBridged: false,
	}
	assetDataPoolOneAssetTwo := types.AssetDataPoolMapping{
		AssetID:   assetTwoID,
		IsBridged: true,
	}
	assetDataPoolOneAssetThree := types.AssetDataPoolMapping{
		AssetID:   assetThreeID,
		IsBridged: true,
	}
	assetDataPoolTwoAssetFour := types.AssetDataPoolMapping{
		AssetID:   assetFourID,
		IsBridged: true,
	}

	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolOne, assetDataPoolTwoAssetFour, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)

	poolOneID := s.CreateNewPool("cmdx", "CMDX-ATOM-CMST", assetOneID, assetTwoID, assetThreeID, assetDataPoolOne)
	poolTwoID := s.CreateNewPool("osmo", "OSMO-ATOM-CMST", assetFourID, assetTwoID, assetThreeID, assetDataPoolTwo)

	s.AddAssetRatesStats(assetThreeID, newDec("0.8"), newDec("0.002"), newDec("0.06"), newDec("0.6"), true, newDec("0.04"), newDec("0.04"), newDec("0.06"), newDec("0.8"), newDec("0.85"), newDec("0.025"), newDec("0.025"), newDec("0.1"), cAssetThreeID)
	s.AddAssetRatesStats(assetOneID, newDec("0.75"), newDec("0.002"), newDec("0.07"), newDec("1.25"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.7"), newDec("0.75"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetOneID)
	//s.AddAssetRatesStats(assetFourID, newDec("0.65"), newDec("0.002"), newDec("0.08"), newDec("1.5"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.6"), newDec("0.65"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetFourID)
	s.AddAssetRatesStats(assetTwoID, newDec("0.5"), newDec("0.002"), newDec("0.08"), newDec("2.0"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.5"), newDec("0.55"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetTwoID)

	pairOneID := s.AddExtendedLendPair(assetTwoID, assetThreeID, false, poolOneID, 1000000)
	pairTwoID := s.AddExtendedLendPair(assetTwoID, assetOneID, false, poolOneID, 1000000)
	pairThreeID := s.AddExtendedLendPair(assetOneID, assetTwoID, false, poolOneID, 1000000)
	pairFourID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolOneID, 1000000)
	pairFiveID := s.AddExtendedLendPair(assetThreeID, assetTwoID, false, poolOneID, 1000000)
	pairSixID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolOneID, 1000000)
	pairSevenID := s.AddExtendedLendPair(assetFourID, assetThreeID, false, poolTwoID, 1000000)
	pairEightID := s.AddExtendedLendPair(assetFourID, assetOneID, false, poolTwoID, 1000000)
	pairNineID := s.AddExtendedLendPair(assetOneID, assetFourID, false, poolTwoID, 1000000)
	pairTenID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolTwoID, 1000000)
	pairElevenID := s.AddExtendedLendPair(assetThreeID, assetFourID, false, poolTwoID, 1000000)
	pairTwelveID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolTwoID, 1000000)
	pairThirteenID := s.AddExtendedLendPair(assetTwoID, assetFourID, true, poolTwoID, 1000000)
	pairFourteenID := s.AddExtendedLendPair(assetThreeID, assetFourID, true, poolTwoID, 1000000)
	pairFifteenID := s.AddExtendedLendPair(assetOneID, assetFourID, true, poolTwoID, 1000000)
	pairSixteenID := s.AddExtendedLendPair(assetFourID, assetTwoID, true, poolOneID, 1000000)
	pairSeventeenID := s.AddExtendedLendPair(assetThreeID, assetTwoID, true, poolOneID, 1000000)
	pairEighteenID := s.AddExtendedLendPair(assetOneID, assetTwoID, true, poolOneID, 1000000)

	s.AddAssetToPair(assetOneID, poolOneID, []uint64{pairThreeID, pairFourID, pairFifteenID})
	s.AddAssetToPair(assetTwoID, poolOneID, []uint64{pairOneID, pairTwoID, pairThirteenID})
	s.AddAssetToPair(assetThreeID, poolOneID, []uint64{pairFiveID, pairSixID, pairFourteenID})
	s.AddAssetToPair(assetFourID, poolTwoID, []uint64{pairSevenID, pairEightID, pairSixteenID})
	s.AddAssetToPair(assetOneID, poolTwoID, []uint64{pairNineID, pairTenID, pairEighteenID})
	s.AddAssetToPair(assetThreeID, poolTwoID, []uint64{pairElevenID, pairTwelveID, pairSeventeenID})

	appOneID := s.CreateNewApp("commodo", "cmmdo")
	appTwoID := s.CreateNewApp("cswap", "cswap")

	testCases := []struct {
		Name               string
		Msg                types.MsgLend
		ExpErr             error
		ExpResp            *types.MsgLendResponse
		QueryResponseIndex uint64
		QueryResponse      *types.LendAsset
		AvailableBalance   sdk.Coins
	}{
		{
			Name:               "asset does not exist",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", 10, sdk.NewCoin("uasset1", sdk.NewInt(100)), poolOneID, 1),
			ExpErr:             types.ErrorAssetDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(0))),
		},
		{
			Name:               "Pool Not Found",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", sdk.NewInt(100)), 3, 1),
			ExpErr:             types.ErrPoolNotFound,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(0))),
		},
		{
			Name:               "App Mapping Id does not exists",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", sdk.NewInt(100)), poolOneID, 10),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(0))),
		},
		{
			Name:               "App Mapping Id mismatch, use the correct App Mapping ID in request",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", sdk.NewInt(100)), poolOneID, appTwoID),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))),
		},
		{
			Name:               "invalid offer coin amount",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset2", sdk.NewInt(100)), poolOneID, appOneID),
			ExpErr:             sdkerrors.Wrapf(types.ErrBadOfferCoinAmount, "uasset2"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))),
		},
		{
			Name:               "Asset Id not defined in the pool",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetFourID, sdk.NewCoin("uasset4", sdk.NewInt(100)), poolOneID, appOneID),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidAssetIDForPool, "4"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))),
		},
		{
			Name:               "Asset Rates Stats not found",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetFourID, sdk.NewCoin("uasset4", sdk.NewInt(100)), poolTwoID, appOneID),
			ExpErr:             sdkerrors.Wrapf(types.ErrorAssetRatesStatsNotFound, "4"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))),
		},
		{
			Name:               "success valid case",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", sdk.NewInt(100)), poolOneID, appOneID),
			ExpErr:             nil,
			ExpResp:            &types.MsgLendResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.LendAsset{
				ID:                 1,
				AssetID:            assetOneID,
				PoolID:             poolOneID,
				Owner:              "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:           sdk.NewCoin("uasset1", sdk.NewInt(100)),
				LendingTime:        time.Time{},
				UpdatedAmountIn:    sdk.NewInt(100),
				AvailableToBorrow:  sdk.NewInt(100),
				Reward_Accumulated: sdk.NewInt(0),
				AppID:              appOneID,
				CPoolName:          "OSMO-ATOM-CMST",
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("ucasset1", newInt(100))),
		},
		{
			Name:               "Duplicate lend Position",
			Msg:                *types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", sdk.NewInt(100)), poolOneID, appOneID),
			ExpErr:             types.ErrorDuplicateLend,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))),
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Lender), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.Lend(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Lender))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))
			}
		})
	}

}

func (s *KeeperTestSuite) TestMsgWithdraw() {

	assetOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	assetTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	assetThreeID := s.CreateNewAsset("ASSETHREE", "uasset3", 2000000)
	assetFourID := s.CreateNewAsset("ASSETFOUR", "uasset4", 2000000)
	cAssetOneID := s.CreateNewAsset("CASSET1", "ucasset1", 1000000)
	cAssetTwoID := s.CreateNewAsset("CASSET2", "ucasset2", 2000000)
	cAssetThreeID := s.CreateNewAsset("CASSET3", "ucasset3", 2000000)
	//cAssetFourID := s.CreateNewAsset("CASSET4", "ucasset4", 2000000)

	var (
		assetDataPoolOne []types.AssetDataPoolMapping
		assetDataPoolTwo []types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := types.AssetDataPoolMapping{
		AssetID:   assetOneID,
		IsBridged: false,
	}
	assetDataPoolOneAssetTwo := types.AssetDataPoolMapping{
		AssetID:   assetTwoID,
		IsBridged: true,
	}
	assetDataPoolOneAssetThree := types.AssetDataPoolMapping{
		AssetID:   assetThreeID,
		IsBridged: true,
	}
	assetDataPoolTwoAssetFour := types.AssetDataPoolMapping{
		AssetID:   assetFourID,
		IsBridged: true,
	}

	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolOne, assetDataPoolTwoAssetFour, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)

	poolOneID := s.CreateNewPool("cmdx", "CMDX-ATOM-CMST", assetOneID, assetTwoID, assetThreeID, assetDataPoolOne)
	poolTwoID := s.CreateNewPool("osmo", "OSMO-ATOM-CMST", assetFourID, assetTwoID, assetThreeID, assetDataPoolTwo)

	s.AddAssetRatesStats(assetThreeID, newDec("0.8"), newDec("0.002"), newDec("0.06"), newDec("0.6"), true, newDec("0.04"), newDec("0.04"), newDec("0.06"), newDec("0.8"), newDec("0.85"), newDec("0.025"), newDec("0.025"), newDec("0.1"), cAssetThreeID)
	s.AddAssetRatesStats(assetOneID, newDec("0.75"), newDec("0.002"), newDec("0.07"), newDec("1.25"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.7"), newDec("0.75"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetOneID)
	//s.AddAssetRatesStats(assetFourID, newDec("0.65"), newDec("0.002"), newDec("0.08"), newDec("1.5"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.6"), newDec("0.65"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetFourID)
	s.AddAssetRatesStats(assetTwoID, newDec("0.5"), newDec("0.002"), newDec("0.08"), newDec("2.0"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.5"), newDec("0.55"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetTwoID)

	pairOneID := s.AddExtendedLendPair(assetTwoID, assetThreeID, false, poolOneID, 1000000)
	pairTwoID := s.AddExtendedLendPair(assetTwoID, assetOneID, false, poolOneID, 1000000)
	pairThreeID := s.AddExtendedLendPair(assetOneID, assetTwoID, false, poolOneID, 1000000)
	pairFourID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolOneID, 1000000)
	pairFiveID := s.AddExtendedLendPair(assetThreeID, assetTwoID, false, poolOneID, 1000000)
	pairSixID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolOneID, 1000000)
	pairSevenID := s.AddExtendedLendPair(assetFourID, assetThreeID, false, poolTwoID, 1000000)
	pairEightID := s.AddExtendedLendPair(assetFourID, assetOneID, false, poolTwoID, 1000000)
	pairNineID := s.AddExtendedLendPair(assetOneID, assetFourID, false, poolTwoID, 1000000)
	pairTenID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolTwoID, 1000000)
	pairElevenID := s.AddExtendedLendPair(assetThreeID, assetFourID, false, poolTwoID, 1000000)
	pairTwelveID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolTwoID, 1000000)
	pairThirteenID := s.AddExtendedLendPair(assetTwoID, assetFourID, true, poolTwoID, 1000000)
	pairFourteenID := s.AddExtendedLendPair(assetThreeID, assetFourID, true, poolTwoID, 1000000)
	pairFifteenID := s.AddExtendedLendPair(assetOneID, assetFourID, true, poolTwoID, 1000000)
	pairSixteenID := s.AddExtendedLendPair(assetFourID, assetTwoID, true, poolOneID, 1000000)
	pairSeventeenID := s.AddExtendedLendPair(assetThreeID, assetTwoID, true, poolOneID, 1000000)
	pairEighteenID := s.AddExtendedLendPair(assetOneID, assetTwoID, true, poolOneID, 1000000)

	s.AddAssetToPair(assetOneID, poolOneID, []uint64{pairThreeID, pairFourID, pairFifteenID})
	s.AddAssetToPair(assetTwoID, poolOneID, []uint64{pairOneID, pairTwoID, pairThirteenID})
	s.AddAssetToPair(assetThreeID, poolOneID, []uint64{pairFiveID, pairSixID, pairFourteenID})
	s.AddAssetToPair(assetFourID, poolTwoID, []uint64{pairSevenID, pairEightID, pairSixteenID})
	s.AddAssetToPair(assetOneID, poolTwoID, []uint64{pairNineID, pairTenID, pairEighteenID})
	s.AddAssetToPair(assetThreeID, poolTwoID, []uint64{pairElevenID, pairTwelveID, pairSeventeenID})

	appOneID := s.CreateNewApp("commodo", "cmmdo")

	msg := types.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", newInt(100)), poolOneID, appOneID)
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100))))
	s.msgServer.Lend(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgWithdraw
		ExpErr             error
		ExpResp            *types.MsgWithdrawResponse
		QueryResponseIndex uint64
		QueryResponse      *types.LendAsset
		AvailableBalance   sdk.Coins
	}{
		{
			Name:               "Withdraw Amount Limit Exceeded",
			Msg:                *types.NewMsgWithdraw("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", 1, sdk.NewCoin("uasset1", sdk.NewInt(100))),
			ExpErr:             types.ErrWithdrawAmountLimitExceeds,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(0))),
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Lender), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.Withdraw(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Lender))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))
			}
		})
	}

}
