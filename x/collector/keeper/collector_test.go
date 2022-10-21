package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	auctionTypes "github.com/comdex-official/comdex/x/auction/types"
	collectorTypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddAppAsset() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
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
				Recipient:     userAddress,
			},
			{
				AssetId:       2,
				GenesisSupply: genesisSupply,
				IsGovToken:    true,
				Recipient:     userAddress,
			},
		},
	}
	// {
	// 	Name:             "commodo",
	// 	ShortName:        "commodo",
	// 	MinGovDeposit:    sdk.NewIntFromUint64(10000000),
	// 	GovTimeInSeconds: 900,
	// 	GenesisToken: []assetTypes.MintGenesisToken{
	// 		{
	// 			3,
	// 			genesisSupply,
	// 			true,
	// 			userAddress,
	// 		},
	// 	},
	// },
	err := assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := assetTypes.Asset{
		Name:      "CMDX",
		Denom:     "ucmdx",
		Decimals:  1000000,
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg2)

	msg3 := assetTypes.Asset{
		Name:      "CMST",
		Denom:     "ucmst",
		Decimals:  1000000,
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg3)

	msg4 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  1000000,
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg4)

	s.Require().NoError(err)
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
	s.auctionKeeper.SetAuctionParams(*ctx, auctionParams)
}

func (s *KeeperTestSuite) TestWasmUpdateCollectorLookupTable() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.TestWasmSetCollectorLookupTableAndAuctionControl()
	for _, tc := range []struct {
		name string
		msg  bindings.MsgUpdateCollectorLookupTable
	}{
		{
			"Wasm Update MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgUpdateCollectorLookupTable{
				AppID:            1,
				AssetID:          2,
				SurplusThreshold: 9999,
				DebtThreshold:    99,
				LSR:              sdk.MustNewDecFromStr("0.001"),
				LotSize:          100,
				BidFactor:        sdk.MustNewDecFromStr("0.00001"),
				DebtLotSize:      300,
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmUpdateCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID, tc.msg.AssetID)
			s.Require().True(found)
			s.Require().Equal(result.AppId, tc.msg.AppID)
			s.Require().Equal(result.CollectorAssetId, tc.msg.AssetID)
			s.Require().Equal(result.SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.LockerSavingRate, tc.msg.LSR)
			s.Require().Equal(result.LotSize, tc.msg.LotSize)
			s.Require().Equal(result.BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.DebtLotSize, tc.msg.DebtLotSize)
		})
	}
}

func (s *KeeperTestSuite) TestWasmSetCollectorLookupTableAndAuctionControl() {
	// userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.AddAppAsset()

	for index, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: 10000000,
				DebtThreshold:    5000000,
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          2000000,
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      2000000,
			},
		},
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 3",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 3,
				SecondaryAssetID: 2,
				SurplusThreshold: 10000000,
				DebtThreshold:    5000000,
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          2000000,
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      2000000,
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTableByApp(*ctx, tc.msg.AppID)
			s.Require().True(found)
			s.Require().Equal(result[index].AppId, tc.msg.AppID)
			s.Require().Equal(result[index].CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result[index].SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result[index].SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result[index].DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result[index].LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result[index].LotSize, tc.msg.LotSize)
			s.Require().Equal(result[index].BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result[index].DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	s.AddAuctionParams()
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 1 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    bool(true),
				IsDebtAuctions:       bool(false),
				IsDistributor:        bool(false),
				AssetOutOraclePrices: bool(false),
				AssetOutPrices:       uint64(1000000),
			},
		},
		{
			"Wasm Add Auction Control AppID 1 AssetID 3",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             uint64(3),
				IsSurplusAuctions:    bool(true),
				IsDebtAuctions:       bool(false),
				IsDistributor:        bool(false),
				AssetOutOraclePrices: bool(false),
				AssetOutPrices:       uint64(100000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetAuctionMappingForApp(*ctx, &tc.msg)
			s.Require().NoError(err)
			result1, found := collectorKeeper.GetAuctionMappingForApp(*ctx, tc.msg.AppID, tc.msg.AssetIDs)
			s.Require().True(found)
			s.Require().Equal(result1.AssetId, tc.msg.AssetIDs)
			s.Require().Equal(result1.IsSurplusAuction, tc.msg.IsSurplusAuctions)
			s.Require().Equal(result1.IsDebtAuction, tc.msg.IsDebtAuctions)
			s.Require().Equal(result1.IsDistributor, tc.msg.IsDistributor)
			s.Require().Equal(result1.IsAuctionActive, false)
			s.Require().Equal(result1.AssetOutOraclePrice, tc.msg.AssetOutOraclePrices)
			s.Require().Equal(result1.AssetOutPrice, tc.msg.AssetOutPrices)
		})
	}
}

