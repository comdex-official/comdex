package common

import (
	"github.com/comdex-official/comdex/x/common/keeper"
	commonTypes "github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	allContracts := k.GetAllContract(ctx)

	for _, data := range allContracts {
		if data.GameType == 1 {
			k.SinglePlayer(ctx, data.ContractAddress, commonTypes.ResolveSinglePlayer, data.GameName)
		} else if data.GameType == 2 {
			k.MultiPlayer(ctx, data.ContractAddress, commonTypes.SetupMultiPlayer, commonTypes.ResolveMultiPlayer, data.GameName)
		} else {
			k.SinglePlayer(ctx, data.ContractAddress, commonTypes.ResolveSinglePlayer, data.GameName)
			k.MultiPlayer(ctx, data.ContractAddress, commonTypes.SetupMultiPlayer, commonTypes.ResolveMultiPlayer, data.GameName)
		}
	}
}
