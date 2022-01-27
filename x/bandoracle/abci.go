package bandoracle

import (
	"fmt"
	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	fmt.Println("bandoracle executed")
	if ctx.BlockHeight()>=39 && ctx.BlockHeight()%20==0{
		msg := types.NewMsgFetchPriceData(
			types.ModuleName,
			37,
			"channel-0",
			&types.FetchPriceCallData{[]string{"BTC","ATOM"}, 1000000} ,
			6,
			3,
			sdk.Coins{sdk.NewCoin("uband", sdk.NewInt(30000))},
			600000,
			600000,
		)
		fmt.Println("beginblocker")
		fmt.Println(msg)
		fmt.Println("{{{{{{{{{{{{{{{{{{{{{{{{{{{{[")
		_, err := k.FetchPrice(ctx, *msg)
		if err != nil {
			return
		}
	}
}
