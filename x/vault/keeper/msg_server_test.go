package keeper_test

import (
	"fmt"
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestMsgCreate() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	pairID3 := s.CreateNewPair(addr1, asseTwoID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)
	extendedVaultPairID4 := s.CreateNewExtendedVaultPair("CMDX-E", appID1, pairID3, true, false)

	testCases := []struct {
		Name               string
		Msg                types.MsgCreateRequest
		ExpErr             error
		ExpResp            *types.MsgCreateResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, 123, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgCreateRequest(
				addr1, 12, extendedVaultPairID1, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID2, extendedVaultPairID1, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgCreateRequest(
				[]byte(""), appID1, extendedVaultPairID1, newInt(10000000), newInt(4000000),
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault creation inactive",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID3, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorVaultCreationInactive,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "stable mint ext pair",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID4, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorCannotCreateStableMintVault,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "AmountOut Greater Than DebtCeiling",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newIntStr("100000000000000000000"), newInt(1000000000000000001),
			),
			ExpErr:             types.ErrorAmountOutGreaterThanDebtCeiling,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error amount out smaller that debt floor",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newInt(10000000), newInt(4000000),
			),
			ExpErr:             types.ErrorAmountOutLessThanDebtFloor,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid collateralization ratio",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newInt(800000000), newInt(200000000),
			),
			ExpErr:             types.ErrorInvalidCollateralizationRatio,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("spendable balance  is smaller than %duasset1: insufficient funds", 1000000000)),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCreateResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
		{
			Name: "success valid case app1 user2",
			Msg: *types.NewMsgCreateRequest(
				addr2, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCreateResponse{},
			QueryResponseIndex: 1,
			QueryResponse: &types.Vault{
				Id:                  2,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr2.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
		{
			Name: "error user vault already exists",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             types.ErrorUserVaultAlreadyExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app2 user1",
			Msg: *types.NewMsgCreateRequest(
				addr1, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCreateResponse{},
			QueryResponseIndex: 2,
			QueryResponse: &types.Vault{
				Id:                  3,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000*2))),
		},
		{
			Name: "success valid case app2 user2",
			Msg: *types.NewMsgCreateRequest(
				addr2, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCreateResponse{},
			QueryResponseIndex: 3,
			QueryResponse: &types.Vault{
				Id:                  4,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr2.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000*2))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.AmountIn)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgCreate(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDeposit() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr2, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr2.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgDepositRequest
		ExpErr             error
		ExpResp            *types.MsgDepositResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgDepositRequest(
				[]byte(""), appID1, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID1, 123, 1, newInt(4000000),
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgDepositRequest(
				addr1, 69, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault creation inactive",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID1, extendedVaultPairID3, 1, newInt(4000000),
			),
			ExpErr:             types.ErrorVaultInactive,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid AppMapping Data",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("App Mapping Id mismatch, use the correct App Mapping ID in request")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID2, extendedVaultPairID2, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("Invalid App Mapping data sent as compared to data exists in vault")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error vault does not exists",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID1, extendedVaultPairID1, 4, newInt(69000000),
			),
			ExpErr:             types.ErrorVaultDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error access unathorized ",
			Msg: *types.NewMsgDepositRequest(
				addr2, appID1, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             types.ErrVaultAccessUnauthorised,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("spendable balance  is smaller than %duasset1: insufficient funds", 69000000)),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgDepositRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgDepositResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1069000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
		{
			Name: "success valid case app2 user2",
			Msg: *types.NewMsgDepositRequest(
				addr2, appID2, extendedVaultPairID2, 2, newInt(69000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgDepositResponse{},
			QueryResponseIndex: 1,
			QueryResponse: &types.Vault{
				Id:                  2,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr2.String(),
				AmountIn:            newInt(1069000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgDeposit(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdraw() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr1, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgWithdrawRequest
		ExpErr             error
		ExpResp            *types.MsgWithdrawResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgWithdrawRequest(
				[]byte(""), appID1, extendedVaultPairID1, 1, newInt(400000000),
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID1, 123, 1, newInt(400000000),
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, 69, extendedVaultPairID1, 1, newInt(400000000),
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(400000000),
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error vault does not exists",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID1, extendedVaultPairID1, 3, newInt(400000000),
			),
			ExpErr:             types.ErrorVaultDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error access unathorized",
			Msg: *types.NewMsgWithdrawRequest(
				addr2, appID1, extendedVaultPairID1, 1, newInt(400000000),
			),
			ExpErr:             types.ErrVaultAccessUnauthorised,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault creation inactive",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID1, extendedVaultPairID3, 1, newInt(4000000),
			),
			ExpErr:             types.ErrorVaultInactive,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid AppMapping Data",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("App Mapping Id mismatch, use the correct App Mapping ID in request")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID2, extendedVaultPairID2, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("Invalid App Mapping data sent as compared to data exists in vault")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid collateralization ratio",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(400000000),
			),
			ExpErr:             types.ErrorInvalidCollateralizationRatio,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgWithdrawResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(950000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(50000000)), sdk.NewCoin("uasset2", newInt(198000000*2))),
		},
		{
			Name: "success valid case app2 user1",
			Msg: *types.NewMsgWithdrawRequest(
				addr1, appID2, extendedVaultPairID2, 2, newInt(50000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgWithdrawResponse{},
			QueryResponseIndex: 1,
			QueryResponse: &types.Vault{
				Id:                  2,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr1.String(),
				AmountIn:            newInt(950000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(100000000)), sdk.NewCoin("uasset2", newInt(198000000*2))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgWithdraw(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDraw() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr1, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgDrawRequest
		ExpErr             error
		ExpResp            *types.MsgDrawResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgDrawRequest(
				[]byte(""), appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID1, 123, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgDrawRequest(
				addr1, 69, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault creation inactive",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID1, extendedVaultPairID3, 1, newInt(4000000),
			),
			ExpErr:             types.ErrorVaultInactive,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid AppMapping Data",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("App Mapping Id mismatch, use the correct App Mapping ID in request")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID2, extendedVaultPairID2, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("Invalid App Mapping data sent as compared to data exists in vault")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error vault does not exists",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID1, extendedVaultPairID1, 3, newInt(50000000),
			),
			ExpErr:             types.ErrorVaultDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error access unathorized",
			Msg: *types.NewMsgDrawRequest(
				addr2, appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrVaultAccessUnauthorised,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid collateralization ratio",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorInvalidCollateralizationRatio,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(10000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgDrawResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(210000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(405900000))), // 198000000*2+(10000000-1%)
		},
		{
			Name: "success valid case app2 user1",
			Msg: *types.NewMsgDrawRequest(
				addr1, appID2, extendedVaultPairID2, 2, newInt(10000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgDrawResponse{},
			QueryResponseIndex: 1,
			QueryResponse: &types.Vault{
				Id:                  2,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(210000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(415800000))), // 198000000*2+(10000000-1%) + (10000000-1%)
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgDraw(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgRepay() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr1, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgRepayRequest
		ExpErr             error
		ExpResp            *types.MsgRepayResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgRepayRequest(
				[]byte(""), appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID1, 123, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgRepayRequest(
				addr1, 69, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error vault does not exists",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID1, extendedVaultPairID1, 3, newInt(50000000),
			),
			ExpErr:             types.ErrorVaultDoesNotExist,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error access unathorized",
			Msg: *types.NewMsgRepayRequest(
				addr2, appID1, extendedVaultPairID1, 1, newInt(50000000),
			),
			ExpErr:             types.ErrVaultAccessUnauthorised,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID1, extendedVaultPairID3, 1, newInt(4000000),
			),
			ExpErr:             types.ErrorInvalidExtendedPairMappingData,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid AppMapping Data",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID2, extendedVaultPairID1, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("App Mapping Id mismatch, use the correct App Mapping ID in request")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID2, extendedVaultPairID2, 1, newInt(4000000),
			),
			ExpErr:             fmt.Errorf(fmt.Sprintf("Invalid App Mapping data sent as compared to data exists in vault")),
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid amount",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(0),
			),
			ExpErr:             types.ErrorInvalidAmount,
			ExpResp:            nil,
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(100000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgRepayResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(100000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(296000000))), // 198000000*2 - 100000000
		},
		{
			Name: "success valid case app2 user1",
			Msg: *types.NewMsgRepayRequest(
				addr1, appID2, extendedVaultPairID2, 2, newInt(100000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgRepayResponse{},
			QueryResponseIndex: 1,
			QueryResponse: &types.Vault{
				Id:                  2,
				AppId:               appID2,
				ExtendedPairVaultID: extendedVaultPairID2,
				Owner:               addr1.String(),
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(100000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(196000000))), // 198000000*2 - 100000000 - 100000000
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.Name, func() {
			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgRepay(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgClose() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr1, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	// add more asset to close the vaults, since 2% is being reduced from asseOut being provided.i.e 2% + 2% assetout of above vaults.
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset2", newInt(4000000))))

	testCases := []struct {
		Name               string
		Msg                types.MsgCloseRequest
		ExpErr             error
		ExpResp            *types.MsgCloseResponse
		AvailableVaultsLen uint64
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgLiquidateRequest(
				[]byte(""), appID1, extendedVaultPairID1, 1,
			),
			ExpErr:             fmt.Errorf("empty address string is not allowed"),
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID1, 123, 1,
			),
			ExpErr:             types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, 69, extendedVaultPairID1, 1,
			),
			ExpErr:             types.ErrorAppMappingDoesNotExist,
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid AppMapping Data",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID2, extendedVaultPairID1, 1,
			),
			ExpErr:           fmt.Errorf(fmt.Sprintf("App Mapping Id mismatch, use the correct App Mapping ID in request")),
			ExpResp:          nil,
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "vault Invalid Extended Pair Mapping Data",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID2, extendedVaultPairID2, 1,
			),
			ExpErr:           fmt.Errorf(fmt.Sprintf("Invalid App Mapping data sent as compared to data exists in vault")),
			ExpResp:          nil,
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID2, extendedVaultPairID1, 1,
			),
			ExpErr:             types.ErrorAppMappingIDMismatch,
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error vault does not exists",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID1, extendedVaultPairID1, 3,
			),
			ExpErr:             types.ErrorVaultDoesNotExist,
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "error access unathorized",
			Msg: *types.NewMsgLiquidateRequest(
				addr2, appID1, extendedVaultPairID1, 1,
			),
			ExpErr:             types.ErrVaultAccessUnauthorised,
			ExpResp:            nil,
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset2", newInt(0))),
		},
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID1, extendedVaultPairID1, 1,
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCloseResponse{},
			AvailableVaultsLen: 1,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000)), sdk.NewCoin("uasset2", newInt(200000000))),
		},
		{
			Name: "success valid case app2 user1",
			Msg: *types.NewMsgLiquidateRequest(
				addr1, appID2, extendedVaultPairID2, 2,
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgCloseResponse{},
			AvailableVaultsLen: 0,
			AvailableBalance:   sdk.NewCoins(sdk.NewCoin("uasset1", newInt(2000000000))),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.Name, func() {
			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgClose(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, int(tc.AvailableVaultsLen))
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDepositAndDraw() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	//asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	//pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	//extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)

	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr2, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr2.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	testCases := []struct {
		Name               string
		Msg                types.MsgDepositAndDrawRequest
		ExpErr             error
		ExpResp            *types.MsgDepositAndDrawResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "success valid case app1 user1",
			Msg: *types.NewMsgDepositAndDrawRequest(
				addr1, appID1, extendedVaultPairID1, 1, newInt(69000000),
			),
			ExpErr:             nil,
			ExpResp:            &types.MsgDepositAndDrawResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1069000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgDepositAndDraw(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				//availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				//s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				//s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateStableMint() {
	addr1 := s.addr(1)
	// addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	appID3 := s.CreateNewApp("appthr")
	// appID4 := s.CreateNewApp("appfou")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, false)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, true, false)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-C", appID3, pairID, true, true)
	// extendedVaultPairID4 := s.CreateNewExtendedVaultPair("CMDX-C", appID4, pairID, true, true)

	testCases := []struct {
		Name             string
		Msg              types.MsgCreateStableMintRequest
		ExpErr           error
		ExpResp          *types.MsgCreateStableMintResponse
		QueryRespIndex   uint64
		QueryResponse    *types.StableMintVault
		AvailableBalance sdk.Coins
	}{
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID1, 123, newInt(10000),
			),
			ExpErr:           types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, 69, extendedVaultPairID1, newInt(10000),
			),
			ExpErr:           types.ErrorAppMappingDoesNotExist,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID2, extendedVaultPairID1, newInt(10000),
			),
			ExpErr:           types.ErrorAppMappingIDMismatch,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgCreateStableMintRequest(
				[]byte(""), appID1, extendedVaultPairID1, newInt(10000),
			),
			ExpErr:           fmt.Errorf("empty address string is not allowed"),
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error non psm pair cannot create stable mint vault",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID1, extendedVaultPairID1, newInt(10000),
			),
			ExpErr:           types.ErrorCannotCreateStableMintVault,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error vault creation inactive",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID2, extendedVaultPairID2, newInt(10000),
			),
			ExpErr:           types.ErrorVaultCreationInactive,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error vault creation inactive",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID2, extendedVaultPairID2, newInt(10000),
			),
			ExpErr:           types.ErrorVaultCreationInactive,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID3, extendedVaultPairID3, newInt(10000),
			),
			ExpErr:           types.ErrorAmountOutLessThanDebtFloor,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "success valid case app3 user1",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID3, extendedVaultPairID3, newInt(1000000000),
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgCreateStableMintResponse{},
			QueryRespIndex: 0,
			QueryResponse: &types.StableMintVault{
				Id:                  1,
				AmountIn:            newInt(1000000000),
				AmountOut:           newInt(1000000000),
				AppId:               3,
				ExtendedPairVaultID: 3,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(990000000))),
		},
		{
			Name: "error stable mint vault already exists",
			Msg: *types.NewMsgCreateStableMintRequest(
				addr1, appID3, extendedVaultPairID3, newInt(10000),
			),
			ExpErr:           types.ErrorAmountOutLessThanDebtFloor, // stable mint vault got registered in above testcase, dk how app state got changed even though error was thrown.
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		// {
		// 	Name: "success valid case app4 user1",
		// 	Msg: *types.NewMsgCreateStableMintRequest(
		// 		addr1, appID4, extendedVaultPairID4, newInt(10000),
		// 	),
		// 	ExpErr:         nil,
		// 	ExpResp:        &types.MsgCreateStableMintResponse{},
		// 	QueryRespIndex: 0,
		// 	QueryResponse: &types.StableMintVault{
		// 		Id:                  1,
		// 		AmountIn:            newInt(10000),
		// 		AmountOut:           newInt(10000),
		// 		AppId:               4,
		// 		ExtendedPairVaultID: 4,
		// 	},
		// 	AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(19800))),
		// },
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.Name, func() {
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgCreateStableMint(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
				s.Require().Equal(tc.QueryResponse.Id, stableMintVaults[tc.QueryRespIndex].Id)
				s.Require().Equal(tc.QueryResponse.AmountIn, stableMintVaults[tc.QueryRespIndex].AmountIn)
				fmt.Println("stableMintVaults[tc.QueryRespIndex].AmountOut", stableMintVaults[tc.QueryRespIndex].AmountOut)
				fmt.Println("tc.QueryResponse.AmountOut", tc.QueryResponse.AmountOut)
				s.Require().Equal(tc.QueryResponse.AmountOut, stableMintVaults[tc.QueryRespIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.AppId, stableMintVaults[tc.QueryRespIndex].AppId)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, stableMintVaults[tc.QueryRespIndex].ExtendedPairVaultID)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDepositStableMint() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, false)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-A", appID1, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-B", appID1, pairID, true, true)
	extendedVaultPairID4 := s.CreateNewExtendedVaultPair("CMDX-D", appID2, pairID, true, true)

	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	_, err := s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
		addr1, appID1, extendedVaultPairID3, newInt(1000000000),
	))
	s.Require().NoError(err)

	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	_, err = s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
		addr1, appID2, extendedVaultPairID4, newInt(1000000000),
	))
	s.Require().NoError(err)

	testCases := []struct {
		Name             string
		Msg              types.MsgDepositStableMintRequest
		ExpErr           error
		ExpResp          *types.MsgDepositStableMintResponse
		QueryRespIndex   uint64
		QueryResponse    *types.StableMintVault
		AvailableBalance sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgDepositStableMintRequest(
				[]byte(""), appID1, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           fmt.Errorf("empty address string is not allowed"),
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, 123, newInt(10000), 1,
			),
			ExpErr:           types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, 69, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAppMappingDoesNotExist,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error vault creation inactive",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           types.ErrorVaultInactive,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error non stable mint vault cannot create stable mint vault",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID2, newInt(10000), 1,
			),
			ExpErr:           types.ErrorCannotCreateStableMintVault,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID2, extendedVaultPairID3, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAppMappingIDMismatch,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid stable mint id",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(10000), 3,
			),
			ExpErr:           types.ErrorVaultDoesNotExist,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAmountOutLessThanDebtFloor,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "success valid case 1 app1 user1",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(2000000000), 1,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgDepositStableMintResponse{},
			QueryRespIndex: 0,
			QueryResponse: &types.StableMintVault{
				Id:                  1,
				AmountIn:            newInt(3000000000),
				AmountOut:           newInt(3000000000),
				AppId:               1,
				ExtendedPairVaultID: 3,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(990000000*4))),
		},
		{
			Name: "success valid 2 case app1 user1",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(1000000000), 1,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgDepositStableMintResponse{},
			QueryRespIndex: 0,
			QueryResponse: &types.StableMintVault{
				Id:                  1,
				AmountIn:            newInt(4000000000),
				AmountOut:           newInt(4000000000),
				AppId:               1,
				ExtendedPairVaultID: 3,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(990000000*5))),
		},
		{
			Name: "success valid 3 case app2 user1",
			Msg: *types.NewMsgDepositStableMintRequest(
				addr1, appID2, extendedVaultPairID4, newInt(9000000000), 2,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgDepositStableMintResponse{},
			QueryRespIndex: 1,
			QueryResponse: &types.StableMintVault{
				Id:                  2,
				AmountIn:            newInt(10000000000),
				AmountOut:           newInt(10000000000),
				AppId:               2,
				ExtendedPairVaultID: 4,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(990000000*14))),
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.Name, func() {
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgDepositStableMint(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
				s.Require().Equal(tc.QueryResponse.Id, stableMintVaults[tc.QueryRespIndex].Id)
				s.Require().Equal(tc.QueryResponse.AmountIn, stableMintVaults[tc.QueryRespIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, stableMintVaults[tc.QueryRespIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.AppId, stableMintVaults[tc.QueryRespIndex].AppId)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, stableMintVaults[tc.QueryRespIndex].ExtendedPairVaultID)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgWithdrawStableMint() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)
	addr3 := s.addr(3)
	addr4 := s.addr(4)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, false)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-A", appID1, pairID, false, true)
	extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-B", appID1, pairID, true, true)
	extendedVaultPairID4 := s.CreateNewExtendedVaultPair("CMDX-D", appID2, pairID, true, true)

	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	_, err := s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
		addr1, appID1, extendedVaultPairID3, newInt(1000000000),
	))
	s.Require().NoError(err)

	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	_, err = s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
		addr1, appID2, extendedVaultPairID4, newInt(1000000000),
	))
	s.Require().NoError(err)

	s.fundAddr(addr3, sdk.NewCoins(sdk.NewCoin("uasset2", newInt(5000000))))
	s.fundAddr(addr4, sdk.NewCoins(sdk.NewCoin("uasset2", newInt(5050000))))

	testCases := []struct {
		Name             string
		Msg              types.MsgWithdrawStableMintRequest
		ExpErr           error
		ExpResp          *types.MsgWithdrawStableMintResponse
		QueryRespIndex   uint64
		QueryResponse    *types.StableMintVault
		AvailableBalance sdk.Coins
	}{
		{
			Name: "error invalid from address",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				[]byte(""), appID1, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           fmt.Errorf("empty address string is not allowed"),
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error extended vault pair does not exists",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, 123, newInt(10000), 1,
			),
			ExpErr:           types.ErrorExtendedPairVaultDoesNotExists,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid appID",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, 69, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAppMappingDoesNotExist,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error vault creation inactive",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID1, newInt(10000), 1,
			),
			ExpErr:           types.ErrorCannotCreateStableMintVault,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error non stable mint vault cannot create stable mint vault",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID2, newInt(10000), 1,
			),
			ExpErr:           types.ErrorCannotCreateStableMintVault,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error appID mismatch",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID2, extendedVaultPairID3, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAppMappingIDMismatch,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid stable mint id",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(10000), 3,
			),
			ExpErr:           types.ErrorAmountOutLessThanDebtFloor,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr2, appID1, extendedVaultPairID3, newInt(10000), 1,
			),
			ExpErr:           types.ErrorAmountOutLessThanDebtFloor,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid withdraw amount",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(9000000000), 1,
			),
			ExpErr:           types.ErrorInvalidAmount,
			ExpResp:          nil,
			QueryRespIndex:   0,
			QueryResponse:    nil,
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "success valid case 1 app1 user1",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(500000000), 1,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgWithdrawStableMintResponse{},
			QueryRespIndex: 0,
			QueryResponse: &types.StableMintVault{
				Id:                  1,
				AmountIn:            newInt(505000000),
				AmountOut:           newInt(505000000),
				AppId:               1,
				ExtendedPairVaultID: 3,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(495000000)), sdk.NewCoin("uasset2", newInt((990000000*2)-500000000))),
		},
		{
			Name: "success valid case 2 case app1 user1",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID1, extendedVaultPairID3, newInt(200000000), 1,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgWithdrawStableMintResponse{},
			QueryRespIndex: 0,
			QueryResponse: &types.StableMintVault{
				Id:                  1,
				AmountIn:            newInt(307000000),
				AmountOut:           newInt(307000000),
				AppId:               1,
				ExtendedPairVaultID: 3,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(693000000)), sdk.NewCoin("uasset2", newInt((990000000*2)-500000000-200000000))),
		},
		{
			Name: "success valid case 3 case app2 user1",
			Msg: *types.NewMsgWithdrawStableMintRequest(
				addr1, appID2, extendedVaultPairID4, newInt(1000000000), 2,
			),
			ExpErr:         nil,
			ExpResp:        &types.MsgWithdrawStableMintResponse{},
			QueryRespIndex: 1,
			QueryResponse: &types.StableMintVault{
				Id:                  2,
				AmountIn:            newInt(10000000),
				AmountOut:           newInt(10000000),
				AppId:               2,
				ExtendedPairVaultID: 4,
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(693000000+990000000)), sdk.NewCoin("uasset2", newInt((990000000*2)-500000000-200000000-1000000000))),
		},
		//{
		//	Name: "success valid case 4 case app2 user3",
		//	Msg: *types.NewMsgWithdrawStableMintRequest(
		//		addr3, appID2, extendedVaultPairID4, newInt(5000000), 2,
		//	),
		//	ExpErr:         nil,
		//	ExpResp:        &types.MsgWithdrawStableMintResponse{},
		//	QueryRespIndex: 1,
		//	QueryResponse: &types.StableMintVault{
		//		Id:                  2,
		//		AmountIn:            newInt(5050000),
		//		AmountOut:           newInt(5050000),
		//		AppId:               2,
		//		ExtendedPairVaultID: 4,
		//	},
		//	AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(4950000))),
		//},
		// {
		// 	Name: "success valid case 5 case app2 user4",
		// 	Msg: *types.NewMsgWithdrawStableMintRequest(
		// 		addr4, appID2, extendedVaultPairID4, newInt(5050000), 3,
		// 	),
		// 	ExpErr:         nil,
		// 	ExpResp:        &types.MsgWithdrawStableMintResponse{},
		// 	QueryRespIndex: 1,
		// 	QueryResponse: &types.StableMintVault{
		// 		Id:                  3,
		// 		AmountIn:            newInt(50500),
		// 		AmountOut:           newInt(50500),
		// 		AppId:               2,
		// 		ExtendedPairVaultID: 4,
		// 	},
		// 	AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset1", newInt(4999500))),
		// },
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.Name, func() {
			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgWithdrawStableMint(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				fmt.Println("resp", resp)
				fmt.Println("tc.ExpResp", tc.ExpResp)
				s.Require().Equal(tc.ExpResp, resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				s.Require().True(tc.AvailableBalance.Equal(availableBalances))

				stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
				s.Require().Equal(tc.QueryResponse.Id, stableMintVaults[tc.QueryRespIndex].Id)
				s.Require().Equal(tc.QueryResponse.AmountIn, stableMintVaults[tc.QueryRespIndex].AmountIn)
				s.Require().Equal(tc.QueryResponse.AmountOut, stableMintVaults[tc.QueryRespIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.AppId, stableMintVaults[tc.QueryRespIndex].AppId)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, stableMintVaults[tc.QueryRespIndex].ExtendedPairVaultID)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgVaultInterestCalc() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	//asseThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	//pairID2 := s.CreateNewPair(addr1, asseOneID, asseThreeID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)
	//extendedVaultPairID3 := s.CreateNewExtendedVaultPair("CMDX-D", appID1, pairID2, false, false)
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	msg := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	msg = types.NewMsgCreateRequest(addr2, appID2, extendedVaultPairID2, newInt(1000000000), newInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr2.String()), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
	s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)

	rewardskeeper, ctx := &s.rewardsKeeper, &s.ctx
	err := rewardskeeper.WhitelistAppIDVault(*ctx, 1)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-02T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(15)

	testCases := []struct {
		Name               string
		Msg                types.MsgVaultInterestCalcRequest
		ExpErr             error
		ExpResp            *types.MsgVaultInterestCalcResponse
		QueryResponseIndex uint64
		QueryResponse      *types.Vault
		AvailableBalance   sdk.Coins
	}{
		{
			Name:               "success valid case app1 user1",
			Msg:                *types.NewMsgVaultInterestCalcRequest(addr1, appID1, 1),
			ExpErr:             nil,
			ExpResp:            &types.MsgVaultInterestCalcResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1069000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
		{
			Name:               "success valid case app1 user1",
			Msg:                *types.NewMsgVaultInterestCalcRequest(addr1, appID1, 1),
			ExpErr:             nil,
			ExpResp:            &types.MsgVaultInterestCalcResponse{},
			QueryResponseIndex: 0,
			QueryResponse: &types.Vault{
				Id:                  1,
				AppId:               appID1,
				ExtendedPairVaultID: extendedVaultPairID1,
				Owner:               addr1.String(),
				AmountIn:            newInt(1069000000),
				AmountOut:           newInt(200000000),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("uasset2", newInt(198000000))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				//s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.From), sdk.NewCoins(sdk.NewCoin("uasset1", tc.Msg.Amount)))
			}

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.msgServer.MsgVaultInterestCalc(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				//availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.From))
				//s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))

				vaults := s.keeper.GetVaults(s.ctx)
				fmt.Println("InterestAccumulated", vaults[tc.QueryResponseIndex].InterestAccumulated)
				s.Require().Len(vaults, 2)
				s.Require().Equal(tc.QueryResponse.Id, vaults[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.Owner, vaults[tc.QueryResponseIndex].Owner)
				//s.Require().Equal(tc.QueryResponse.AmountIn, vaults[tc.QueryResponseIndex].AmountIn)
				//s.Require().Equal(tc.QueryResponse.AmountOut, vaults[tc.QueryResponseIndex].AmountOut)
				s.Require().Equal(tc.QueryResponse.ExtendedPairVaultID, vaults[tc.QueryResponseIndex].ExtendedPairVaultID)
				s.Require().Equal(tc.QueryResponse.AppId, vaults[tc.QueryResponseIndex].AppId)
			}
		})
	}
}
