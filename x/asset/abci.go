package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"

	mv510 "github.com/comdex-official/comdex/app/upgrades/mainnet/v510"
	"github.com/comdex-official/comdex/x/asset/keeper"
)

func BeginBlocker(ctx sdk.Context,
	_ abci.RequestBeginBlock,
	k keeper.Keeper,
	bank bankkeeper.Keeper,
	staking stakingkeeper.Keeper,
	mint mintkeeper.Keeper,
) {
	if ctx.BlockHeight() == 55 {
		mv510.MintLostTokens(ctx, bank, staking, mint)
	}
}
