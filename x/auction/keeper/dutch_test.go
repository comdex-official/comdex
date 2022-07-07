package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctionTypes "github.com/comdex-official/comdex/x/auction/types"
	liquidationTypes "github.com/comdex-official/comdex/x/liquidation/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper1 "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (s *KeeperTestSuite) AddPairAndExtendedPairVault1() {

	assetKeeper, liquidationKeeper, ctx := &s.assetKeeper, &s.liquidationKeeper, &s.ctx

	for _, tc := range []struct {
		name              string
		pair              assetTypes.Pair
		extendedPairVault bindings.MsgAddExtendedPairsVault
		symbol1           string
		symbol2           string
	}{
		{"Add Pair , Extended Pair Vault : cmdx cmst",
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
				DebtCeiling:         1000000000000,
				DebtFloor:           1000000,
				IsStableMintVault:   false,
				MinCr:               sdk.MustNewDecFromStr("1.5"),
				PairName:            "CMDX-B",
				AssetOutOraclePrice: true,
				AssetOutPrice:       1000000,
				MinUsdValueLeft:     1000000,
			},
			"ucmdx",
			"ucmst",
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddPairsRecords(*ctx, tc.pair)
			s.Require().NoError(err)

			err = assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.extendedPairVault)
			s.Require().NoError(err)

			err = liquidationKeeper.WhitelistAppID(*ctx, 1)
			s.Require().NoError(err)

			s.SetInitialOraclePriceForSymbols(tc.symbol1, tc.symbol2)
		})
	}
}

func (s *KeeperTestSuite) SetOraclePrice(symbol string, price uint64) {
	var (
		store = s.app.MarketKeeper.Store(s.ctx)
		key   = markettypes.PriceForMarketKey(symbol)
	)
	value := s.app.AppCodec().MustMarshal(
		&protobuftypes.UInt64Value{
			Value: price,
		},
	)
	store.Set(key, value)
}

