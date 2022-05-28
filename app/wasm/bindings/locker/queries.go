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
