package keeper

import (
	"time"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Dutch Auction Model Followed for auction type true
	if liquidationData.AuctionType {

		//Trigger Dutch Auction
		err := k.DutchAuctionActivator(ctx, liquidationData)
		if err != nil {

		}

	} else {
		//Trigger English Auction
		err := k.EnglishAuctionActivator(ctx, liquidationData)
		if err != nil {

		}
	}

	return nil
}

func (k Keeper) DutchAuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Getting previous auction ID
	auctionID := k.GetAuctionID(ctx)

	//Price Calculation Function to determine auction different stage price
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)

	if !liquidationWhitelistingAppData.IsDutchActivated {
		return types.ErrDutchAuctionDisabled
	}
	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	pair, _ := k.asset.GetPair(ctx, liquidationData.PairId)

	twaDataCollateral, found := k.market.GetTwa(ctx, pair.AssetIn)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	twaDataDebt, found := k.market.GetTwa(ctx, pair.AssetOut)
	if !found || !twaDataDebt.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	//Checking if DEBT  token is CMST  then setting its price to $1 , else all tokens price will come from oracle.
	if liquidationData.IsDebtCmst {
		twaDataDebt.Twa = 1000000
	}

	//Some params will come from the specific app and they could be configured by them ,
	//rest of the params like auction duration and fees and other params will be based on comdex to edit based on governance
	//Understanding different Params:
	//Premium : Initial Price i.e price of the collateral at which the auction will start
	//Discount: Final Price , i.e less than the oracle price of the collateral asset and at this , auction would end
	//Decrement Factor:     Linear decrease in the price of the collateral every block is governed by this.
	CollateralTokenInitialPrice := k.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), dutchAuctionParams.Premium)
	// CollateralTokenEndPrice := k.getOutflowTokenEndPrice(CollateralTokenInitialPrice, dutchAuctionParams.Cusp)
	auctionParams, _ := k.GetAuctionParams(ctx)

	//Saving liquidation data to the auction struct
	auctionData := types.Auctions{
		AuctionId:                   auctionID + 1,
		CollateralToken:             liquidationData.CollateralToken,
		DebtToken:                   liquidationData.DebtToken,
		CollateralTokenAuctionPrice: CollateralTokenInitialPrice,
		CollateralTokenOraclePrice:  sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa))),
		DebtTokenOraclePrice:        sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa))),
		LockedVaultId:               liquidationData.LockedVaultId,
		StartTime:                   ctx.BlockTime(),
		EndTime:                     ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AppId:                       liquidationData.AppId,
		AuctionType:                 liquidationData.AuctionType,
	}

	k.SetAuctionID(ctx, auctionData.AuctionId)
	k.SetAuction(ctx, auctionData)

	return nil
}

// AUCTIONITERATOR
// -> DUCTHAUCTIONITERATOR
// -> ENGLISHAUCTIONITERATOR
func (k Keeper) EnglishAuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	return nil

}
