package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/comdex-official/comdex/app/wasm/bindings"
)

func (s *KeeperTestSuite) WasmSetCollectorLookupTableAndAuctionControlForSurplus() {
	// userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 2 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            2,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdkmath.NewInt(10000000),
				DebtThreshold:    sdkmath.NewInt(5000000),
				LockerSavingRate: sdkmath.LegacyMustNewDecFromStr("0.1"),
				LotSize:          sdkmath.NewInt(200000),
				BidFactor:        sdkmath.LegacyMustNewDecFromStr("0.01"),
				DebtLotSize:      sdkmath.NewInt(2000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID)
			s.Require().True(found)
			s.Require().Equal(result.AppId, tc.msg.AppID)
			s.Require().Equal(result.CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result.SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result.SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result.LotSize, tc.msg.LotSize)
			s.Require().Equal(result.BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	// s.AddAuctionParams()
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 2 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                2,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    true,
				IsDebtAuctions:       false,
				IsDistributor:        false,
				AssetOutOraclePrices: false,
				AssetOutPrices:       uint64(1000000),
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

func (s *KeeperTestSuite) WasmSetCollectorLookupTableAndAuctionControlForDebt() {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 2 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            2,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdkmath.NewInt(10000000),
				DebtThreshold:    sdkmath.NewInt(5000000),
				LockerSavingRate: sdkmath.LegacyMustNewDecFromStr("0.1"),
				LotSize:          sdkmath.NewInt(200000),
				BidFactor:        sdkmath.LegacyMustNewDecFromStr("0.01"),
				DebtLotSize:      sdkmath.NewInt(2000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID)
			s.Require().True(found)
			s.Require().Equal(result.AppId, tc.msg.AppID)
			s.Require().Equal(result.CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result.SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result.SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result.LotSize, tc.msg.LotSize)
			s.Require().Equal(result.BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 2 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                2,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    false,
				IsDebtAuctions:       true,
				IsDistributor:        false,
				AssetOutOraclePrices: false,
				AssetOutPrices:       uint64(1000000),
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
