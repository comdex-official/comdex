package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenmintTypes "github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddAppAsset() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	genesisSupply := sdkmath.NewIntFromUint64(9000000)
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdkmath.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				3,
				genesisSupply,
				true,
				userAddress,
			},
			{
				2,
				genesisSupply,
				true,
				userAddress,
			},
		},
	}
	err := assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := assetTypes.AppData{
		Name:             "commodo",
		ShortName:        "comdo",
		MinGovDeposit:    sdkmath.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				3,
				genesisSupply,
				true,
				userAddress,
			},
		},
	}
	err = assetKeeper.AddAppRecords(*ctx, msg2)
	s.Require().NoError(err)

	msg20 := assetTypes.AppData{
		Name:             "cool",
		ShortName:        "cool",
		MinGovDeposit:    sdkmath.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken:     nil,
	}
	err = assetKeeper.AddAppRecords(*ctx, msg20)
	s.Require().NoError(err)

	msg3 := assetTypes.Asset{
		Name:      "CMDX",
		Denom:     "ucmdx",
		Decimals:  sdkmath.NewInt(1000000),
		IsOnChain: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg3)
	s.Require().NoError(err)

	msg4 := assetTypes.Asset{
		Name:      "CMST",
		Denom:     "ucmst",
		Decimals:  sdkmath.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg4)
	s.Require().NoError(err)

	msg5 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  sdkmath.NewInt(1000000),
		IsOnChain: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg5)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestMsgMintNewTokens() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	tokenmintKeeper, ctx := &s.tokenmintKeeper, &s.ctx
	wctx := sdk.WrapSDKContext(*ctx)
	s.AddAppAsset()
	server := keeper.NewMsgServer(*tokenmintKeeper)
	for _, tc := range []struct {
		name          string
		msg           tokenmintTypes.MsgMintNewTokensRequest
		expectedError bool
	}{
		{
			"Mint New Tokens : App ID : 1, Asset ID : 3",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 3,
			},
			false,
		},
		{
			"Mint New Tokens : App ID : 1, Asset ID : 2",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 2,
			},
			false,
		},
		{
			"Duplicate Failure Mint New Tokens : App ID : 1, Asset ID : 3",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 3,
			},
			true,
		},
		{
			"App Mapping Does Not Exists",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   11,
				AssetId: 3,
			},
			true,
		},
		{
			"Asset Not WhiteListed For Genesis Minting",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 1,
			},
			true,
		},
		{
			"Mint New Tokens : App ID : 2, Asset ID : 3",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   2,
				AssetId: 3,
			},
			false,
		},
	} {
		s.Run(tc.name, func() {
			genesisSupply := sdkmath.NewIntFromUint64(9000000)
			asset, found := s.assetKeeper.GetAsset(*ctx, tc.msg.AssetId)
			s.Require().True(found)
			previousCoin, err := s.getBalance(userAddress, asset.Denom)
			_, err = server.MsgMintNewTokens(wctx, &tc.msg)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				currentCoin, err := s.getBalance(userAddress, asset.Denom)
				s.Require().NoError(err)
				ActualAmountMinted := currentCoin.Amount.Sub(previousCoin.Amount)
				s.Require().Equal(ActualAmountMinted, genesisSupply)
				req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
					AppId:   tc.msg.AppId,
					AssetId: tc.msg.AssetId,
				}
				res, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
				s.Require().NoError(err)
				s.Require().Equal(res.MintedTokens.AssetId, tc.msg.AssetId)
				s.Require().Equal(res.MintedTokens.GenesisSupply, ActualAmountMinted)
				s.Require().Equal(res.MintedTokens.CurrentSupply, ActualAmountMinted)
			}
		})
	}
	result := s.tokenmintKeeper.GetTotalTokenMinted(*ctx)
	// validates no. of apps
	s.Require().Equal(len(result), 2)
	// validates no of assets under app id 1
	s.Require().Equal(len(result[0].MintedTokens), 2)
	// validates no of assets under app id 2
	s.Require().Equal(len(result[1].MintedTokens), 1)
}

