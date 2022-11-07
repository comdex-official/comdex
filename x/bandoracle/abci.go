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
			if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero && ctx.BlockHeight() != block {
				if !k.GetCheckFlag(ctx) {
					msg := k.GetFetchPriceMsg(ctx)
					_, err := k.FetchPrice(ctx, msg)
					if err != nil {
						return err
					}
					k.SetTempFetchPriceID(ctx, 0)
					k.SetCheckFlag(ctx, true)
					k.SetOracleValidationResult(ctx, false)
				} else {
					msg := k.GetFetchPriceMsg(ctx)
					_, err := k.FetchPrice(ctx, msg)
					if err != nil {
						return err
					}
					id := k.GetLastFetchPriceID(ctx)
					req := k.GetTempFetchPriceID(ctx)
					res := k.OraclePriceValidationByRequestID(ctx, req)
					//By default discard height -1 - set while adding band proposal
					//addd new parameter in kvv store to save the accepted discard height
					//one more bool value to save the result of the operation---bydefault false
					if !res && discardheight.IsLT.Zero() {
						discardHeight := ctx.BlockHeight()
					} else if res && discardheight.GT.Zero() {
						if (ctx.BlockHeight() - discardHeight).LT(acceptedDifference) {
							//No issues
							dicardHeight:=-1

						} else if (ctx.BlockHeight() - discardHeight).GTE(acceptedDifference) {
							discardBool:=true
							discardHeight:=-1



						}

					}
					k.SetOracleValidationResult(ctx, res)
					k.SetTempFetchPriceID(ctx, id)
				}
			}
		}
		return nil
	})
}

//if discardBool true------> setCounter to 0 , price to false   ----should be a if condition at the top----> at the end of condition set Discard Bool to false
//if discardBool false-------> nothing to do.