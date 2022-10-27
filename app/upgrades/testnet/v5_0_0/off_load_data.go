package v5_0_0

import (
	"fmt"

	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"

	// assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migration Plan
// First Step: Migrate all pools to new updated proto and update the poolID counter
// Second Step: Migrate all lend positions, add up total lend by asset and update it in new proto, also append lendID in new struct and Update counter to last lend position
// Third Step: Migrate all borrow to new proto and make is liquidated flag false
// Fourth Step: Migrate all liquidated borrow to new borrow struct and make is_liquidated field to true
// Fifth Step: Start Auction for liquidated Borrow Positions
// Sixth Step: Correct AppID in Auction Params
// Seventh Step: In App Proposal migrate data to new type

// TODO: while testing revert back kv stores for pair, Asset_rates_stats & asset_pair mapping, Also check all queries

func MigrateData(ctx sdk.Context, k lendkeeper.Keeper) error {
	// err := FuncMigrateApp(ctx, k)
	// if err != nil {
	// 	return err
	// }

	err := FuncMigratePool(ctx, k)
	if err != nil {
		return err
	}

	err = FuncMigrateLend(ctx, k)
	if err != nil {
		return err
	}

	err = FuncMigrateBorrow(ctx, k)
	if err != nil {
		return err
	}

	err = FuncMigrateAuctionParams(ctx, k)
	if err != nil {
		return err
	}

	err = FuncMigrateLiquidatedBorrow(ctx, k)
	if err != nil {
		return err
	}

	return nil
}

// FuncMigratePool - First Step: Migrate all pools to new updated proto
func FuncMigratePool(ctx sdk.Context, k lendkeeper.Keeper) error {
	oldPools := k.OldGetPools(ctx)
	var (
		assetDataPoolOne []*types.AssetDataPoolMapping
		assetDataPoolTwo []*types.AssetDataPoolMapping
		assetData        []*types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := &types.AssetDataPoolMapping{
		AssetID:          1,
		AssetTransitType: 3,
		SupplyCap:        uint64(5000000000000000000),
	}
	assetDataPoolOneAssetTwo := &types.AssetDataPoolMapping{
		AssetID:          2,
		AssetTransitType: 1,
		SupplyCap:        uint64(1000000000000000000),
	}
	assetDataPoolOneAssetThree := &types.AssetDataPoolMapping{
		AssetID:          3,
		AssetTransitType: 2,
		SupplyCap:        uint64(5000000000000000000),
	}
	assetDataPoolTwoAssetFour := &types.AssetDataPoolMapping{
		AssetID:          4,
		AssetTransitType: 1,
		SupplyCap:        uint64(3000000000000000000),
	}
	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolTwo, assetDataPoolTwoAssetFour, assetDataPoolOneAssetOne, assetDataPoolOneAssetThree)

	for _, j := range oldPools {
		if j.PoolID == 1 {
			assetData = assetDataPoolOne
		} else {
			assetData = assetDataPoolTwo
		}
		newPool := types.Pool{
			PoolID:       j.PoolID,
			ModuleName:   j.ModuleName,
			CPoolName:    j.CPoolName,
			ReserveFunds: j.ReserveFunds,
			AssetData:    assetData,
		}

		for _, v := range newPool.AssetData {
			var assetStats types.PoolAssetLBMapping
			assetStats.PoolID = newPool.PoolID
			assetStats.AssetID = v.AssetID
			assetStats.TotalBorrowed = sdk.ZeroInt()
			assetStats.TotalStableBorrowed = sdk.ZeroInt()
			assetStats.TotalLend = sdk.ZeroInt()
			assetStats.TotalInterestAccumulated = sdk.ZeroInt()
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
			k.UpdateAPR(ctx, newPool.PoolID, v.AssetID)
			reserveBuybackStats, found := k.GetReserveBuybackAssetData(ctx, v.AssetID)
			if !found {
				reserveBuybackStats.AssetID = v.AssetID
				reserveBuybackStats.ReserveAmount = sdk.ZeroInt()
				reserveBuybackStats.BuybackAmount = sdk.ZeroInt()
				k.SetReserveBuybackAssetData(ctx, reserveBuybackStats)
			}
		}

		k.SetPool(ctx, newPool)
		k.SetPoolID(ctx, newPool.PoolID)
	}
	return nil
}

// FuncMigrateLend - Second Step: Migrate all lend positions, add up total lend by asset and update it in new proto, also append lendID in new struct
func FuncMigrateLend(ctx sdk.Context, k lendkeeper.Keeper) error {
	allLends := k.OldGetAllLend(ctx)
	for _, v := range allLends {
		if v.AmountIn.Amount.LTE(sdk.ZeroInt()) || v.AvailableToBorrow.LT(sdk.ZeroInt()) {
			continue
		}
		newLend := types.LendAsset{
			ID:                  v.ID,
			AssetID:             v.AssetID,
			PoolID:              v.PoolID,
			Owner:               v.Owner,
			AmountIn:            v.AmountIn,
			LendingTime:         v.LendingTime,
			AvailableToBorrow:   v.AvailableToBorrow,
			AppID:               v.AppID,
			GlobalIndex:         v.GlobalIndex,
			LastInteractionTime: v.LastInteractionTime,
			CPoolName:           v.CPoolName,
		}
		k.UpdateLendStats(ctx, v.AssetID, v.PoolID, v.AmountIn.Amount, true) // update global lend data in poolAssetLBMappingData
		k.SetUserLendIDCounter(ctx, newLend.ID)
		k.SetLend(ctx, newLend)

		// making UserAssetLendBorrowMapping for user
		var mappingData types.UserAssetLendBorrowMapping
		mappingData.Owner = newLend.Owner
		mappingData.LendId = newLend.ID
		mappingData.PoolId = v.PoolID
		mappingData.BorrowId = nil
		k.SetUserLendBorrowMapping(ctx, mappingData)

		// Adding Lend ID mapping to poolAssetLBMappingData
		poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, v.PoolID, v.AssetID)
		poolAssetLBMappingData.LendIds = append(poolAssetLBMappingData.LendIds, newLend.ID)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
	}

	return nil
}

