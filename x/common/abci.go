package common

import (
	"fmt"
	"github.com/comdex-official/comdex/x/common/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

	Msg := []byte(`{"resolve_bet":{}}`)

	allContracts := k.GetAllContract(ctx)
	logger := k.Logger(ctx)

	for _, data := range allContracts {
		logger.Info(fmt.Sprintf("Game Id %d contract call", data.GameId))
		_ = k.SudoContractCall(ctx, data.ContractAddr, Msg)
	}

}