func (s *KeeperTestSuite) TestSetNetFeesCollected() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.TestWasmSetCollectorLookupTableAndAuctionControl()
	negNumber, _ := sdk.NewIntFromString("-100")
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		fee           sdk.Int
		errorExpected bool
	}{
		{
			"Set net fees collected : AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(100),
			false,
		},
		{
			"Set net fees collected : cannot be negative AppID 1 AssetID 2",
			1,
			2,
			negNumber,
			true,
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.SetNetFeeCollectedData(*ctx, tc.appID, tc.assetID, tc.fee)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				netFeesData, found := collectorKeeper.GetNetFeeCollectedData(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(tc.appID, netFeesData.AppId)
				s.Require().Equal(tc.assetID, netFeesData.AssetId)
				s.Require().Equal(tc.fee, netFeesData.NetFeesCollected)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddNetFeesCollected() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.TestSetNetFeesCollected()
	negNumber, _ := sdk.NewIntFromString("-100")
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		fee           sdk.Int
		errorExpected bool
	}{
		{
			"Add net fees collected : AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(974),
			false,
		},
		{
			"Add net fees collected : AppID 1 AssetID 2",
			1,
			2,
			negNumber,
			true,
		},
	} {
		s.Run(tc.name, func() {
			netFeesData1, found := collectorKeeper.GetNetFeeCollectedData(*ctx, tc.appID, tc.assetID)
			s.Require().True(found)
			err := collectorKeeper.SetNetFeeCollectedData(*ctx, tc.appID, tc.assetID, tc.fee)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				netFeesData, found := collectorKeeper.GetNetFeeCollectedData(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(tc.appID, netFeesData.AppId)
				s.Require().Equal(tc.assetID, netFeesData.AssetId)
				s.Require().Equal(netFeesData1.NetFeesCollected.Add(tc.fee), netFeesData.NetFeesCollected)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDecreaseNetFeesCollected() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.TestSetNetFeesCollected()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		fee           sdk.Int
		errorExpected bool
	}{
		{
			"Decrease net fees collected : AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(52),
			false,
		},
		{
			"Decrease net fees collected : Net fees cannot be negative AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(102),
			true,
		},
	} {
		s.Run(tc.name, func() {
			netFeesData1, found := collectorKeeper.GetNetFeeCollectedData(*ctx, tc.appID, tc.assetID)
			s.Require().True(found)
			err := collectorKeeper.DecreaseNetFeeCollectedData(*ctx, tc.appID, tc.assetID, tc.fee, netFeesData1)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				netFeesData, found := collectorKeeper.GetNetFeeCollectedData(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(tc.appID, netFeesData.AppId)
				s.Require().Equal(tc.assetID, netFeesData.AssetId)
				s.Require().Equal(netFeesData1.NetFeesCollected.Sub(tc.fee), netFeesData.NetFeesCollected)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGetAmountFromCollector() {
	collectorKeeper, auctionKeeper, ctx := &s.collectorKeeper, &s.auctionKeeper, &s.ctx
	s.TestSetNetFeesCollected()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		GetAmount     sdk.Int
		FundAmount    uint64
		denom         string
		errorExpected bool
	}{
		{
			"Get Amount From Collector : AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(52),
			100,
			"ucmst",
			false,
		},
		{
			"Get Amount From Collector : Insufficient Balance AppID 1 AssetID 2",
			1,
			2,
			sdk.NewIntFromUint64(101),
			100,
			"ucmst",
			true,
		},
	} {
		s.Run(tc.name, func() {
			err := auctionKeeper.FundModule(*ctx, "auctionV1", tc.denom, tc.FundAmount)
			s.Require().NoError(err)
			err = s.app.BankKeeper.SendCoinsFromModuleToModule(*ctx, "auctionV1", "collectorV1", sdk.NewCoins(sdk.NewCoin(tc.denom, sdk.NewIntFromUint64(tc.FundAmount))))
			s.Require().NoError(err)
			beforeCollectorBalance := auctionKeeper.GetModuleAccountBalance(*ctx, "collectorV1", tc.denom)
			returnAmount, err := collectorKeeper.GetAmountFromCollector(*ctx, tc.appID, tc.assetID, tc.GetAmount)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.GetAmount, returnAmount)
				auctionBalance := auctionKeeper.GetModuleAccountBalance(*ctx, "auctionV1", tc.denom)
				s.Require().Equal(tc.GetAmount, auctionBalance)
				currentCollectorBalance := auctionKeeper.GetModuleAccountBalance(*ctx, "collectorV1", tc.denom)
				s.Require().Equal(currentCollectorBalance, beforeCollectorBalance.Sub(tc.GetAmount))
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateCollector() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.AddAppAsset()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		collectorData collectorTypes.CollectorData
		errorExpected bool
	}{
		{
			name:    "UpdateCollector : AppID 1 AssetID 2",
			appID:   1,
			assetID: 2,
			collectorData: collectorTypes.CollectorData{
				CollectedStabilityFee:       sdk.NewIntFromUint64(100),
				CollectedClosingFee:         sdk.NewIntFromUint64(200),
				CollectedOpeningFee:         sdk.NewIntFromUint64(300),
				LiquidationRewardsCollected: sdk.NewIntFromUint64(400),
			},
			errorExpected: false,
		},
		{
			name:    "UpdateCollector : AppID 1 AssetID 3",
			appID:   1,
			assetID: 3,
			collectorData: collectorTypes.CollectorData{
				CollectedStabilityFee:       sdk.NewIntFromUint64(100),
				CollectedClosingFee:         sdk.NewIntFromUint64(200),
				CollectedOpeningFee:         sdk.NewIntFromUint64(300),
				LiquidationRewardsCollected: sdk.NewIntFromUint64(500),
			},
			errorExpected: false,
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.UpdateCollector(*ctx,
				tc.appID,
				tc.assetID,
				tc.collectorData.CollectedStabilityFee,
				tc.collectorData.CollectedClosingFee,
				tc.collectorData.CollectedOpeningFee,
				tc.collectorData.LiquidationRewardsCollected)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				collectorData, found := collectorKeeper.GetCollectorDataForAppIDAssetID(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(tc.collectorData, collectorData)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddUpdateCollector() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx
	s.TestUpdateCollector()
	for _, tc := range []struct {
		name          string
		appID         uint64
		assetID       uint64
		collectorData collectorTypes.CollectorData
		errorExpected bool
	}{
		{
			name:    "Add UpdateCollector : AppID 1 AssetID 2",
			appID:   1,
			assetID: 2,
			collectorData: collectorTypes.CollectorData{
				CollectedStabilityFee:       sdk.NewIntFromUint64(100),
				CollectedClosingFee:         sdk.NewIntFromUint64(200),
				CollectedOpeningFee:         sdk.NewIntFromUint64(300),
				LiquidationRewardsCollected: sdk.NewIntFromUint64(400),
			},
			errorExpected: false,
		},
	} {
		s.Run(tc.name, func() {
			beforeCollectorData, found := collectorKeeper.GetCollectorDataForAppIDAssetID(*ctx, tc.appID, tc.assetID)
			s.Require().True(found)
			err := collectorKeeper.UpdateCollector(*ctx,
				tc.appID,
				tc.assetID,
				tc.collectorData.CollectedStabilityFee,
				tc.collectorData.CollectedClosingFee,
				tc.collectorData.CollectedOpeningFee,
				tc.collectorData.LiquidationRewardsCollected)
			if tc.errorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				currentCollectorData, found := collectorKeeper.GetCollectorDataForAppIDAssetID(*ctx, tc.appID, tc.assetID)
				s.Require().True(found)
				s.Require().Equal(beforeCollectorData.CollectedClosingFee.Add(tc.collectorData.CollectedClosingFee), currentCollectorData.CollectedClosingFee)
				s.Require().Equal(beforeCollectorData.CollectedStabilityFee.Add(tc.collectorData.CollectedStabilityFee), currentCollectorData.CollectedStabilityFee)
				s.Require().Equal(beforeCollectorData.CollectedOpeningFee.Add(tc.collectorData.CollectedOpeningFee), currentCollectorData.CollectedOpeningFee)
				s.Require().Equal(beforeCollectorData.LiquidationRewardsCollected.Add(tc.collectorData.LiquidationRewardsCollected), currentCollectorData.LiquidationRewardsCollected)

			}
		})
	}
}
