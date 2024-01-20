package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestCreatePair() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	testCases := []struct {
		Name               string
		Msg                types.MsgCreatePair
		ExpErr             error
		ExpResp            *types.Pair
		QueryResponseIndex uint64
		QueryResponse      *types.Pair
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error app id invalid",
			Msg: *types.NewMsgCreatePair(
				69, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr:             errorsmod.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error both non whitelisted denoms",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, "dummy1", "dummy2",
			),
			ExpErr:             errorsmod.Wrapf(types.ErrAssetNotWhiteListed, "asset with denom  %s is not white listed", "dummy1"),
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error first non whitelisted denom",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, "dummy1", asset2.Denom,
			),
			ExpErr:             errorsmod.Wrapf(types.ErrAssetNotWhiteListed, "asset with denom  %s is not white listed", "dummy1"),
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error second non whitelisted denom",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset1.Denom, "dummy2",
			),
			ExpErr:             errorsmod.Wrapf(types.ErrAssetNotWhiteListed, "asset with denom  %s is not white listed", "dummy2"),
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error insufficient pair creation fee",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr:             errorsmod.Wrap(errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "spendable balance  is smaller than 2000000000ucmdx"), "insufficient pair creation fee"),
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case app1 pair1",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr: nil,
			ExpResp: &types.Pair{
				Id:                      1,
				BaseCoinDenom:           asset1.Denom,
				QuoteCoinDenom:          asset2.Denom,
				EscrowAddress:           "cosmos1url34vfv5a5a7esm2aapklqelh2mzuwe34vvgc94eh752t8mcjeqla0n0v",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos19a7w3ferywxjst035636dzktx94xyh22u64pwee3fl62sennw5hsw8erx3",
				AppId:                   appID1,
			},
			QueryResponseIndex: 0,
			QueryResponse: &types.Pair{
				Id:                      1,
				BaseCoinDenom:           asset1.Denom,
				QuoteCoinDenom:          asset2.Denom,
				EscrowAddress:           "cosmos1url34vfv5a5a7esm2aapklqelh2mzuwe34vvgc94eh752t8mcjeqla0n0v",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos19a7w3ferywxjst035636dzktx94xyh22u64pwee3fl62sennw5hsw8erx3",
				AppId:                   appID1,
			},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "success valid case app1 pair2",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset2.Denom, asset1.Denom,
			),
			ExpErr: nil,
			ExpResp: &types.Pair{
				Id:                      2,
				BaseCoinDenom:           asset2.Denom,
				QuoteCoinDenom:          asset1.Denom,
				EscrowAddress:           "cosmos174fjyc3w25m7pku6sww9w03e9s8phphe4h88euxypsp5ekjngdkqp596l4",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos1sfhtyclc7mvz356z578yc44jqr2pe3tpmcswjngqc3mvh7r2kmcs9pd44e",
				AppId:                   appID1,
			},
			QueryResponseIndex: 1,
			QueryResponse: &types.Pair{
				Id:                      2,
				BaseCoinDenom:           asset2.Denom,
				QuoteCoinDenom:          asset1.Denom,
				EscrowAddress:           "cosmos174fjyc3w25m7pku6sww9w03e9s8phphe4h88euxypsp5ekjngdkqp596l4",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos1sfhtyclc7mvz356z578yc44jqr2pe3tpmcswjngqc3mvh7r2kmcs9pd44e",
				AppId:                   appID1,
			},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error pair already exists 1",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr:             types.ErrPairAlreadyExists,
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error pair already exists 2",
			Msg: *types.NewMsgCreatePair(
				appID1, addr1, asset2.Denom, asset1.Denom,
			),
			ExpErr:             types.ErrPairAlreadyExists,
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case app2 pair1",
			Msg: *types.NewMsgCreatePair(
				appID2, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr: nil,
			ExpResp: &types.Pair{
				Id:                      1,
				BaseCoinDenom:           asset1.Denom,
				QuoteCoinDenom:          asset2.Denom,
				EscrowAddress:           "cosmos1mjjgv53lgef35ywsfd4hwduyyptqreev2k9ngh2p84plhrc69srqdmdydt",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos1v9h2ymqyhu34py8shdl59j62gw3h4qykxns76s0epedyq890ca9qvvjw87",
				AppId:                   appID2,
			},
			QueryResponseIndex: 0,
			QueryResponse: &types.Pair{
				Id:                      1,
				BaseCoinDenom:           asset1.Denom,
				QuoteCoinDenom:          asset2.Denom,
				EscrowAddress:           "cosmos1mjjgv53lgef35ywsfd4hwduyyptqreev2k9ngh2p84plhrc69srqdmdydt",
				LastOrderId:             0,
				LastPrice:               nil,
				CurrentBatchId:          1,
				SwapFeeCollectorAddress: "cosmos1v9h2ymqyhu34py8shdl59j62gw3h4qykxns76s0epedyq890ca9qvvjw87",
				AppId:                   appID2,
			},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error pair already exists 3",
			Msg: *types.NewMsgCreatePair(
				appID2, addr1, asset1.Denom, asset2.Denom,
			),
			ExpErr:             types.ErrPairAlreadyExists,
			ExpResp:            &types.Pair{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), params.PairCreationFee)
			}

			resp, err := s.keeper.CreatePair(s.ctx, &tc.Msg, false)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, &resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Creator))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)
				collectedPairCreationFee := s.getBalances(sdk.MustAccAddressFromBech32(params.FeeCollectorAddress))
				s.Require().True(sdk.NewCoins(sdk.NewCoin(params.PairCreationFee[0].Denom, params.PairCreationFee[0].Amount.Mul(sdkmath.NewInt(int64(tc.QueryResponseIndex+1))))).Equal(collectedPairCreationFee))

				pairs := s.keeper.GetAllPairs(s.ctx, tc.Msg.AppId)
				s.Require().Len(pairs, int(tc.QueryResponseIndex)+1)
				s.Require().Equal(tc.QueryResponse.Id, pairs[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.BaseCoinDenom, pairs[tc.QueryResponseIndex].BaseCoinDenom)
				s.Require().Equal(tc.QueryResponse.QuoteCoinDenom, pairs[tc.QueryResponseIndex].QuoteCoinDenom)
				s.Require().Equal(tc.QueryResponse.EscrowAddress, pairs[tc.QueryResponseIndex].EscrowAddress)
				s.Require().Equal(tc.QueryResponse.LastOrderId, pairs[tc.QueryResponseIndex].LastOrderId)
				s.Require().Equal(tc.QueryResponse.LastPrice, pairs[tc.QueryResponseIndex].LastPrice)
				s.Require().Equal(tc.QueryResponse.CurrentBatchId, pairs[tc.QueryResponseIndex].CurrentBatchId)
				s.Require().Equal(tc.QueryResponse.SwapFeeCollectorAddress, pairs[tc.QueryResponseIndex].SwapFeeCollectorAddress)
				s.Require().Equal(tc.QueryResponse.AppId, pairs[tc.QueryResponseIndex].AppId)

			}
		})
	}
}