// FuncMigrateBorrow - Third Step: Migrate all borrow to new proto and make is liquidated flag false
func FuncMigrateBorrow(ctx sdk.Context, k lendkeeper.Keeper) error {
	oldBorrows := k.OldGetAllBorrow(ctx)
	for _, v := range oldBorrows {
		if v.AmountIn.Amount.LTE(sdk.ZeroInt()) || v.AmountOut.Amount.LT(sdk.ZeroInt()) {
			continue
		}
		newBorrow := types.BorrowAsset{
			ID:                  v.ID,
			LendingID:           v.LendingID,
			IsStableBorrow:      v.IsStableBorrow,
			PairID:              v.PairID,
			AmountIn:            v.AmountIn,
			AmountOut:           v.AmountOut,
			BridgedAssetAmount:  v.BridgedAssetAmount,
			BorrowingTime:       v.BorrowingTime,
			StableBorrowRate:    v.StableBorrowRate,
			InterestAccumulated: sdk.NewDecFromInt(v.Interest_Accumulated),
			GlobalIndex:         v.GlobalIndex,
			ReserveGlobalIndex:  v.ReserveGlobalIndex,
			LastInteractionTime: v.LastInteractionTime,
			CPoolName:           v.CPoolName,
			IsLiquidated:        false,
		}
		lend, _ := k.GetLend(ctx, v.LendingID)
		pair, _ := k.GetLendPair(ctx, v.PairID)
		k.UpdateBorrowStats(ctx, pair, newBorrow.IsStableBorrow, v.AmountOut.Amount, true)

		poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
		poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, newBorrow.ID)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

		k.SetUserBorrowIDCounter(ctx, newBorrow.ID)
		k.SetBorrow(ctx, newBorrow)

		mappingData, _ := k.GetUserLendBorrowMapping(ctx, lend.Owner, v.LendingID)
		mappingData.BorrowId = append(mappingData.BorrowId, newBorrow.ID)
		k.SetUserLendBorrowMapping(ctx, mappingData)
	}
	return nil
}

