package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k *msgServer) TriggerAuction(ctx sdk.Context, appid uint64) {

	_, _ = k.GetAppidToAssetCollectorMapping(ctx, appid)
	_, _ = k.GetCollectorLookupTable(ctx, appid)

	// check for app_id in both get calls
	// match asset id's in both
	// check if net > surplus threshold + lot_size -> surplus auction
	// check if net < debt threshold - lot_size -> debt auction
	// update historical data

}
