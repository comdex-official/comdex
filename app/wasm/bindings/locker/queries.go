package locker

import (
	lockerKeeper "github.com/comdex-official/comdex/x/locker/keeper"
)

type QueryPlugin struct {
	lockerKeeper *lockerKeeper.Keeper
}

func NewQueryPlugin(
	lockerKeeper *lockerKeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		lockerKeeper: lockerKeeper,
	}
}

/*func (qp QueryPlugin) GetQueryState(addr, denom, blockheight, target string) error {
	poolData, err := qp.lockerKeeper.QueryState(addr, denom, blockheight, target)
	if err != nil {
		return sdkerrors.Wrap(err, "gamm get pool")
	}
	var poolAssets []struct{

	}
	poolAssets.Assets = poolData.GetTotalPoolLiquidity(ctx)
	poolAssets.Shares = sdk.Coin{
		Denom:  gammtypes.GetPoolShareDenom(poolID),
		Amount: poolData.GetTotalShares(),
	}
	return &poolAssets, nil
}*/
