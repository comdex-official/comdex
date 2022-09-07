package keeper_test

import (
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperTestSuite) TestQueryAllVaults() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
	}{
		{
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
	}

	for _, v := range newVaults {
		msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
		s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
		_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
		s.Require().NoError(err)
	}

	allVaults, err := s.querier.QueryAllVaults(sdk.WrapSDKContext(s.ctx), &types.QueryAllVaultsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(len(newVaults), len(allVaults.Vault))

	for i, v := range allVaults.Vault {
		s.Require().Equal(newVaults[i].Address.String(), v.Owner)
		s.Require().Equal(newVaults[i].AppID, v.AppId)
		s.Require().Equal(newVaults[i].ExtendedVaultPairID, v.ExtendedPairVaultID)
		s.Require().Equal(newVaults[i].AmountIn, v.AmountIn)
		s.Require().Equal(newVaults[i].AmountOut, v.AmountOut)
	}
}

func (s *KeeperTestSuite) TestQueryAllVaultsByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
	}{
		{
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
		{
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
		},
	}

	for _, v := range newVaults {
		msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
		s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
		_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
		s.Require().NoError(err)
	}

	allVaults, err := s.querier.QueryAllVaultsByApp(sdk.WrapSDKContext(s.ctx), &types.QueryAllVaultsByAppRequest{AppId: appID1})
	s.Require().NoError(err)
	s.Require().Equal(2, len(allVaults.Vault))

	allVaults, err = s.querier.QueryAllVaultsByApp(sdk.WrapSDKContext(s.ctx), &types.QueryAllVaultsByAppRequest{AppId: appID2})
	s.Require().NoError(err)
	s.Require().Equal(2, len(allVaults.Vault))
}

