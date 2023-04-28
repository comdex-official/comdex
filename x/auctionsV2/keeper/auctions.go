package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (k Keeper) AuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	auctionType, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)

	//Dutch Auction Model Followed for auction type 0
	if auctionType.AuctionType {

		//Trigger Dutch Auction
		err := k.DutchAuctionActivator(ctx, liquidationData)
		if err != nil {

		}

	} else {
		//Trigger English Auction
	}

	return nil
}

func (k Keeper) DutchAuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Getting previous auction ID
	auctionID := k.GetAuctionID(ctx)

	//Price Calculation Function to determine auction different stage price
	dutchAuctionParams, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)
	dap := dutchAuctionParams.DutchAuctionParam
	extPair, _ := k.asset.GetPairsVault(ctx, liquidationData.ExtendedPairId)
	pair, _ := k.asset.GetPair(ctx, extPair.PairId)
	twaDataCollateral, found := k.market.GetTwa(ctx, pair.AssetIn)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	twaDataDebt, found := k.market.GetTwa(ctx, pair.AssetOut)
	if !found || !twaDataDebt.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	CollateralTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), dap.Buffer)
	CollateralTokenEndPrice := k.getOutflowTokenEndPrice(CollateralTokenInitialPrice, dap.Cusp)
	auctionParams, _ := k.GetAuctionParams(ctx)

	//Saving liquidation data to the auction struct
	auctionData := types.Auctions{
		AuctionId:                   auctionID + 1,
		CollateralToken:             liquidationData.CollateralToken,
		DebtToken:                   liquidationData.DebtToken,
		CollateralTokenInitialPrice: CollateralTokenInitialPrice,
		CollateralTokenCurrentPrice: sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa))),
		CollateralTokenEndPrice:     CollateralTokenEndPrice,
		DebtTokenCurrentPrice:       sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa))),
		LockedVaultId:               liquidationData.LockedVaultId,
		StartTime:                   ctx.BlockTime(),
		EndTime:                     ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AppId:                       liquidationData.AppId,
		AuctionType:                 liquidationData.AuctionType,
	}

	return nil
}

//AUCTIONITERATOR
// -> DUCTHAUCTIONITERATOR
// -> ENGLISHAUCTIONITERATOR
