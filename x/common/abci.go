package common

import (
	// "encoding/hex"
	"fmt"
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/common/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Msg := []byte(`{"resolve_bet":{"app_hash":"` + hex.EncodeToString(ctx.BlockHeader().AppHash) + `", "block_time":"` + ctx.BlockTime().String() + `"}}`)

	Msg := []byte(`{"resolve_bet":{}}`)

	allContracts := k.GetAllContract(ctx)
	logger := k.Logger(ctx)

	for _, data := range allContracts {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			err := k.SudoContractCall(ctx, data.ContractAddr, Msg)
			if err != nil {
				logger.Error(fmt.Sprintf("Game Id %d contract call error", data.GameId))
				return err
			} 
			logger.Info(fmt.Sprintf("Game Id %d contract call", data.GameId))
			return nil
		})
	}

}
