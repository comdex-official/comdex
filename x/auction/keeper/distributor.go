package keeper

import (
	esmtypes "github.com/comdex-official/comdex/x/esm/types"

	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) DistributorActivator(ctx sdk.Context, data collectortypes.CollectorAuctionLookupTable,
	inData collectortypes.AssetIdToAuctionLookupTable, killSwitchParams esmtypes.KillSwitchParams, status bool) error {
	if !inData.IsSurplusAuction && !inData.IsAuctionActive && !killSwitchParams.BreakerEnable && !status {
		// to do
		// reduce coin from collector
		// send coin to contract
	}
	return nil
}
