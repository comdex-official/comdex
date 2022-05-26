package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
)

func RegisterCustomPlugins(
	locker lockerkeeper.Keeper,
) []wasmkeeper.Option {

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(locker),
	)

	return []wasm.Option{
		messengerDecoratorOpt,
	}
}
