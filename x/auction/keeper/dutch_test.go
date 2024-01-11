package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auction"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctionTypes "github.com/comdex-official/comdex/x/auction/types"
	collectorTypes "github.com/comdex-official/comdex/x/collector/types"
	liquidationTypes "github.com/comdex-official/comdex/x/liquidation/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper1 "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddPairAndExtendedPairVault1() {
	assetKeeper, liquidationKeeper, ctx := &s.assetKeeper, &s.liquidationKeeper, &s.ctx

	for _, tc := range []struct {
		name              string
		pair              assetTypes.Pair
		extendedPairVault bindings.MsgAddExtendedPairsVault
		asset1            uint64
		asset2            uint64
	}{
		{
			"Add Pair , Extended Pair Vault : cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
			bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              1,
				StabilityFee:        sdk.MustNewDecFromStr("0.01"),
				ClosingFee:          sdk.MustNewDecFromStr("0"),
				LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
				DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
				IsVaultActive:       true,
				DebtCeiling:         sdk.NewInt(1000000000000),
				DebtFloor:           sdk.NewInt(1000000),
				IsStableMintVault:   false,
				MinCr:               sdk.MustNewDecFromStr("1.5"),
				PairName:            "CMDX-B",
				AssetOutOraclePrice: true,
				AssetOutPrice:       1000000,
				MinUsdValueLeft:     1000000,
			},
			1,
			2,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddPairsRecords(*ctx, tc.pair)
			s.Require().NoError(err)

			err = assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.extendedPairVault)
			s.Require().NoError(err)

			err = liquidationKeeper.WasmWhitelistAppIDLiquidation(*ctx, 1)
			s.Require().NoError(err)

			s.SetInitialOraclePriceForID(tc.asset1, tc.asset2)
		})
	}
}

func (s *KeeperTestSuite) SetOraclePrice(assetID uint64, price uint64) {
	market := markettypes.TimeWeightedAverage{
		AssetID:       assetID,
		ScriptID:      12,
		Twa:           price,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{price},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market)
}

func (s *KeeperTestSuite) SetInitialOraclePriceForID(asset1 uint64, asset2 uint64) {
	s.SetOraclePrice(asset1, 2000000)
	s.SetOraclePrice(asset2, 1000000)
}

func (s *KeeperTestSuite) ChangeOraclePrice(asset uint64) {
	s.SetOraclePrice(asset, 1000000)
}

