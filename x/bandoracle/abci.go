package bandoracle

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		block := k.GetLastBlockHeight(ctx)
		if block != types.Int64Zero {
			if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero {
				req := k.GetTempFetchPriceID(ctx)
				res := k.OraclePriceValidationByRequestID(ctx, req)
				k.SetOracleValidationResult(ctx, res)
			}
			if ctx.BlockHeight()%types.Int64Twenty-types.Int64One == types.Int64Zero && ctx.BlockHeight() > block+types.Int64Eleven {
				id := k.GetLastFetchPriceID(ctx)
				k.SetTempFetchPriceID(ctx, id)
				msg := k.GetFetchPriceMsg(ctx)
				_, err := k.FetchPrice(ctx, msg)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
