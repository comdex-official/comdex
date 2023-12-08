package bandoracle

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	block := k.GetLastBlockHeight(ctx)
	if block != types.Int64Zero {
		// if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero && ctx.BlockHeight() != block {
		if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero {
			if !k.GetCheckFlag(ctx) {
				msg := k.GetFetchPriceMsg(ctx)
				_, err := k.FetchPrice(ctx, msg)
				if err != nil {
					ctx.Logger().Error("Error in Fetch Price in 1st condition")
				}
				k.SetTempFetchPriceID(ctx, 0)
				k.SetCheckFlag(ctx, true)
				k.SetOracleValidationResult(ctx, false)
			} else {
				msg := k.GetFetchPriceMsg(ctx)
				_, err := k.FetchPrice(ctx, msg)
				if err != nil {
					ctx.Logger().Error("Error in Fetch Price in 2nd condition")
				}
				id := k.GetLastFetchPriceID(ctx)
				req := k.GetTempFetchPriceID(ctx)
				res := k.OraclePriceValidationByRequestID(ctx, req)
				discardData := k.GetDiscardData(ctx)
				// By default discard height -1 - set while adding band proposal
				// addd new parameter in kvv store to save the accepted discard height
				// one more bool value to save the result of the operation---bydefault false
				if !res && discardData.BlockHeight < 0 {
					discardData.BlockHeight = ctx.BlockHeight()
				} else if res && discardData.BlockHeight > 0 {
					if (ctx.BlockHeight() - discardData.BlockHeight) < msg.AcceptedHeightDiff {
						// No issues
						discardData.BlockHeight = -1
					} else if (ctx.BlockHeight() - discardData.BlockHeight) >= msg.AcceptedHeightDiff {
						discardData.DiscardBool = true
						discardData.BlockHeight = -1
					}
				}
				k.SetDiscardData(ctx, discardData)
				k.SetOracleValidationResult(ctx, res)
				k.SetTempFetchPriceID(ctx, id)
			}
		}
	}
}

// if discardBool true------> setCounter to 0 , price to false   ----should be a if condition at the top----> at the end of condition set Discard Bool to false
// if discardBool false-------> nothing to do.