func (s *KeeperTestSuite) CreateVault() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	server := vaultKeeper1.NewMsgServer(*vaultKeeper)

	for index, tc := range []struct {
		name string
		msg  vaultTypes.MsgCreateRequest
	}{
		{
			"Create Vault : AppID 1 extended pair 1 user address 1",
			vaultTypes.MsgCreateRequest{
				From:                userAddress1,
				AppId:               1,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
		{
			"Create Vault : AppID 1 extended pair 1 user address 2",
			vaultTypes.MsgCreateRequest{
				From:                userAddress2,
				AppId:               1,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := server.MsgCreate(sdk.WrapSDKContext(*ctx), &tc.msg)
			s.Require().NoError(err)
			res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
			s.Require().NoError(err)
			_, err = s.vaultQuerier.QueryVaultInfoByVaultID(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryVaultInfoByVaultIDRequest{Id: res.Vault[index].Id})
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) GetVaultCount() int {
	ctx := &s.ctx
	res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
	s.Require().NoError(err)
	return len(res.Vault)
}

func (s *KeeperTestSuite) GetVaultCountForExtendedPairIDbyAppID(appID, extID uint64) int {
	vaultKeeper1, ctx := &s.vaultKeeper, &s.ctx
	res, found := vaultKeeper1.GetAppExtendedPairVaultMappingData(*ctx, appID, extID)
	s.Require().True(found)
	return len(res.VaultIds)
}

func (s *KeeperTestSuite) GetAssetPrice(id uint64) sdk.Dec {
	marketKeeper, ctx := &s.marketKeeper, &s.ctx
	price, found := marketKeeper.GetTwa(*ctx, id)
	s.Require().True(found)
	price1 := sdk.NewDecFromInt(sdk.NewIntFromUint64(price.Twa))
	return price1
}

func (s *KeeperTestSuite) AddAppAsset() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	addr, err := sdk.AccAddressFromBech32(userAddress1)
	s.Require().NoError(err)
	addr2, err := sdk.AccAddressFromBech32(userAddress2)
	s.Require().NoError(err)
	genesisSupply := sdk.NewIntFromUint64(1000000)
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				AssetId:       3,
				GenesisSupply: genesisSupply,
				IsGovToken:    true,
				Recipient:     userAddress1,
			},
			{
				AssetId:       2,
				GenesisSupply: genesisSupply,
				IsGovToken:    true,
				Recipient:     userAddress1,
			},
		},
	}
	err = assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name string
		msg  assetTypes.Asset
	}{
		{
			"Add Asset 1",
			assetTypes.Asset{
				Name:          "CMDX",
				Denom:         "ucmdx",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
		{
			"Add Asset 2",
			assetTypes.Asset{
				Name:          "CMST",
				Denom:         "ucmst",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
		{
			"Add Asset 3",
			assetTypes.Asset{
				Name:          "HARBOR",
				Denom:         "uharbor",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAssetRecords(*ctx, tc.msg)
			s.Require().NoError(err)
			s.fundAddr(addr, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
			s.fundAddr(addr2, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
		})
	}
}

func (s *KeeperTestSuite) LiquidateVaults1() {
	liquidationKeeper, vaultKeeper, ctx := &s.liquidationKeeper, &s.vaultKeeper, &s.ctx
	s.CreateVault()
	currentVaultsCount := 2
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1), 1), currentVaultsCount)
	beforeVault, found := vaultKeeper.GetVault(*ctx, 1)
	s.Require().True(found)

	s.AddAuctionParams()

	// Liquidation shouldn't happen as price not changed
	err := liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(0))

	// Liquidation should happen as price changed
	s.ChangeOraclePrice(1)
	err = liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	id = liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(2))
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount-2)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1), 1), currentVaultsCount-2)

	lockedVault := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, beforeVault.Id)
	s.Require().Equal(lockedVault[0].ExtendedPairId, beforeVault.ExtendedPairVaultID)
	s.Require().Equal(lockedVault[0].Owner, beforeVault.Owner)
	s.Require().Equal(lockedVault[0].AmountIn, beforeVault.AmountIn)
	s.Require().Equal(lockedVault[0].AmountOut, beforeVault.AmountOut)
	s.Require().Equal(lockedVault[0].UpdatedAmountOut, sdk.ZeroInt())
	s.Require().Equal(lockedVault[0].Initiator, liquidationTypes.ModuleName)
	s.Require().Equal(lockedVault[0].IsAuctionInProgress, true)
	s.Require().Equal(lockedVault[0].IsAuctionComplete, false)
	s.Require().Equal(lockedVault[0].SellOffHistory, []string(nil))
	price, err := s.app.MarketKeeper.CalcAssetPrice(*ctx, uint64(1), beforeVault.AmountIn)
	s.Require().NoError(err)
	s.Require().Equal(lockedVault[0].CollateralToBeAuctioned, price)
	s.Require().Equal(lockedVault[0].CrAtLiquidation, sdk.NewDecFromInt(lockedVault[0].AmountIn).Mul(s.GetAssetPrice(1)).Quo(sdk.NewDecFromInt(lockedVault[0].AmountOut).Mul(s.GetAssetPrice(2))))
}

func (s *KeeperTestSuite) AddAuctionParams() {
	ctx := &s.ctx
	auctionParams := auctionTypes.AuctionParams{
		AppId:                  1,
		AuctionDurationSeconds: 300,
		Buffer:                 sdk.MustNewDecFromStr("1.2"),
		Cusp:                   sdk.MustNewDecFromStr("0.6"),
		Step:                   sdk.NewIntFromUint64(1),
		PriceFunctionType:      1,
		SurplusId:              1,
		DebtId:                 2,
		DutchId:                3,
		BidDurationSeconds:     300,
	}
	s.app.AuctionKeeper.SetAuctionParams(*ctx, auctionParams)
}

