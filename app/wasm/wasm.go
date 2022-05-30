package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	locker2 "github.com/comdex-official/comdex/app/wasm/bindings/locker"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
)

func RegisterCustomPlugins(
	locker lockerkeeper.Keeper,
) []wasmkeeper.Option {

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		locker2.CustomMessageDecorator(locker),
	)

	return []wasm.Option{
		messengerDecoratorOpt,
	}
}
