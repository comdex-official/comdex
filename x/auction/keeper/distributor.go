package keeper

import (
	esmtypes "github.com/comdex-official/comdex/x/esm/types"

	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) DistributorActivator(ctx sdk.Context, data collectortypes.AppAssetIdToAuctionLookupTable, killSwitchParams esmtypes.KillSwitchParams, status bool) error {
	if !data.IsSurplusAuction && !data.IsAuctionActive && !killSwitchParams.BreakerEnable && !status {
		// to do
		// reduce coin from collector
		// send coin to contract
	}
	return nil
}