func (s *KeeperTestSuite) TestDutchActivator() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	addr1, err := sdk.AccAddressFromBech32(userAddress1)
	s.Require().NoError(err)
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.LiquidateVaults1()
	s.AddAuctionParams()
	k, liquidationKeeper, ctx := &s.keeper, &s.liquidationKeeper, &s.ctx

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)

	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	lockedVault, found := liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	err = k.DutchActivator(*ctx, lockedVault)
	s.Require().Error(err)
	dutchAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(dutchAuction.AppId, lockedVault.AppId)
	s.Require().Equal(dutchAuction.AuctionId, auctionId)
	s.Require().Equal(dutchAuction.AuctionMappingId, auctionMappingId)
	s.Require().Equal(dutchAuction.OutflowTokenInitAmount.Amount, lockedVault.AmountIn)
	s.Require().Equal(dutchAuction.OutflowTokenCurrentAmount.Amount, lockedVault.AmountIn)
	s.Require().Equal(dutchAuction.InflowTokenCurrentAmount.Amount, sdk.ZeroInt())

	inFlowTokenTargetAmount := lockedVault.AmountOut
	mulfactor := sdk.NewDecFromInt(inFlowTokenTargetAmount).Mul(dutchAuction.LiquidationPenalty)
	inFlowTokenTargetAmount = inFlowTokenTargetAmount.Add(mulfactor.TruncateInt()).Add(lockedVault.InterestAccumulated)

	s.Require().Equal(dutchAuction.InflowTokenTargetAmount.Amount, inFlowTokenTargetAmount)

	s.Require().Equal(dutchAuction.VaultOwner, addr1)
	s.Require().Equal(dutchAuction.AuctionStatus, auctionTypes.AuctionStartNoBids)

	assetOutPrice, found := s.marketKeeper.GetTwa(*ctx, dutchAuction.AssetOutId)
	s.Require().True(found)
	assetInPrice, found := s.marketKeeper.GetTwa(*ctx, dutchAuction.AssetInId)
	s.Require().True(found)
	s.Require().Equal(dutchAuction.OutflowTokenCurrentPrice, sdk.NewDecFromInt(sdk.NewIntFromUint64(assetOutPrice.Twa)).Mul(sdk.MustNewDecFromStr("1.2")))
	s.Require().Equal(dutchAuction.OutflowTokenEndPrice, dutchAuction.OutflowTokenInitialPrice.Mul(sdk.MustNewDecFromStr("0.6")))
	s.Require().Equal(dutchAuction.InflowTokenCurrentPrice, sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice.Twa)))
}

func (s *KeeperTestSuite) TestDutchBid() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	s.TestDutchActivator()
	k, ctx := &s.keeper, &s.ctx
	appID := uint64(1)
	auctionMappingID := uint64(3)
	auctionID := uint64(1)
	biddingID := uint64(1)

	server := auctionKeeper.NewMsgServiceServer(*k)

	for _, tc := range []struct {
		name            string
		msg             auctionTypes.MsgPlaceDutchBidRequest
		advanceSeconds  int64
		isErrorExpected bool
	}{
		{
			"Place Dutch Bid : bid more than collateral auctioned max 1.4",
			auctionTypes.MsgPlaceDutchBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1000001ucmdx"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			0,
			true,
		},
		{
			"Place Dutch Bid : invalid bid denom max 1.4",
			auctionTypes.MsgPlaceDutchBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1000ucmst"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			0,
			true,
		},
		{
			"Place Dutch Bid : leaving dust should fail max 1.4",
			auctionTypes.MsgPlaceDutchBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1ucmdx"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			250,
			true,
		},
		{
			"Place Dutch Bid : collateral max price less than current price 1000ucmdx max 1.1",
			auctionTypes.MsgPlaceDutchBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1000ucmdx"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			0,
			true,
		},
		{
			"Place Dutch Bid : collateral max price more than current price 1000ucmdx max 1.4",
			auctionTypes.MsgPlaceDutchBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1000ucmdx"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			0,
			false,
		},
	} {
		s.Run(tc.name, func() {
			s.advanceseconds(tc.advanceSeconds)
			auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
			beforeAuction, err := k.GetDutchAuction(*ctx, appID, auctionMappingID, auctionID)
			s.Require().NoError(err)
			beforeCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
			s.Require().NoError(err)
			beforeCmstBalance, err := s.getBalance(userAddress1, "ucmst")
			s.Require().NoError(err)

			// dont expect error
			_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.isErrorExpected {
				s.advanceseconds(301 - tc.advanceSeconds)
				auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
			} else {
				s.Require().NoError(err)

				afterCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
				s.Require().NoError(err)
				afterCmstBalance, err := s.getBalance(userAddress1, "ucmst")
				s.Require().NoError(err)

				afterAuction, err := k.GetDutchAuction(*ctx, appID, auctionMappingID, auctionID)
				s.Require().NoError(err)

				userBid, err := k.GetDutchUserBidding(*ctx, userAddress1, appID, biddingID)
				userReceivableAmount := sdk.NewDecFromInt(tc.msg.Amount.Amount).Mul(beforeAuction.OutflowTokenCurrentPrice).Quo(beforeAuction.InflowTokenCurrentPrice).TruncateInt()
				userOutflowCoin := sdk.NewCoin("ucmst", userReceivableAmount)
				userInflowCoin := tc.msg.Amount
				s.Require().Equal(beforeAuction.OutflowTokenCurrentAmount.Sub(userInflowCoin), afterAuction.OutflowTokenCurrentAmount)
				s.Require().Equal(beforeAuction.InflowTokenCurrentAmount.Add(userOutflowCoin), afterAuction.InflowTokenCurrentAmount)
				s.Require().Equal(beforeCmdxBalance.Add(userInflowCoin), afterCmdxBalance)
				s.Require().Equal(beforeCmstBalance.Sub(userOutflowCoin), afterCmstBalance)
				s.Require().Equal(userBid.BiddingId, biddingID)
				s.Require().Equal(userBid.AppId, appID)
				s.Require().Equal(userBid.AuctionId, auctionID)
				s.Require().Equal(userBid.BiddingStatus, auctionTypes.SuccessBiddingStatus)
				s.Require().Equal(userBid.AuctionStatus, auctionTypes.ActiveAuctionStatus)
				s.Require().Equal(userBid.Bidder, userAddress1)
				s.Require().Equal(userBid.AuctionMappingId, auctionMappingID)
				s.Require().Equal(userBid.OutflowTokenAmount, userOutflowCoin)
				s.Require().Equal(userBid.InflowTokenAmount, userInflowCoin)
			}
		})
	}
}

