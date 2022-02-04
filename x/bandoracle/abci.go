package bandoracle

import (
	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	if ctx.BlockHeight() >= 99 && ctx.BlockHeight()%20 == 0 {
		msg := types.NewMsgFetchPriceData(
			types.ModuleName,
			112,
			"channel-0",
			&types.FetchPriceCallData{[]string{"ATOM"}, 1000000},
			6,
			3,
			sdk.Coins{sdk.NewCoin("uband", sdk.NewInt(250000))},
			600000,
			600000,
		)
		_, err := k.FetchPrice(ctx, *msg)
		if err != nil {
			return
		}
	}
}
