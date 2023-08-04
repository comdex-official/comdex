package v12

import (
	"fmt"
	auctionkeeperold "github.com/comdex-official/comdex/x/auction/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	auctionkeeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationkeeperold "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v4/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v4/types"
)

// An error occurred during the creation of the CMST/STJUNO pair, as it was mistakenly created in the Harbor app (ID-2) instead of the cSwap app (ID-1).
// As a result, the transaction fee was charged to the creator of the pair, who is entitled to a refund.
// The provided code is designed to initiate the refund process.
// The transaction hash for the pair creation is EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F.
// MintScan - https://www.mintscan.io/comdex/txs/EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F

func CreateUpgradeHandlerV12(
	mm *module.Manager,
	configurator module.Configurator,
	icqkeeper *icqkeeper.Keeper,
	liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	collectorKeeper collectorkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
	auctionKeeperOld auctionkeeperold.Keeper,
	liquidationKeeperOld liquidationkeeperold.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying main net upgrade - v.12.0.0")

		icqparams := icqtypes.DefaultParams()
		icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
		icqkeeper.SetParams(ctx, icqparams)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		InitializeStates(ctx, liquidationKeeper, auctionKeeper)
		Refund(ctx, bankKeeper, collectorKeeper)
		ClearFaultyAuctions(ctx, lendKeeper, auctionKeeperOld, liquidationKeeperOld, bankKeeper)
		return vm, err
	}
}

func InitializeStates(
	ctx sdk.Context,
	liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) {
	dutchAuctionParams := liquidationtypes.DutchAuctionParam{
		Premium:         newDec("1.2"),
		Discount:        newDec("0.7"),
		DecrementFactor: sdk.NewInt(1),
	}
	englishAuctionParams := liquidationtypes.EnglishAuctionParam{DecrementFactor: sdk.NewInt(1)}

	harborParams := liquidationtypes.LiquidationWhiteListing{
		AppId:               2,
		Initiator:           false,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: &englishAuctionParams,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	commodoParams := liquidationtypes.LiquidationWhiteListing{
		AppId:               3,
		Initiator:           false,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: nil,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	liquidationKeeper.SetLiquidationWhiteListing(ctx, harborParams)
	liquidationKeeper.SetLiquidationWhiteListing(ctx, commodoParams)

	appReserveFundsTxDataHbr, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 2)
	if !found {
		appReserveFundsTxDataHbr.AppId = 2
	}
	appReserveFundsTxDataHbr.AssetTxData = append(appReserveFundsTxDataHbr.AssetTxData, liquidationtypes.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataHbr)

	appReserveFundsTxDataCmdo, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 3)
	if !found {
		appReserveFundsTxDataCmdo.AppId = 3
	}
	appReserveFundsTxDataCmdo.AssetTxData = append(appReserveFundsTxDataCmdo.AssetTxData, liquidationtypes.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataCmdo)

	auctionParams := auctionsV2types.AuctionParams{
		AuctionDurationSeconds: 18000,
		Step:                   newDec("0.1"),
		WithdrawalFee:          newDec("0.0"),
		ClosingFee:             newDec("0.0"),
		MinUsdValueLeft:        100000,
		BidFactor:              newDec("0.1"),
		LiquidationPenalty:     newDec("0.1"),
		AuctionBonus:           newDec("0.0"),
	}
	auctionKeeper.SetAuctionParams(ctx, auctionParams)
	auctionKeeper.SetParams(ctx, auctionsV2types.Params{})
	auctionKeeper.SetAuctionID(ctx, 0)
	auctionKeeper.SetUserBidID(ctx, 0)

}

func ClearFaultyAuctions(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
	auctionKeeper auctionkeeperold.Keeper,
	liquidationKeeper liquidationkeeperold.Keeper,
	bankKeeper bankkeeper.Keeper,
) {
	//Send Inflow_token_target_amount to the pool
	//Subtract Inflow_token_target_amount from borrow Position
	//Add the Borrowed amount in poolLBMapping
	//Delete Auction
	//Update BorrowPosition Is liquidated -> false

	// get all the current auctions
	dutchAuctions := auctionKeeper.GetDutchLendAuctions(ctx, 3)
	for _, dutchAuction := range dutchAuctions {
		cPoolModuleName := lendtypes.ModuleAcc1
		reserveModuleName := lendtypes.ModuleName
		//send debt from reserve to the pool
		err := bankKeeper.SendCoinsFromModuleToModule(ctx, reserveModuleName, cPoolModuleName, sdk.NewCoins(dutchAuction.InflowTokenTargetAmount))
		if err != nil {
			fmt.Println(err)
		}
		//send collateral to the reserve from auction module outflow_token_current_amount
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, reserveModuleName, sdk.NewCoins(dutchAuction.OutflowTokenCurrentAmount))
		if err != nil {
			fmt.Println(err)
		}

		borrowPos := lendKeeper.GetBorrowByUserAndAssetID(ctx, dutchAuction.VaultOwner.String(), dutchAuction.InflowTokenTargetAmount.Denom, dutchAuction.AssetOutId)
		borrowPos.AmountOut.Amount = borrowPos.AmountOut.Amount.Sub(dutchAuction.InflowTokenTargetAmount.Amount)
		borrowPos.IsLiquidated = false
		lendKeeper.SetBorrow(ctx, borrowPos)

		poolAssetLBMappingData, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 1, dutchAuction.AssetInId)

		poolAssetLBMappingData.TotalBorrowed = poolAssetLBMappingData.TotalBorrowed.Add(borrowPos.AmountOut.Amount)
		lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
		lockedVault, found := liquidationKeeper.GetLockedVault(ctx, 3, dutchAuction.LockedVaultId)
		if found {
			liquidationKeeper.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
		}
		err = auctionKeeper.SetHistoryDutchLendAuction(ctx, dutchAuction)
		if err != nil {
			fmt.Println(err)
		}
		err = auctionKeeper.DeleteDutchLendAuction(ctx, dutchAuction)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func newDec(i string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(i)
	return dec
}
