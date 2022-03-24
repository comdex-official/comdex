package rewards

import (
	"github.com/comdex-official/comdex/x/rewards/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	k.EnableMintingRewards(ctx)
	k.DisableMintingRewards(ctx)
	k.TriggerRewards(ctx)
}
