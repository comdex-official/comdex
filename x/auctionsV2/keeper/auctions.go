package keeper

import(
	
	// "github.com/comdex-official/comdex/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

)



func (k Keeper) AuctionActivator(ctx sdk.Context,lockedVault liquidationtypes.LockedVault) error {

	//Using app id provided in the lockedVault Struct
	// Using the app id to fetch the app data whitelisted in the liquidation module to find the auction type selected by the app.
	appWhitelistedData:=k.liquidation.GetLiquidationWhiteListing(ctx,lockedVault.AppId)

	



	
	return nil
}