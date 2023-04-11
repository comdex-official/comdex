package keeper

import (

	// "github.com/comdex-official/comdex/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	auctionType, _ := k.liquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)

	//Dutch Auction Model Followed for auction type 0
	if auctionType.AuctionType== 0{
		//Trigger Dutch Auction

	}

	return nil
}