// FuncMigrateLiquidatedBorrow - Fourth Step & Fifth Step : Migrate all liquidated borrow to new borrow struct and make is_liquidated field to true
func FuncMigrateLiquidatedBorrow(ctx sdk.Context, k lendkeeper.Keeper) error {
	liqBorrow := k.Liquidation.GetLockedVaultByApp(ctx, 3)
	for _, v := range liqBorrow {
		if v.AmountIn.LTE(sdk.ZeroInt()) || v.AmountOut.LT(sdk.ZeroInt()) {
			continue
		}
		borrowMetaData := v.GetBorrowMetaData()
		pair, _ := k.GetLendPair(ctx, v.ExtendedPairId)
		assetIn, _ := k.Asset.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.Asset.GetAsset(ctx, pair.AssetOut)
		amountIn := sdk.NewCoin(assetIn.Denom, v.AmountIn)
		amountOut := sdk.NewCoin(assetOut.Denom, v.AmountOut)
		pool, _ := k.GetPool(ctx, pair.AssetOutPoolID)
		assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
		reserveGlobalIndex, _ := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)

		newBorrow := types.BorrowAsset{
			ID:                  v.OriginalVaultId,
			LendingID:           borrowMetaData.LendingId,
			IsStableBorrow:      borrowMetaData.IsStableBorrow,
			PairID:              v.ExtendedPairId,
			AmountIn:            amountIn,
			AmountOut:           amountOut,
			BridgedAssetAmount:  borrowMetaData.BridgedAssetAmount,
			BorrowingTime:       ctx.BlockTime(),
			StableBorrowRate:    borrowMetaData.StableBorrowRate,
			InterestAccumulated: sdk.NewDecFromInt(v.InterestAccumulated),
			GlobalIndex:         assetStats.BorrowApr,
			ReserveGlobalIndex:  reserveGlobalIndex,
			LastInteractionTime: ctx.BlockTime(),
			CPoolName:           pool.CPoolName,
			IsLiquidated:        true,
		}
		lend, _ := k.GetLend(ctx, newBorrow.LendingID)
		k.UpdateBorrowStats(ctx, pair, newBorrow.IsStableBorrow, v.AmountOut, true)

		poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
		poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, newBorrow.ID)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

		k.SetUserBorrowIDCounter(ctx, newBorrow.ID)
		k.SetBorrow(ctx, newBorrow)

		mappingData, _ := k.GetUserLendBorrowMapping(ctx, lend.Owner, newBorrow.LendingID)
		mappingData.BorrowId = append(mappingData.BorrowId, newBorrow.ID)
		k.SetUserLendBorrowMapping(ctx, mappingData)

		assetInTwA, found := k.Market.GetTwa(ctx, assetIn.Id)
		if !found || !assetInTwA.IsPriceActive {
			ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), v.LockedVaultId)
			return nil
		}
		assetInPrice := assetInTwA.Twa

		assetOutTwA, found := k.Market.GetTwa(ctx, assetOut.Id)
		if !found || !assetOutTwA.IsPriceActive {
			ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), v.LockedVaultId)
			return nil
		}
		assetOutPrice := assetOutTwA.Twa
		//assetInPrice is the collateral price
		////Here collateral to be auctioned is received in ucollateral*uusd so inorder to get back amount we divide with uusd of assetIn
		AssetInPrice := sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice))
		if AssetInPrice.Equal(sdk.ZeroDec()) {
			ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), v.LockedVaultId)
			return nil
		}
		AssetOutPrice := sdk.NewDecFromInt(sdk.NewIntFromUint64(assetOutPrice))
		if AssetOutPrice.Equal(sdk.ZeroDec()) {
			ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), v.LockedVaultId)
			return nil
		}

		outflowToken := sdk.NewCoin(assetIn.Denom, v.CollateralToBeAuctioned.Quo(AssetInPrice).TruncateInt())
		inflowToken := sdk.NewCoin(assetOut.Denom, v.CollateralToBeAuctioned.Quo(AssetOutPrice).TruncateInt())
		liquidationPenalty, _ := sdk.NewDecFromStr("0.05")
		fmt.Println("before k.auction.StartLendDutchAuction")
		err := k.Auction.StartLendDutchAuction(ctx, outflowToken, inflowToken, 3, assetOut.Id, assetIn.Id, v.LockedVaultId, v.Owner, liquidationPenalty)
		if err != nil {
			return err
		}
	}
	return nil
}

// FuncMigrateAuctionParams -  Sixth Step: Correct AppID in Auction Params
func FuncMigrateAuctionParams(ctx sdk.Context, k lendkeeper.Keeper) error {
	buffer, _ := sdk.NewDecFromStr("1.2")
	cusp, _ := sdk.NewDecFromStr("0.4")
	auctionParams := types.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 21600,
		Buffer:                 buffer,
		Cusp:                   cusp,
		Step:                   sdk.NewIntFromUint64(360),
		PriceFunctionType:      1,
		DutchId:                3,
		BidDurationSeconds:     10800,
	}
	err := k.AddAuctionParamsData(ctx, auctionParams)
	if err != nil {
		return err
	}
	return nil
}

// func FuncMigrateApp(ctx sdk.Context, k lendkeeper.Keeper) error {
// 	app1 := assettypes.AppData{
// 		Id:               1,
// 		Name:             "CSWAP",
// 		ShortName:        "cswap",
// 		MinGovDeposit:    sdk.ZeroInt(),
// 		GovTimeInSeconds: 300,
// 		GenesisToken:     nil,
// 	}
// 	k.Asset.SetApp(ctx, app1)

// 	genesisToken := assettypes.MintGenesisToken{
// 		AssetId:       9,
// 		GenesisSupply: sdk.NewIntFromUint64(1000000000000000),
// 		IsGovToken:    true,
// 		Recipient:     "comdex1unvvj23q89dlgh82rdtk5su7akdl5932reqarg",
// 	}
// 	var gToken []assettypes.MintGenesisToken
// 	gToken = append(gToken, genesisToken)
// 	app2 := assettypes.AppData{
// 		Id:               2,
// 		Name:             "HARBOR",
// 		ShortName:        "hbr",
// 		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
// 		GovTimeInSeconds: 300,
// 		GenesisToken:     gToken,
// 	}
// 	k.Asset.SetApp(ctx, app2)

// 	app3 := assettypes.AppData{
// 		Id:               3,
// 		Name:             "commodo",
// 		ShortName:        "cmdo",
// 		MinGovDeposit:    sdk.ZeroInt(),
// 		GovTimeInSeconds: 0,
// 		GenesisToken:     nil,
// 	}
// 	k.Asset.SetApp(ctx, app3)
// 	k.Asset.SetAppID(ctx, 3)

// 	return nil
// }