func (s *KeeperTestSuite) TestQueryVault() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultRequest
		ExpErr              error
		ExpResp             *types.QueryVaultResponse
	}{
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultRequest{Id: 1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultResponse{
				Vault: types.Vault{
					Id:                  1,
					AppId:               appID1,
					ExtendedPairVaultID: extendedVaultPairID1,
					Owner:               addr1.String(),
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(200000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultRequest{Id: 2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultResponse{
				Vault: types.Vault{
					Id:                  2,
					AppId:               appID1,
					ExtendedPairVaultID: extendedVaultPairID1,
					Owner:               addr2.String(),
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(200000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultRequest{Id: 3},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultResponse{
				Vault: types.Vault{
					Id:                  3,
					AppId:               appID2,
					ExtendedPairVaultID: extendedVaultPairID2,
					Owner:               addr1.String(),
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(200000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultRequest{Id: 4},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultResponse{
				Vault: types.Vault{
					Id:                  4,
					AppId:               appID2,
					ExtendedPairVaultID: extendedVaultPairID2,
					Owner:               addr2.String(),
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(200000000),
				},
			},
		},
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultRequest{Id: 0},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultResponse{},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vault, err := s.querier.QueryVault(sdk.WrapSDKContext(s.ctx), v.Req)

		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.Vault.Id, vault.Vault.Id)
			s.Require().Equal(v.ExpResp.Vault.AppId, vault.Vault.AppId)
			s.Require().Equal(v.ExpResp.Vault.ExtendedPairVaultID, vault.Vault.ExtendedPairVaultID)
			s.Require().Equal(v.ExpResp.Vault.Owner, vault.Vault.Owner)
			s.Require().Equal(v.ExpResp.Vault.AmountIn, vault.Vault.AmountIn)
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultInfoByVaultID() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultInfoByVaultIDRequest
		ExpErr              error
		ExpResp             *types.QueryVaultInfoByVaultIDResponse
	}{
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoByVaultIDRequest{Id: 1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoByVaultIDResponse{
				VaultsInfo: types.VaultInfo{
					Id:               1,
					ExtendedPairID:   extendedVaultPairID1,
					Owner:            addr1.String(),
					Collateral:       newInt(1000000000),
					Debt:             newInt(200000000),
					ExtendedPairName: "CMDX-C",
					AssetInDenom:     "uasset1",
					AssetOutDenom:    "uasset2",
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoByVaultIDRequest{Id: 2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoByVaultIDResponse{
				VaultsInfo: types.VaultInfo{
					Id:               2,
					ExtendedPairID:   extendedVaultPairID1,
					Owner:            addr2.String(),
					Collateral:       newInt(1000000000),
					Debt:             newInt(200000000),
					ExtendedPairName: "CMDX-C",
					AssetInDenom:     "uasset1",
					AssetOutDenom:    "uasset2",
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoByVaultIDRequest{Id: 3},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoByVaultIDResponse{
				VaultsInfo: types.VaultInfo{
					Id:               3,
					ExtendedPairID:   extendedVaultPairID2,
					Owner:            addr1.String(),
					Collateral:       newInt(1000000000),
					Debt:             newInt(200000000),
					ExtendedPairName: "CMDX-C",
					AssetInDenom:     "uasset1",
					AssetOutDenom:    "uasset2",
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoByVaultIDRequest{Id: 4},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoByVaultIDResponse{
				VaultsInfo: types.VaultInfo{
					Id:               4,
					ExtendedPairID:   extendedVaultPairID2,
					Owner:            addr2.String(),
					Collateral:       newInt(1000000000),
					Debt:             newInt(200000000),
					ExtendedPairName: "CMDX-C",
					AssetInDenom:     "uasset1",
					AssetOutDenom:    "uasset2",
				},
			},
		},
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultInfoByVaultIDRequest{Id: 0},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultInfoByVaultIDResponse{},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultInfo, err := s.querier.QueryVaultInfoByVaultID(sdk.WrapSDKContext(s.ctx), v.Req)

		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.VaultsInfo.Id, vaultInfo.VaultsInfo.Id)
			s.Require().Equal(v.ExpResp.VaultsInfo.ExtendedPairID, vaultInfo.VaultsInfo.ExtendedPairID)
			s.Require().Equal(v.ExpResp.VaultsInfo.Owner, vaultInfo.VaultsInfo.Owner)
			s.Require().Equal(v.ExpResp.VaultsInfo.Collateral, vaultInfo.VaultsInfo.Collateral)
			s.Require().Equal(v.ExpResp.VaultsInfo.Debt, vaultInfo.VaultsInfo.Debt)
			s.Require().Equal(v.ExpResp.VaultsInfo.ExtendedPairName, vaultInfo.VaultsInfo.ExtendedPairName)
			s.Require().Equal(v.ExpResp.VaultsInfo.AssetInDenom, vaultInfo.VaultsInfo.AssetInDenom)
			s.Require().Equal(v.ExpResp.VaultsInfo.AssetOutDenom, vaultInfo.VaultsInfo.AssetOutDenom)
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultInfoOfOwnerByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultInfoOfOwnerByAppRequest
		ExpErr              error
		ExpResp             *types.QueryVaultInfoOfOwnerByAppResponse
	}{
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoOfOwnerByAppRequest{AppId: appID1, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoOfOwnerByAppResponse{
				VaultsInfo: []types.VaultInfo{
					{
						Id:               1,
						ExtendedPairID:   extendedVaultPairID1,
						Owner:            addr1.String(),
						Collateral:       newInt(1000000000),
						Debt:             newInt(200000000),
						ExtendedPairName: "CMDX-C",
						AssetInDenom:     "uasset1",
						AssetOutDenom:    "uasset2",
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoOfOwnerByAppRequest{AppId: appID1, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoOfOwnerByAppResponse{
				VaultsInfo: []types.VaultInfo{
					{
						Id:               2,
						ExtendedPairID:   extendedVaultPairID1,
						Owner:            addr2.String(),
						Collateral:       newInt(1000000000),
						Debt:             newInt(200000000),
						ExtendedPairName: "CMDX-C",
						AssetInDenom:     "uasset1",
						AssetOutDenom:    "uasset2",
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoOfOwnerByAppRequest{AppId: appID2, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoOfOwnerByAppResponse{
				VaultsInfo: []types.VaultInfo{
					{
						Id:               3,
						ExtendedPairID:   extendedVaultPairID2,
						Owner:            addr1.String(),
						Collateral:       newInt(1000000000),
						Debt:             newInt(200000000),
						ExtendedPairName: "CMDX-C",
						AssetInDenom:     "uasset1",
						AssetOutDenom:    "uasset2",
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultInfoOfOwnerByAppRequest{AppId: appID2, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultInfoOfOwnerByAppResponse{
				VaultsInfo: []types.VaultInfo{
					{
						Id:               4,
						ExtendedPairID:   extendedVaultPairID2,
						Owner:            addr2.String(),
						Collateral:       newInt(1000000000),
						Debt:             newInt(200000000),
						ExtendedPairName: "CMDX-C",
						AssetInDenom:     "uasset1",
						AssetOutDenom:    "uasset2",
					},
				},
			},
		},
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultInfoOfOwnerByAppRequest{AppId: appID2, Owner: "comdex..."},
			ExpErr:            status.Errorf(codes.NotFound, "Address is not correct"),
			ExpResp:           nil,
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultInfo, err := s.querier.QueryVaultInfoOfOwnerByApp(sdk.WrapSDKContext(s.ctx), v.Req)

		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].Id, vaultInfo.VaultsInfo[0].Id)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].ExtendedPairID, vaultInfo.VaultsInfo[0].ExtendedPairID)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].Owner, vaultInfo.VaultsInfo[0].Owner)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].Collateral, vaultInfo.VaultsInfo[0].Collateral)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].Debt, vaultInfo.VaultsInfo[0].Debt)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].ExtendedPairName, vaultInfo.VaultsInfo[0].ExtendedPairName)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].AssetInDenom, vaultInfo.VaultsInfo[0].AssetInDenom)
			s.Require().Equal(v.ExpResp.VaultsInfo[0].AssetOutDenom, vaultInfo.VaultsInfo[0].AssetOutDenom)
		}
	}
}

func (s *KeeperTestSuite) TestQueryAllVaultsByAppAndExtendedPair() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryAllVaultsByAppAndExtendedPairRequest
		ExpErr              error
		ExpResp             *types.QueryAllVaultsByAppAndExtendedPairResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryAllVaultsByAppAndExtendedPairRequest{AppId: 69, ExtendedPairId: 12},
			ExpErr:            types.ErrorAppExtendedPairDataDoesNotExists,
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultsByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultsByAppAndExtendedPairResponse{
				Vault: []types.Vault{
					{
						Id:                  1,
						AppId:               appID1,
						ExtendedPairVaultID: extendedVaultPairID1,
						Owner:               addr1.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultsByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultsByAppAndExtendedPairResponse{
				Vault: []types.Vault{
					{
						Id:                  1,
						AppId:               appID1,
						ExtendedPairVaultID: extendedVaultPairID1,
						Owner:               addr1.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
					{
						Id:                  2,
						AppId:               appID1,
						ExtendedPairVaultID: extendedVaultPairID1,
						Owner:               addr2.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultsByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultsByAppAndExtendedPairResponse{
				Vault: []types.Vault{
					{
						Id:                  3,
						AppId:               appID2,
						ExtendedPairVaultID: extendedVaultPairID2,
						Owner:               addr1.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultsByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultsByAppAndExtendedPairResponse{
				Vault: []types.Vault{
					{
						Id:                  3,
						AppId:               appID2,
						ExtendedPairVaultID: extendedVaultPairID2,
						Owner:               addr1.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
					{
						Id:                  4,
						AppId:               appID2,
						ExtendedPairVaultID: extendedVaultPairID2,
						Owner:               addr2.String(),
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(200000000),
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaults, err := s.querier.QueryAllVaultsByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)

		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(len(v.ExpResp.Vault), len(vaults.Vault))
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultIDOfOwnerByExtendedPairAndApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest
		ExpErr              error
		ExpResp             *types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1, Owner: addr1.String()},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{
				Vault_Id: 1,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{
				Vault_Id: 2,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{
				Vault_Id: 3,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{
				Vault_Id: 4,
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultIDs, err := s.querier.QueryVaultIDOfOwnerByExtendedPairAndApp(sdk.WrapSDKContext(s.ctx), v.Req)

		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.Vault_Id, vaultIDs.Vault_Id)
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultIdsByAppInAllExtendedPairs() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultIdsByAppInAllExtendedPairsRequest
		ExpErr              error
		ExpResp             *types.QueryVaultIdsByAppInAllExtendedPairsResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultIdsByAppInAllExtendedPairsRequest{AppId: appID1},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultIdsByAppInAllExtendedPairsResponse{VaultIds: nil},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIdsByAppInAllExtendedPairsRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
				VaultIds: []uint64{1},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIdsByAppInAllExtendedPairsRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
				VaultIds: []uint64{1, 2},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIdsByAppInAllExtendedPairsRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
				VaultIds: []uint64{3},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultIdsByAppInAllExtendedPairsRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
				VaultIds: []uint64{3, 4},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultIDs, err := s.querier.QueryVaultIdsByAppInAllExtendedPairs(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(len(v.ExpResp.VaultIds), len(vaultIDs.VaultIds))
		}
	}
}

func (s *KeeperTestSuite) TestQueryAllVaultIdsByAnOwner() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryAllVaultIdsByAnOwnerRequest
		ExpErr              error
		ExpResp             *types.QueryAllVaultIdsByAnOwnerResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryAllVaultIdsByAnOwnerRequest{Owner: "comdex..."},
			ExpErr:            status.Errorf(codes.NotFound, "Address is not correct"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultIdsByAnOwnerRequest{Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultIdsByAnOwnerResponse{
				VaultIds: []uint64{1},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultIdsByAnOwnerRequest{Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultIdsByAnOwnerResponse{
				VaultIds: []uint64{2},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultIdsByAnOwnerRequest{Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultIdsByAnOwnerResponse{
				VaultIds: []uint64{1, 3},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryAllVaultIdsByAnOwnerRequest{Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryAllVaultIdsByAnOwnerResponse{
				VaultIds: []uint64{2, 4},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultIDs, err := s.querier.QueryAllVaultIdsByAnOwner(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(len(v.ExpResp.VaultIds), len(vaultIDs.VaultIds))
		}
	}
}

func (s *KeeperTestSuite) TestQueryTokenMintedByAppAndExtendedPair() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryTokenMintedByAppAndExtendedPairRequest
		ExpErr              error
		ExpResp             *types.QueryTokenMintedByAppAndExtendedPairResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedByAppAndExtendedPairResponse{
				TokenMinted: sdk.NewInt(200000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedByAppAndExtendedPairResponse{
				TokenMinted: sdk.NewInt(400000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedByAppAndExtendedPairResponse{
				TokenMinted: sdk.NewInt(200000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedByAppAndExtendedPairResponse{
				TokenMinted: sdk.NewInt(400000000),
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		tokenMinted, err := s.querier.QueryTokenMintedByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.TokenMinted, tokenMinted.TokenMinted)
		}
	}
}

func (s *KeeperTestSuite) TestQueryTokenMintedAssetWiseByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryTokenMintedAssetWiseByAppRequest
		ExpErr              error
		ExpResp             *types.QueryTokenMintedAssetWiseByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryTokenMintedAssetWiseByAppRequest{AppId: appID1},
			ExpErr:            nil,
			ExpResp:           &types.QueryTokenMintedAssetWiseByAppResponse{},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedAssetWiseByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedAssetWiseByAppResponse{
				MintedData: []types.MintedDataMap{
					{AssetDenom: "uasset2", MintedAmount: newInt(200000000)},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedAssetWiseByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedAssetWiseByAppResponse{
				MintedData: []types.MintedDataMap{
					{AssetDenom: "uasset2", MintedAmount: newInt(400000000)},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedAssetWiseByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedAssetWiseByAppResponse{
				MintedData: []types.MintedDataMap{
					{AssetDenom: "uasset2", MintedAmount: newInt(200000000)},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTokenMintedAssetWiseByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTokenMintedAssetWiseByAppResponse{
				MintedData: []types.MintedDataMap{
					{AssetDenom: "uasset2", MintedAmount: newInt(400000000)},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		tokenMinted, err := s.querier.QueryTokenMintedAssetWiseByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.MintedData, tokenMinted.MintedData)
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultCountByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultCountByAppRequest
		ExpErr              error
		ExpResp             *types.QueryVaultCountByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultCountByAppRequest{AppId: 69},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultCountByAppRequest{AppId: appID1},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultCountByAppResponse{},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppResponse{
				VaultCount: 1,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppResponse{
				VaultCount: 2,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppResponse{
				VaultCount: 1,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppResponse{
				VaultCount: 2,
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultCount, err := s.querier.QueryVaultCountByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.VaultCount, vaultCount.VaultCount)
		}
	}
}

func (s *KeeperTestSuite) TestQueryVaultCountByAppAndExtendedPair() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryVaultCountByAppAndExtendedPairRequest
		ExpErr              error
		ExpResp             *types.QueryVaultCountByAppAndExtendedPairResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: 69},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:            nil,
			ExpResp:           &types.QueryVaultCountByAppAndExtendedPairResponse{},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppAndExtendedPairResponse{
				VaultCount: 1,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppAndExtendedPairResponse{
				VaultCount: 2,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppAndExtendedPairResponse{
				VaultCount: 1,
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryVaultCountByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryVaultCountByAppAndExtendedPairResponse{
				VaultCount: 2,
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vaultCount, err := s.querier.QueryVaultCountByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.VaultCount, vaultCount.VaultCount)
		}
	}
}

func (s *KeeperTestSuite) TestQueryTotalValueLockedByAppAndExtendedPair() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryTotalValueLockedByAppAndExtendedPairRequest
		ExpErr              error
		ExpResp             uint64
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           0,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: 69, ExtendedPairId: 12},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           0,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:            nil,
			ExpResp:           0,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp:             1000000000,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp:             2000000000,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp:             1000000000,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTotalValueLockedByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp:             2000000000,
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		vauleLocked, err := s.querier.QueryTotalValueLockedByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			var received uint64
			if vauleLocked.ValueLocked == nil {
				received = 0
			} else {
				received = vauleLocked.ValueLocked.Uint64()
			}
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp, received)
		}
	}
}

func (s *KeeperTestSuite) TestQueryExtendedPairIDsByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryExtendedPairIDsByAppRequest
		ExpErr              error
		ExpResp             *types.QueryExtendedPairIDsByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryExtendedPairIDsByAppRequest{AppId: appID1},
			ExpErr:            nil,
			ExpResp:           &types.QueryExtendedPairIDsByAppResponse{},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairIDsByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairIDsByAppResponse{
				ExtendedPairIds: []uint64{extendedVaultPairID1},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairIDsByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairIDsByAppResponse{
				ExtendedPairIds: []uint64{extendedVaultPairID1},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairIDsByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairIDsByAppResponse{
				ExtendedPairIds: []uint64{extendedVaultPairID2},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairIDsByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairIDsByAppResponse{
				ExtendedPairIds: []uint64{extendedVaultPairID2},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		extendedPairIDs, err := s.querier.QueryExtendedPairIDsByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.ExtendedPairIds, extendedPairIDs.ExtendedPairIds)
		}
	}
}

func (s *KeeperTestSuite) TestQueryStableVaultByVaultID() {
	addr1 := s.addr(1)
	// addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, true, true)

	newVaults := []struct {
		SkipStableMintVaultCreation bool
		Address                     sdk.AccAddress
		AppID                       uint64
		ExtendedVaultPairID         uint64
		AmountIn                    sdk.Int
		Req                         *types.QueryStableVaultByVaultIDRequest
		ExpErr                      error
		ExpResp                     *types.QueryStableVaultByVaultIDResponse
	}{
		{
			SkipStableMintVaultCreation: true,
			Req:                         nil,
			ExpErr:                      status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:                     nil,
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID1,
			ExtendedVaultPairID:         extendedVaultPairID1,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByVaultIDRequest{StableVaultId: 1},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByVaultIDResponse{
				StableMintVault: &types.StableMintVault{
					Id:                  1,
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(1000000000),
					AppId:               appID1,
					ExtendedPairVaultID: extendedVaultPairID1,
				},
			},
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID2,
			ExtendedVaultPairID:         extendedVaultPairID2,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByVaultIDRequest{StableVaultId: 2},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByVaultIDResponse{
				StableMintVault: &types.StableMintVault{
					Id:                  2,
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(1000000000),
					AppId:               appID2,
					ExtendedPairVaultID: extendedVaultPairID2,
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipStableMintVaultCreation {
			s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
			_, err := s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
				v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn,
			))
			s.Require().NoError(err)
		}
		stableMintVault, err := s.querier.QueryStableVaultByVaultID(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.StableMintVault.Id, stableMintVault.StableMintVault.Id)
			s.Require().Equal(v.ExpResp.StableMintVault.AmountIn, stableMintVault.StableMintVault.AmountIn)
			s.Require().Equal(v.ExpResp.StableMintVault.AmountOut, stableMintVault.StableMintVault.AmountOut)
			s.Require().Equal(v.ExpResp.StableMintVault.AppId, stableMintVault.StableMintVault.AppId)
			s.Require().Equal(v.ExpResp.StableMintVault.ExtendedPairVaultID, stableMintVault.StableMintVault.ExtendedPairVaultID)
		}
	}
}

func (s *KeeperTestSuite) TestQueryStableVaultByApp() {
	addr1 := s.addr(1)
	// addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, true, true)

	newVaults := []struct {
		SkipStableMintVaultCreation bool
		Address                     sdk.AccAddress
		AppID                       uint64
		ExtendedVaultPairID         uint64
		AmountIn                    sdk.Int
		Req                         *types.QueryStableVaultByAppRequest
		ExpErr                      error
		ExpResp                     *types.QueryStableVaultByAppResponse
	}{
		{
			SkipStableMintVaultCreation: true,
			Req:                         nil,
			ExpErr:                      status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:                     nil,
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID1,
			ExtendedVaultPairID:         extendedVaultPairID1,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByAppRequest{AppId: appID1},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByAppResponse{
				StableMintVault: []types.StableMintVault{
					{
						Id:                  1,
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(1000000000),
						AppId:               appID1,
						ExtendedPairVaultID: extendedVaultPairID1,
					},
				},
			},
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID2,
			ExtendedVaultPairID:         extendedVaultPairID2,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByAppRequest{AppId: appID2},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByAppResponse{
				StableMintVault: []types.StableMintVault{
					{
						Id:                  2,
						AmountIn:            newInt(1000000000),
						AmountOut:           newInt(1000000000),
						AppId:               appID2,
						ExtendedPairVaultID: extendedVaultPairID2,
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipStableMintVaultCreation {
			s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
			_, err := s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
				v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn,
			))
			s.Require().NoError(err)
		}
		stableMintVaults, err := s.querier.QueryStableVaultByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.StableMintVault[0].Id, stableMintVaults.StableMintVault[0].Id)
			s.Require().Equal(v.ExpResp.StableMintVault[0].AmountIn, stableMintVaults.StableMintVault[0].AmountIn)
			s.Require().Equal(v.ExpResp.StableMintVault[0].AmountOut, stableMintVaults.StableMintVault[0].AmountOut)
			s.Require().Equal(v.ExpResp.StableMintVault[0].AppId, stableMintVaults.StableMintVault[0].AppId)
			s.Require().Equal(v.ExpResp.StableMintVault[0].ExtendedPairVaultID, stableMintVaults.StableMintVault[0].ExtendedPairVaultID)
		}
	}
}

func (s *KeeperTestSuite) TestQueryStableVaultByAppAndExtendedPair() {
	addr1 := s.addr(1)
	// addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, true, true)

	newVaults := []struct {
		SkipStableMintVaultCreation bool
		Address                     sdk.AccAddress
		AppID                       uint64
		ExtendedVaultPairID         uint64
		AmountIn                    sdk.Int
		Req                         *types.QueryStableVaultByAppAndExtendedPairRequest
		ExpErr                      error
		ExpResp                     *types.QueryStableVaultByAppAndExtendedPairResponse
	}{
		{
			SkipStableMintVaultCreation: true,
			Req:                         nil,
			ExpErr:                      status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:                     nil,
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID1,
			ExtendedVaultPairID:         extendedVaultPairID1,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByAppAndExtendedPairResponse{
				StableMintVault: &types.StableMintVault{
					Id:                  1,
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(1000000000),
					AppId:               appID1,
					ExtendedPairVaultID: extendedVaultPairID1,
				},
			},
		},
		{
			SkipStableMintVaultCreation: false,
			Address:                     addr1,
			AppID:                       appID2,
			ExtendedVaultPairID:         extendedVaultPairID2,
			AmountIn:                    newInt(1000000000),
			Req:                         &types.QueryStableVaultByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:                      nil,
			ExpResp: &types.QueryStableVaultByAppAndExtendedPairResponse{
				StableMintVault: &types.StableMintVault{
					Id:                  2,
					AmountIn:            newInt(1000000000),
					AmountOut:           newInt(1000000000),
					AppId:               appID2,
					ExtendedPairVaultID: extendedVaultPairID2,
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipStableMintVaultCreation {
			s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000))))
			_, err := s.msgServer.MsgCreateStableMint(sdk.WrapSDKContext(s.ctx), types.NewMsgCreateStableMintRequest(
				v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn,
			))
			s.Require().NoError(err)
		}
		stableMintVault, err := s.querier.QueryStableVaultByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.StableMintVault.Id, stableMintVault.StableMintVault.Id)
			s.Require().Equal(v.ExpResp.StableMintVault.AmountIn, stableMintVault.StableMintVault.AmountIn)
			s.Require().Equal(v.ExpResp.StableMintVault.AmountOut, stableMintVault.StableMintVault.AmountOut)
			s.Require().Equal(v.ExpResp.StableMintVault.AppId, stableMintVault.StableMintVault.AppId)
			s.Require().Equal(v.ExpResp.StableMintVault.ExtendedPairVaultID, stableMintVault.StableMintVault.ExtendedPairVaultID)
		}
	}
}

func (s *KeeperTestSuite) TestQueryExtendedPairVaultMappingByAppAndExtendedPair() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest
		ExpErr              error
		ExpResp             *types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{AppId: 69},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{
				ExtendedPairVaultMapping: &types.AppExtendedPairVaultMappingData{
					ExtendedPairId:         extendedVaultPairID1,
					VaultIds:               []uint64{1},
					TokenMintedAmount:      newInt(200000000),
					CollateralLockedAmount: newInt(1000000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{AppId: appID1, ExtendedPairId: extendedVaultPairID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{
				ExtendedPairVaultMapping: &types.AppExtendedPairVaultMappingData{
					ExtendedPairId:         extendedVaultPairID1,
					VaultIds:               []uint64{1, 2},
					TokenMintedAmount:      newInt(400000000),
					CollateralLockedAmount: newInt(2000000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{
				ExtendedPairVaultMapping: &types.AppExtendedPairVaultMappingData{
					ExtendedPairId:         extendedVaultPairID2,
					VaultIds:               []uint64{3},
					TokenMintedAmount:      newInt(200000000),
					CollateralLockedAmount: newInt(1000000000),
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{AppId: appID2, ExtendedPairId: extendedVaultPairID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{
				ExtendedPairVaultMapping: &types.AppExtendedPairVaultMappingData{
					ExtendedPairId:         extendedVaultPairID2,
					VaultIds:               []uint64{3, 4},
					TokenMintedAmount:      newInt(400000000),
					CollateralLockedAmount: newInt(2000000000),
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		extendedPairVault, err := s.querier.QueryExtendedPairVaultMappingByAppAndExtendedPair(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping.ExtendedPairId, extendedPairVault.ExtendedPairVaultMapping.ExtendedPairId)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping.VaultIds, extendedPairVault.ExtendedPairVaultMapping.VaultIds)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping.TokenMintedAmount, extendedPairVault.ExtendedPairVaultMapping.TokenMintedAmount)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping.CollateralLockedAmount, extendedPairVault.ExtendedPairVaultMapping.CollateralLockedAmount)
		}
	}
}

func (s *KeeperTestSuite) TestQueryExtendedPairVaultMappingByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryExtendedPairVaultMappingByAppRequest
		ExpErr              error
		ExpResp             *types.QueryExtendedPairVaultMappingByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryExtendedPairVaultMappingByAppRequest{AppId: 69},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppResponse{
				ExtendedPairVaultMapping: []types.AppExtendedPairVaultMappingData{
					{
						ExtendedPairId:         extendedVaultPairID1,
						VaultIds:               []uint64{1},
						TokenMintedAmount:      newInt(200000000),
						CollateralLockedAmount: newInt(1000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppResponse{
				ExtendedPairVaultMapping: []types.AppExtendedPairVaultMappingData{
					{
						ExtendedPairId:         extendedVaultPairID1,
						VaultIds:               []uint64{1, 2},
						TokenMintedAmount:      newInt(400000000),
						CollateralLockedAmount: newInt(2000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppResponse{
				ExtendedPairVaultMapping: []types.AppExtendedPairVaultMappingData{
					{
						ExtendedPairId:         extendedVaultPairID2,
						VaultIds:               []uint64{3},
						TokenMintedAmount:      newInt(200000000),
						CollateralLockedAmount: newInt(1000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryExtendedPairVaultMappingByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryExtendedPairVaultMappingByAppResponse{
				ExtendedPairVaultMapping: []types.AppExtendedPairVaultMappingData{
					{
						ExtendedPairId:         extendedVaultPairID2,
						VaultIds:               []uint64{3, 4},
						TokenMintedAmount:      newInt(400000000),
						CollateralLockedAmount: newInt(2000000000),
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		extendedPairVault, err := s.querier.QueryExtendedPairVaultMappingByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping[0].ExtendedPairId, extendedPairVault.ExtendedPairVaultMapping[0].ExtendedPairId)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping[0].VaultIds, extendedPairVault.ExtendedPairVaultMapping[0].VaultIds)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping[0].TokenMintedAmount, extendedPairVault.ExtendedPairVaultMapping[0].TokenMintedAmount)
			s.Require().Equal(v.ExpResp.ExtendedPairVaultMapping[0].CollateralLockedAmount, extendedPairVault.ExtendedPairVaultMapping[0].CollateralLockedAmount)
		}
	}
}

func (s *KeeperTestSuite) TestQueryTVLByAppOfAllExtendedPairs() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryTVLByAppOfAllExtendedPairsRequest
		ExpErr              error
		ExpResp             *types.QueryTVLByAppOfAllExtendedPairsResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppOfAllExtendedPairsRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppOfAllExtendedPairsResponse{
				Tvldata: []types.TvlLockedDataMap{
					{
						AssetDenom:             "uasset1",
						CollateralLockedAmount: newInt(1000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppOfAllExtendedPairsRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppOfAllExtendedPairsResponse{
				Tvldata: []types.TvlLockedDataMap{
					{
						AssetDenom:             "uasset1",
						CollateralLockedAmount: newInt(2000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppOfAllExtendedPairsRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppOfAllExtendedPairsResponse{
				Tvldata: []types.TvlLockedDataMap{
					{
						AssetDenom:             "uasset1",
						CollateralLockedAmount: newInt(1000000000),
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppOfAllExtendedPairsRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppOfAllExtendedPairsResponse{
				Tvldata: []types.TvlLockedDataMap{
					{
						AssetDenom:             "uasset1",
						CollateralLockedAmount: newInt(2000000000),
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		tvlData, err := s.querier.QueryTVLByAppOfAllExtendedPairs(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.Tvldata[0].AssetDenom, tvlData.Tvldata[0].AssetDenom)
			s.Require().Equal(v.ExpResp.Tvldata[0].CollateralLockedAmount, tvlData.Tvldata[0].CollateralLockedAmount)
		}
	}
}

func (s *KeeperTestSuite) TestQueryTVLByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryTVLByAppRequest
		ExpErr              error
		ExpResp             *types.QueryTVLByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryTVLByAppRequest{AppId: 69},
			ExpErr:            status.Errorf(codes.NotFound, "App does not exist for id %d", 69),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppResponse{
				CollateralLocked: newInt(1000000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppResponse{
				CollateralLocked: newInt(2000000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppResponse{
				CollateralLocked: newInt(1000000000),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryTVLByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryTVLByAppResponse{
				CollateralLocked: newInt(2000000000),
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		tvlData, err := s.querier.QueryTVLByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.CollateralLocked, tvlData.CollateralLocked)
		}
	}
}

func (s *KeeperTestSuite) TestQueryUserMyPositionByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryUserMyPositionByAppRequest
		ExpErr              error
		ExpResp             *types.QueryUserMyPositionByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation: true,
			Req:               &types.QueryUserMyPositionByAppRequest{AppId: appID1, Owner: "comdex..."},
			ExpErr:            status.Errorf(codes.NotFound, "Address is not correct"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserMyPositionByAppRequest{AppId: appID1, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserMyPositionByAppResponse{
				CollateralLocked:  newInt(1000000000),
				TotalDue:          newInt(400000000),
				AvailableToBorrow: newInt(34782608),
				AverageCrRatio:    sdk.MustNewDecFromStr("2.5"),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserMyPositionByAppRequest{AppId: appID1, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserMyPositionByAppResponse{
				CollateralLocked:  newInt(1000000000),
				TotalDue:          newInt(400000000),
				AvailableToBorrow: newInt(34782608),
				AverageCrRatio:    sdk.MustNewDecFromStr("2.5"),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserMyPositionByAppRequest{AppId: appID2, Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserMyPositionByAppResponse{
				CollateralLocked:  newInt(1000000000),
				TotalDue:          newInt(400000000),
				AvailableToBorrow: newInt(34782608),
				AverageCrRatio:    sdk.MustNewDecFromStr("2.5"),
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserMyPositionByAppRequest{AppId: appID2, Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserMyPositionByAppResponse{
				CollateralLocked:  newInt(1000000000),
				TotalDue:          newInt(400000000),
				AvailableToBorrow: newInt(34782608),
				AverageCrRatio:    sdk.MustNewDecFromStr("2.5"),
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		userPositionData, err := s.querier.QueryUserMyPositionByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.CollateralLocked, userPositionData.CollateralLocked)
			s.Require().Equal(v.ExpResp.TotalDue, userPositionData.TotalDue)
			s.Require().Equal(v.ExpResp.AvailableToBorrow, userPositionData.AvailableToBorrow)
			s.Require().Equal(v.ExpResp.AverageCrRatio, userPositionData.AverageCrRatio)
		}
	}
}

func (s *KeeperTestSuite) TestQueryUserExtendedPairTotalData() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryUserExtendedPairTotalDataRequest
		ExpErr              error
		ExpResp             *types.QueryUserExtendedPairTotalDataResponse
	}{
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserExtendedPairTotalDataRequest{Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserExtendedPairTotalDataResponse{
				UserTotalData: []types.OwnerAppExtendedPairVaultMappingData{
					{
						Owner:          addr1.String(),
						AppId:          appID1,
						ExtendedPairId: extendedVaultPairID1,
						VaultId:        1,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserExtendedPairTotalDataRequest{Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserExtendedPairTotalDataResponse{
				UserTotalData: []types.OwnerAppExtendedPairVaultMappingData{
					{
						Owner:          addr2.String(),
						AppId:          appID1,
						ExtendedPairId: extendedVaultPairID1,
						VaultId:        2,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserExtendedPairTotalDataRequest{Owner: addr1.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserExtendedPairTotalDataResponse{
				UserTotalData: []types.OwnerAppExtendedPairVaultMappingData{
					{Owner: addr1.String(),
						AppId:          appID1,
						ExtendedPairId: extendedVaultPairID1,
						VaultId:        1,
					},
					{Owner: addr1.String(),
						AppId:          appID2,
						ExtendedPairId: extendedVaultPairID2,
						VaultId:        3,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryUserExtendedPairTotalDataRequest{Owner: addr2.String()},
			ExpErr:              nil,
			ExpResp: &types.QueryUserExtendedPairTotalDataResponse{
				UserTotalData: []types.OwnerAppExtendedPairVaultMappingData{
					{
						Owner:          addr2.String(),
						AppId:          appID1,
						ExtendedPairId: extendedVaultPairID1,
						VaultId:        2,
					},
					{
						Owner:          addr2.String(),
						AppId:          appID2,
						ExtendedPairId: extendedVaultPairID2,
						VaultId:        4,
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		userTotalData, err := s.querier.QueryUserExtendedPairTotalData(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.UserTotalData, userTotalData.UserTotalData)
		}
	}
}

func (s *KeeperTestSuite) TestQueryPairsLockedAndMintedStatisticByApp() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)
	extendedVaultPairID2 := s.CreateNewExtendedVaultPair("CMDX-C", appID2, pairID, false, true)

	newVaults := []struct {
		SkipVaultCreation   bool
		Address             sdk.AccAddress
		AppID               uint64
		ExtendedVaultPairID uint64
		AmountIn            sdk.Int
		AmountOut           sdk.Int
		Req                 *types.QueryPairsLockedAndMintedStatisticByAppRequest
		ExpErr              error
		ExpResp             *types.QueryPairsLockedAndMintedStatisticByAppResponse
	}{
		{
			SkipVaultCreation: true,
			Req:               nil,
			ExpErr:            status.Error(codes.InvalidArgument, "request cannot be empty"),
			ExpResp:           nil,
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryPairsLockedAndMintedStatisticByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryPairsLockedAndMintedStatisticByAppResponse{
				PairStatisticData: []types.PairStatisticData{
					{
						AssetInDenom:        "uasset1",
						AssetOutDenom:       "uasset2",
						CollateralAmount:    newInt(1000000000),
						MintedAmount:        newInt(200000000),
						ExtendedPairVaultID: extendedVaultPairID1,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID1,
			ExtendedVaultPairID: extendedVaultPairID1,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryPairsLockedAndMintedStatisticByAppRequest{AppId: appID1},
			ExpErr:              nil,
			ExpResp: &types.QueryPairsLockedAndMintedStatisticByAppResponse{
				PairStatisticData: []types.PairStatisticData{
					{
						AssetInDenom:        "uasset1",
						AssetOutDenom:       "uasset2",
						CollateralAmount:    newInt(2000000000),
						MintedAmount:        newInt(400000000),
						ExtendedPairVaultID: extendedVaultPairID1,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr1,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryPairsLockedAndMintedStatisticByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryPairsLockedAndMintedStatisticByAppResponse{
				PairStatisticData: []types.PairStatisticData{
					{
						AssetInDenom:        "uasset1",
						AssetOutDenom:       "uasset2",
						CollateralAmount:    newInt(1000000000),
						MintedAmount:        newInt(200000000),
						ExtendedPairVaultID: extendedVaultPairID2,
					},
				},
			},
		},
		{
			SkipVaultCreation:   false,
			Address:             addr2,
			AppID:               appID2,
			ExtendedVaultPairID: extendedVaultPairID2,
			AmountIn:            newInt(1000000000),
			AmountOut:           newInt(200000000),
			Req:                 &types.QueryPairsLockedAndMintedStatisticByAppRequest{AppId: appID2},
			ExpErr:              nil,
			ExpResp: &types.QueryPairsLockedAndMintedStatisticByAppResponse{
				PairStatisticData: []types.PairStatisticData{
					{
						AssetInDenom:        "uasset1",
						AssetOutDenom:       "uasset2",
						CollateralAmount:    newInt(2000000000),
						MintedAmount:        newInt(400000000),
						ExtendedPairVaultID: extendedVaultPairID2,
					},
				},
			},
		},
	}

	for _, v := range newVaults {
		if !v.SkipVaultCreation {
			msg := types.NewMsgCreateRequest(v.Address, v.AppID, v.ExtendedVaultPairID, v.AmountIn, v.AmountOut)
			s.fundAddr(sdk.MustAccAddressFromBech32(v.Address.String()), sdk.NewCoins(sdk.NewCoin("uasset1", v.AmountIn)))
			_, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), msg)
			s.Require().NoError(err)
		}
		pairsStatistics, err := s.querier.QueryPairsLockedAndMintedStatisticByApp(sdk.WrapSDKContext(s.ctx), v.Req)
		if v.ExpErr != nil {
			s.Require().Error(err)
			s.Require().EqualError(v.ExpErr, err.Error())
		} else {
			s.Require().NoError(err)
			s.Require().Equal(v.ExpResp.PairStatisticData, pairsStatistics.PairStatisticData)
		}
	}
}