func (s *KeeperTestSuite) SetInitialOraclePriceForSymbols(asset1 string, asset2 string) {
	s.SetOraclePrice(asset1, 2000000)
	s.SetOraclePrice(asset2, 1000000)
}
func (s *KeeperTestSuite) ChangeOraclePrice(asset string) {
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
		{"Create Vault : AppID 1 extended pair 1 user address 1",
			vaultTypes.MsgCreateRequest{
				From:                userAddress1,
				AppId:               1,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
		{"Create Vault : AppID 1 extended pair 1 user address 2",
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

func (s *KeeperTestSuite) GetVaultCountForExtendedPairIDbyAppID(appID uint64) int {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	res, found := liquidationKeeper.GetAppExtendedPairVaultMapping(*ctx, appID)
	s.Require().True(found)
	return len(res.ExtendedPairVaults[0].VaultIds)
}

func (s *KeeperTestSuite) GetAssetPrice(id uint64) sdk.Dec {
	marketKeeper, ctx := &s.marketKeeper, &s.ctx
	price, found := marketKeeper.GetPriceForAsset(*ctx, id)
	s.Require().True(found)
	price1 := sdk.NewDecFromInt(sdk.NewIntFromUint64(price))
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
	msg1 := []assetTypes.AppData{{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				3,
				&genesisSupply,
				true,
				userAddress1,
			},
			{
				2,
				&genesisSupply,
				true,
				userAddress1,
			},
		},
	},
		{
			Name:             "commodo",
			ShortName:        "commodo",
			MinGovDeposit:    sdk.NewIntFromUint64(10000000),
			GovTimeInSeconds: 900,
			GenesisToken: []assetTypes.MintGenesisToken{
				{
					3,
					&genesisSupply,
					true,
					userAddress1,
				},
			},
		},
	}
	err = assetKeeper.AddAppRecords(*ctx, msg1...)
	s.Require().NoError(err)

	for index, tc := range []struct {
		name string
		msg  assetTypes.Asset
	}{
		{"Add Asset 1",
			assetTypes.Asset{Name: "CMDX",
				Denom:     "ucmdx",
				Decimals:  1000000,
				IsOnChain: true},
		},
		{"Add Asset 2",
			assetTypes.Asset{Name: "CMST",
				Denom:     "ucmst",
				Decimals:  1000000,
				IsOnChain: true},
		},
		{"Add Asset 3",
			assetTypes.Asset{Name: "HARBOR",
				Denom:     "uharbor",
				Decimals:  1000000,
				IsOnChain: true},
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAssetRecords(*ctx, tc.msg)
			s.Require().NoError(err)
			s.marketKeeper.SetMarketForAsset(*ctx, uint64(index+1), tc.msg.Denom)
			market := markettypes.Market{
				Symbol:   tc.msg.Denom,
				ScriptID: 12,
				Rates:    1000000,
			}
			s.app.MarketKeeper.SetMarket(s.ctx, market)
			res := s.app.MarketKeeper.HasMarketForAsset(s.ctx, uint64(index+1))
			s.Require().True(res)
			s.fundAddr(addr, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
			s.fundAddr(addr2, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
		})
	}

}

func (s *KeeperTestSuite) LiquidateVaults1() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.CreateVault()
	currentVaultsCount := 2
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1)), currentVaultsCount)
	beforeVault, found := liquidationKeeper.GetVault(*ctx, "cswap1")
	s.Require().True(found)

	// Liquidation shouldn't happen as price not changed
	err := liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(0))

	// Liquidation should happen as price changed
	s.ChangeOraclePrice("ucmdx")
	err = liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	err = liquidationKeeper.UpdateLockedVaults(*ctx)
	s.Require().NoError(err)
	id = liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(2))
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount-2)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1)), currentVaultsCount-2)

	lockedVault := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, beforeVault.Id)
	s.Require().Equal(lockedVault[0].ExtendedPairId, beforeVault.ExtendedPairVaultID)
	s.Require().Equal(lockedVault[0].Owner, beforeVault.Owner)
	s.Require().Equal(lockedVault[0].AmountIn, beforeVault.AmountIn)
	s.Require().Equal(lockedVault[0].AmountOut, beforeVault.AmountOut)
	s.Require().Equal(lockedVault[0].UpdatedAmountOut, beforeVault.AmountOut.Add(beforeVault.InterestAccumulated).Add(beforeVault.ClosingFeeAccumulated))
	s.Require().Equal(lockedVault[0].Initiator, liquidationTypes.ModuleName)
	s.Require().Equal(lockedVault[0].IsAuctionInProgress, false)
	s.Require().Equal(lockedVault[0].IsAuctionComplete, false)
	s.Require().Equal(lockedVault[0].SellOffHistory, []string(nil))
	price, found := s.app.MarketKeeper.GetPriceForAsset(*ctx, uint64(1))
	s.Require().True(found)
	s.Require().Equal(lockedVault[0].CollateralToBeAuctioned, beforeVault.AmountIn.ToDec().Mul(sdk.NewIntFromUint64(price).ToDec()))
	s.Require().Equal(lockedVault[0].CurrentCollaterlisationRatio, lockedVault[0].AmountIn.ToDec().Mul(s.GetAssetPrice(1)).Quo(lockedVault[0].UpdatedAmountOut.ToDec().Mul(s.GetAssetPrice(2))))
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

	err = k.DutchActivator(*ctx)
	s.Require().NoError(err)

	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	auction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	lockedVault, found := liquidationKeeper.GetLockedVault(*ctx, 1)
	s.Require().True(found)

	s.Require().Equal(auction.AppId, lockedVault.AppId)
	s.Require().Equal(auction.AuctionId, auctionId)
	s.Require().Equal(auction.AuctionMappingId, auctionMappingId)
	s.Require().Equal(auction.OutflowTokenInitAmount.Amount, lockedVault.AmountIn)
	s.Require().Equal(auction.OutflowTokenCurrentAmount.Amount, lockedVault.AmountIn)
	s.Require().Equal(auction.InflowTokenCurrentAmount.Amount, sdk.ZeroInt())

	inFlowTokenTargetAmount := lockedVault.AmountOut
	mulfactor := inFlowTokenTargetAmount.ToDec().Mul(auction.LiquidationPenalty)
	inFlowTokenTargetAmount = inFlowTokenTargetAmount.Add(mulfactor.TruncateInt())

	s.Require().Equal(auction.InflowTokenTargetAmount.Amount, inFlowTokenTargetAmount)

	s.Require().Equal(auction.VaultOwner, addr1)
	s.Require().Equal(auction.AuctionStatus, auctionTypes.AuctionStartNoBids)

	assetOutPrice, found := s.marketKeeper.GetPriceForAsset(*ctx, auction.AssetOutId)
	s.Require().True(found)
	assetInPrice, found := s.marketKeeper.GetPriceForAsset(*ctx, auction.AssetInId)
	s.Require().True(found)
	s.Require().Equal(auction.OutflowTokenCurrentPrice, sdk.NewDecFromInt(sdk.NewIntFromUint64(assetOutPrice)).Mul(sdk.MustNewDecFromStr("1.2")))
	s.Require().Equal(auction.OutflowTokenEndPrice, auction.OutflowTokenInitialPrice.Mul(sdk.MustNewDecFromStr("0.6")))
	s.Require().Equal(auction.InflowTokenCurrentPrice, sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice)))
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
	beforeAuction, err := k.GetDutchAuction(*ctx, appID, auctionMappingID, auctionID)
	s.Require().NoError(err)
	beforeCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	s.Require().NoError(err)
	beforeCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)

	//expect error as max price is less than current price of outflow token
	_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx),
		&auctionTypes.MsgPlaceDutchBidRequest{
			AuctionId:        1,
			Bidder:           userAddress1,
			Amount:           ParseCoin("1000ucmdx"),
			Max:              sdk.MustNewDecFromStr("1.1"),
			AppId:            appID,
			AuctionMappingId: auctionMappingID,
		})
	s.Require().Error(err)

	//dont expect error
	_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx),
		&auctionTypes.MsgPlaceDutchBidRequest{
			AuctionId:        1,
			Bidder:           userAddress1,
			Amount:           ParseCoin("1000ucmdx"),
			Max:              sdk.MustNewDecFromStr("1.4"),
			AppId:            appID,
			AuctionMappingId: auctionMappingID,
		})
	s.Require().NoError(err)

	afterCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	s.Require().NoError(err)
	afterCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)

	afterAuction, err := k.GetDutchAuction(*ctx, appID, auctionMappingID, auctionID)
	s.Require().NoError(err)

	userBid, err := k.GetDutchUserBidding(*ctx, userAddress1, appID, biddingID)
	userOutflowCoin := ParseCoin("1200ucmst")
	userInflowCoin := ParseCoin("1000ucmdx")
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
	//beforeCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	//s.Require().NoError(err)
	beforeCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)

	_, err = server.MsgPlaceDutchBid(sdk.WrapSDKContext(*ctx),
		&auctionTypes.MsgPlaceDutchBidRequest{
			AuctionId:        1,
			Bidder:           userAddress1,
			Amount:           beforeAuction.OutflowTokenCurrentAmount,
			Max:              sdk.MustNewDecFromStr("1.2"),
			AppId:            appId,
			AuctionMappingId: auctionMappingId,
		})
	s.Require().NoError(err)

	_, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().Error(err)
	//afterCmdxBalance, err := s.getBalance(userAddress1, "ucmdx")
	//s.Require().NoError(err)
	afterCmstBalance, err := s.getBalance(userAddress1, "ucmst")
	s.Require().NoError(err)
	afterAuction, err := k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	userOutflowCoin := beforeAuction.InflowTokenTargetAmount.Sub(beforeAuction.InflowTokenCurrentAmount)
	//userInflowCoin := beforeAuction.OutflowTokenCurrentAmount
	//s.Require().Equal(beforeAuction.OutflowTokenCurrentAmount.Sub(userInflowCoin), afterAuction.OutflowTokenCurrentAmount)
	s.Require().Equal(afterAuction.InflowTokenTargetAmount, afterAuction.InflowTokenCurrentAmount)
	//s.Require().Equal(beforeCmdxBalance.Add(userInflowCoin), afterCmdxBalance)
	s.Require().Equal(beforeCmstBalance.Sub(userOutflowCoin), afterCmstBalance)

}