func (s *KeeperTestSuite) TestCloseDutchAuction() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	s.TestDutchBid()
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	server := auctionKeeper.NewMsgServiceServer(*k)
	beforeAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	beforeCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	s.Require().NoError(err)
	beforeCmdxBalanceOfOwner, err := s.getBalance(beforeAuction.VaultOwner.String(), "ucmdx")
	s.Require().NoError(err)
	beforeCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)

	_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx),
		&auctionTypes.MsgPlaceDutchBidRequest{
			AuctionId:        1,
			Bidder:           userAddress1,
			Amount:           beforeAuction.OutflowTokenCurrentAmount,
			AppId:            appId,
			AuctionMappingId: auctionMappingId,
		})
	s.Require().NoError(err)

	_, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().Error(err)
	afterCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	s.Require().NoError(err)
	afterCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)
	afterAuction, err := k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	userOutflowCoin := beforeAuction.InflowTokenTargetAmount.Sub(beforeAuction.InflowTokenCurrentAmount)
	userInflowCoin := beforeAuction.OutflowTokenCurrentAmount.Sub(afterAuction.OutflowTokenCurrentAmount)
	s.Require().Equal(beforeAuction.OutflowTokenCurrentAmount.Sub(userInflowCoin), afterAuction.OutflowTokenCurrentAmount)
	s.Require().Equal(afterAuction.InflowTokenTargetAmount, afterAuction.InflowTokenCurrentAmount)
	getVaultOwnerBal, err := s.getBalance(string(afterAuction.VaultOwner.String()), "ucmdx")
	amtToOwner := getVaultOwnerBal.Sub((beforeCmdxBalanceOfOwner).Add(userInflowCoin))
	s.Require().Equal(beforeCmdxBalance.Add(userInflowCoin).Add(amtToOwner), afterCmdxBalance)
	s.Require().Equal(beforeCmstBalance.Sub(userOutflowCoin), afterCmstBalance)
}