func (s *KeeperTestSuite) TestMintNewTokensForApp() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	tokenmintKeeper, ctx := &s.tokenmintKeeper, &s.ctx
	wctx := sdk.WrapSDKContext(*ctx)
	s.TestMsgMintNewTokens()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		address       string
		mintAmount    sdkmath.Int
		expectedError bool
	}{
		{
			"Mint New Tokens : App ID : 1, Asset ID : 2",
			1,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(423),
			false,
		},
		{
			"App Mapping Does Not Exists",
			33,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(423),
			true,
		},
		{
			"App Mapping Does Not Exists",
			3,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(423),
			true,
		},
	} {
		s.Run(tc.name, func() {
			asset, found := s.assetKeeper.GetAsset(*ctx, tc.assetID)
			s.Require().True(found)
			previousCoin, err2 := s.getBalance(userAddress, asset.Denom)
			s.Require().NoError(err2)
			err := tokenmintKeeper.MintNewTokensForApp(*ctx, tc.appID, tc.assetID, tc.address, tc.mintAmount)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
					AppId:   tc.appID,
					AssetId: tc.assetID,
				}
				beforeTokenMint, err1 := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
				s.Require().NoError(err1)
				currentCoin, err := s.getBalance(userAddress, asset.Denom)
				s.Require().NoError(err)
				ActualAmountMinted := currentCoin.Amount.Sub(previousCoin.Amount)
				fmt.Println("ActualAmountMinted....11..", ActualAmountMinted)
				fmt.Println("tc.mintAmount...11..", tc.mintAmount)
				s.Require().Equal(ActualAmountMinted, tc.mintAmount)
				req2 := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
					AppId:   tc.appID,
					AssetId: tc.assetID,
				}
				res, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req2)
				s.Require().NoError(err)
				s.Require().Equal(res.MintedTokens.AssetId, tc.assetID)
				s.Require().Equal(res.MintedTokens.GenesisSupply, beforeTokenMint.MintedTokens.GenesisSupply)
				s.Require().Equal(tc.mintAmount, ActualAmountMinted)
				//s.Require().Equal(res.MintedTokens.CurrentSupply, beforeTokenMint.MintedTokens.CurrentSupply.Add(tc.mintAmount))
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnTokensForApp() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	tokenmintKeeper, ctx := &s.tokenmintKeeper, &s.ctx
	wctx := sdk.WrapSDKContext(*ctx)
	s.TestMsgMintNewTokens()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		address       string
		sendAmount    sdkmath.Int
		burnAmount    sdkmath.Int
		expectedError bool
	}{
		{
			"Burn Tokens : App ID : 1, Asset ID : 2",
			1,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(423),
			sdkmath.NewIntFromUint64(423),
			false,
		},
		{
			"Burn Tokens Insuffient balance failure: App ID : 1, Asset ID : 2",
			1,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(422),
			sdkmath.NewIntFromUint64(423),
			true,
		},
	} {
		s.Run(tc.name, func() {
			asset, found := s.assetKeeper.GetAsset(*ctx, tc.assetID)
			s.Require().True(found)
			sender, err := sdk.AccAddressFromBech32(tc.address)
			s.Require().NoError(err)
			err = s.app.BankKeeper.SendCoinsFromAccountToModule(*ctx, sender, "tokenmint", sdk.NewCoins(sdk.NewCoin(asset.Denom, tc.sendAmount)))
			s.Require().NoError(err)
			req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
				AppId:   tc.appID,
				AssetId: tc.assetID,
			}
			beforeTokenMint, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
			s.Require().NoError(err)
			beforeTokenMintBalance := s.auctionKeeper.GetModuleAccountBalance(*ctx, "tokenmint", asset.Denom)
			err = tokenmintKeeper.BurnTokensForApp(*ctx, tc.appID, tc.assetID, tc.burnAmount)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				afterTokenMintBalance := s.auctionKeeper.GetModuleAccountBalance(*ctx, "tokenmint", asset.Denom)
				s.Require().NoError(err)
				req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
					AppId:   tc.appID,
					AssetId: tc.assetID,
				}
				res, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
				s.Require().NoError(err)
				s.Require().Equal(res.MintedTokens.AssetId, tc.assetID)
				s.Require().Equal(res.MintedTokens.GenesisSupply, beforeTokenMint.MintedTokens.GenesisSupply)
				s.Require().Equal(res.MintedTokens.CurrentSupply.Add(tc.burnAmount), beforeTokenMint.MintedTokens.CurrentSupply)
				s.Require().Equal(beforeTokenMintBalance.Sub(afterTokenMintBalance), tc.burnAmount)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnGovTokensForApp() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	tokenmintKeeper, ctx := &s.tokenmintKeeper, &s.ctx
	wctx := sdk.WrapSDKContext(*ctx)
	s.TestMsgMintNewTokens()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		address       string
		burnAmount    sdkmath.Int
		expectedError bool
	}{
		{
			"Burn Gov Tokens For App: App ID : 1, Asset ID : 2",
			1,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(423),
			false,
		},
		{
			"Burn Gov Tokens For App Insuffient balance failure: App ID : 1, Asset ID : 2",
			1,
			2,
			userAddress,
			sdkmath.NewIntFromUint64(9000001),
			true,
		},
	} {
		s.Run(tc.name, func() {
			asset, found := s.assetKeeper.GetAsset(*ctx, tc.assetID)
			s.Require().True(found)
			sender, err := sdk.AccAddressFromBech32(tc.address)
			s.Require().NoError(err)

			req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
				AppId:   tc.appID,
				AssetId: tc.assetID,
			}
			beforeTokenMint, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
			s.Require().NoError(err)
			beforeUserBalance, err := s.getBalance(tc.address, asset.Denom)
			s.Require().NoError(err)
			beforeTokenMintBalance := s.auctionKeeper.GetModuleAccountBalance(*ctx, "tokenmint", asset.Denom)
			err = tokenmintKeeper.BurnGovTokensForApp(*ctx, tc.appID, sender, sdk.NewCoin(asset.Denom, tc.burnAmount))

			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				afterTokenMintBalance := s.auctionKeeper.GetModuleAccountBalance(*ctx, "tokenmint", asset.Denom)
				s.Require().NoError(err)
				afterUserBalance, err := s.getBalance(tc.address, asset.Denom)
				s.Require().NoError(err)
				req := tokenmintTypes.QueryTokenMintedByAppAndAssetRequest{
					AppId:   tc.appID,
					AssetId: tc.assetID,
				}
				res, err := s.querier.QueryTokenMintedByAppAndAsset(wctx, &req)
				s.Require().NoError(err)
				currentSupply, found := s.tokenmintKeeper.GetAssetDataInTokenMintByAppSupply(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(res.MintedTokens.AssetId, tc.assetID)
				s.Require().Equal(res.MintedTokens.GenesisSupply, beforeTokenMint.MintedTokens.GenesisSupply)
				s.Require().Equal(sdkmath.NewInt(currentSupply).Add(tc.burnAmount), beforeTokenMint.MintedTokens.CurrentSupply)
				s.Require().Equal(beforeTokenMintBalance, afterTokenMintBalance)
				s.Require().Equal(beforeUserBalance.Amount.Sub(afterUserBalance.Amount), tc.burnAmount)
			}
		})
	}
}