func (s *KeeperTestSuite) TestRestartDutchAuction() {
	//userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	s.TestDutchBid()
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	//server := auctionKeeper.NewMsgServiceServer(*k)
	_, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.advanceseconds(200)

	err = k.DutchActivator(*ctx)
	s.Require().NoError(err)

	_, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

}

//
//func (s *KeeperTestSuite) TestDecToInt() {
//	num := sdk.MustNewDecFromStr("133333333.999999")
//	fmt.Println(num)
//	res := num.TruncateInt()
//	fmt.Println(res)
//}

//
//func (s *KeeperTestSuite) TestLinearPriceFunction() {
//	k := &s.keeper
//
//	top := sdk.MustNewDecFromStr("100").Mul(sdk.MustNewDecFromStr("1"))
//	tau := sdk.NewIntFromUint64(60)
//	dur := sdk.NewInt(0)
//	for n := 0; n <= 60; n++ {
//		fmt.Println("top tau dur seconds")
//		fmt.Println(top, tau, dur)
//		fmt.Println("price")
//		price := k.GetPriceFromLinearDecreaseFunction(top, tau, dur)
//		fmt.Println(price)
//		dur = dur.Add(sdk.NewInt(1))
//
//	}
//}

//
//func (s *KeeperTestSuite) BidDutchAuction() {
//	k, ctx := &s.keeper, &s.ctx
//
//	max := sdk.MustNewDecFromStr("12")
//	//create auction
//	s.TestDutchActivator()
//	userTokens := ParseCoin("2500denom2")
//
//	bid := ParseCoin("5denom1")
//	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
//	s.Require().NoError(err)
//
//	//fund bidder
//	s.fundAddr(bidder, userTokens)
//
//	//place bid
//	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
//	s.Require().NoError(err)
//	_, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	//check if user balance reduced
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
//	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
//
//	_, err = k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//
//	//close auction by advancing time
//	//s.advanceseconds(advanceSeconds)
//	//k.CloseDutchAuction(*ctx, auction)
//
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom1"))
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
//	//check if user got collateral
//	s.Require().Equal(sdk.NewInt(10), k.GetBalance(*ctx, bidder, "denom1").Amount)
//
//	//check status of bid
//	_, err = k.GetHistoryDutchUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//
//	//get closed auction
//	_, err = k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//}
//
//func (s *KeeperTestSuite) BidsDutchAuction() {
//	k, ctx := &s.keeper, &s.ctx
//	appId := uint64(1)
//	auctionMappingId := uint64(3)
//	auctionId := uint64(1)
//	max := sdk.MustNewDecFromStr("12")
//	//create auction
//	s.TestCreateDutchAuction()
//	userTokens := ParseCoin("1000denom2")
//
//	bid := ParseCoin("10denom1")
//	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
//	s.Require().NoError(err)
//
//	//fund bidder
//	s.fundAddr(bidder, userTokens)
//
//	//place bid
//	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
//	s.Require().NoError(err)
//	auction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	//check if user balance reduced
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
//	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
//
//	_, err = k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//
//	//place bid
//	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
//	s.Require().NoError(err)
//	auction, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	//check if user balance reduced
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
//	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
//
//	_, err = k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//
//	//close auction by advancing time
//	s.advanceseconds(advanceSeconds)
//	err = k.CloseDutchAuction(*ctx, auction)
//	if err != nil {
//		return
//	}
//
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom1"))
//	//check if user got collateral
//	s.Require().Equal(sdk.NewInt(10), k.GetBalance(*ctx, bidder, "denom1").Amount)
//
//	//check status of bid
//	_, err = k.GetHistoryDutchUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//
//	//get closed auction
//	_, err = k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//}
//
//func getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Int {
//	result1 := tau.Sub(dur)
//	result2 := top.MulInt(result1)
//	result3 := result2.Quo(tau.ToDec())
//	return result3.TruncateInt()
//}
//func (s *KeeperTestSuite) LinearPriceFunction() {
//	top := sdk.MustNewDecFromStr("100").Mul(sdk.MustNewDecFromStr("1"))
//	tau := sdk.NewIntFromUint64(60)
//	dur := sdk.NewInt(0)
//	for n := 0; n <= 60; n++ {
//		fmt.Println("top tau dur seconds")
//		fmt.Println(top, tau, dur)
//		fmt.Println("price")
//		price := getPriceFromLinearDecreaseFunction(top, tau, dur)
//		fmt.Println(price)
//		dur = dur.Add(sdk.NewInt(1))
//
//	}
//}
