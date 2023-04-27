package keeper

import (
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

 //Saving liquidation data to the auction struct
	auctionData := types.Auctions{
		AuctionId: auctionID+1,
		CollateralToken: liquidationData.CollateralToken,
		DebtToken: liquidationData.DebtToken,
		CollateralTokenInitialPrice: 0,
		CollateralTokenCurrentPrice: 0,
		CollateralTokenEndPrice: 0,
		DebtTokenCurrentPrice: 0,
		LockedVaultId: liquidationData.LockedVaultId,
		StartTime: 0,
		EndTime: 0,
		AppId: liquidationData.AppId,
		AuctionType: liquidationData.AuctionType,
		




	}

	return nil
}

//AUCTIONITERATOR
// -> DUCTHAUCTIONITERATOR
// -> ENGLISHAUCTIONITERATOR
