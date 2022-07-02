package keeper

import (
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleAddESMTriggerParamsRecords(ctx sdk.Context, p *types.ESMTriggerParamsProposal) error {
	return k.AddESMTriggerParamsRecords(ctx, p.EsmTriggerParams)
}