func (s *KeeperTestSuite) TestCloseDutchAuctionWithProtocolLoss() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	s.TestDutchBid()
	k, liquidationKeeper, ctx := &s.keeper, &s.liquidationKeeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	lockedVault, found := liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	err := k.DutchActivator(*ctx, lockedVault)
	s.Require().Error(err)
	server := auctionKeeper.NewMsgServiceServer(*k)
	beforeAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	beforeCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)
	s.advanceseconds(250)

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)

	err1 := k.FundModule(*ctx, auctionTypes.ModuleName, "ucmst", 10000000)
	s.Require().NoError(err1)
	err = s.app.BankKeeper.SendCoinsFromModuleToModule(*ctx, auctionTypes.ModuleName, collectorTypes.ModuleName, sdk.NewCoins(sdk.NewCoin("ucmst", sdk.NewInt(10000000))))
	s.Require().NoError(err)

	err = s.app.CollectorKeeper.SetNetFeeCollectedData(*ctx, 1, 2, sdk.NewInt(10000000))
	s.Require().NoError(err)

	_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx),
		&auctionTypes.MsgPlaceDutchBidRequest{
			AuctionId:        1,
			Bidder:           userAddress1,
			Amount:           beforeAuction.OutflowTokenCurrentAmount,
			AppId:            appId,
			AuctionMappingId: auctionMappingId,
		})
	s.Require().NoError(err)

	_, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().Error(err)

	afterCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)
	afterAuction, err := k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	userOutflowCoin := afterAuction.InflowTokenCurrentAmount.Sub(beforeAuction.InflowTokenCurrentAmount)
	s.Require().True(afterAuction.InflowTokenCurrentAmount.IsLT(afterAuction.InflowTokenTargetAmount))

	s.Require().Equal(beforeCmstBalance.Sub(userOutflowCoin), afterCmstBalance)

	// verify loss
	stats, found := k.GetProtocolStat(*ctx, appId, 2)
	s.Require().True(found)
	loss := sdk.NewDecFromInt(afterAuction.InflowTokenTargetAmount.Sub(afterAuction.InflowTokenCurrentAmount).Amount)
	s.Require().Equal(loss, stats.Loss)
}

func (s *KeeperTestSuite) TestRestartDutchAuction() {
	s.TestDutchBid()
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	dutchAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	// exact the auction duration
	s.advanceseconds(300)

	startPrice := dutchAuction.OutflowTokenCurrentPrice
	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)

	dutchAuction, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(dutchAuction.OutflowTokenCurrentPrice, startPrice.Mul(sdk.MustNewDecFromStr("0.6")))

	// full the auction duration RESTART
	s.advanceseconds(1)
	beforeAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
	s.Require().NoError(err)
	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
	s.Require().NoError(err)

	afterAuction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(beforeAuction.AuctionId, afterAuction.AuctionId)
	s.Require().Equal(beforeAuction.OutflowTokenInitAmount, afterAuction.OutflowTokenInitAmount)
	s.Require().Equal(beforeAuction.OutflowTokenCurrentAmount, afterAuction.OutflowTokenCurrentAmount)
	s.Require().Equal(beforeAuction.InflowTokenTargetAmount, afterAuction.InflowTokenTargetAmount)
	s.Require().Equal(beforeAuction.InflowTokenCurrentAmount, afterAuction.InflowTokenCurrentAmount)
	s.Require().Equal(beforeAuction.OutflowTokenInitialPrice, afterAuction.OutflowTokenInitialPrice)
	s.Require().Equal(beforeAuction.OutflowTokenEndPrice, afterAuction.OutflowTokenEndPrice)
	s.Require().Equal(beforeAuction.InflowTokenCurrentPrice, afterAuction.InflowTokenCurrentPrice)
	s.Require().Equal(beforeAuction.AuctionStatus, afterAuction.AuctionStatus)
	s.Require().Equal(beforeAuction.AuctionMappingId, afterAuction.AuctionMappingId)
	s.Require().Equal(beforeAuction.AppId, afterAuction.AppId)
	s.Require().Equal(beforeAuction.AssetInId, afterAuction.AssetInId)
	s.Require().Equal(beforeAuction.AssetOutId, afterAuction.AssetOutId)
	s.Require().Equal(beforeAuction.LockedVaultId, afterAuction.LockedVaultId)
	s.Require().Equal(beforeAuction.VaultOwner, afterAuction.VaultOwner)
	s.Require().Equal(beforeAuction.LiquidationPenalty, afterAuction.LiquidationPenalty)
	s.Require().Equal(afterAuction.OutflowTokenCurrentPrice, startPrice)

	// half the auction duration
	s.advanceseconds(150)

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)

	dutchAuction, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(dutchAuction.OutflowTokenCurrentPrice, startPrice.Mul(sdk.MustNewDecFromStr("0.8")))
}
