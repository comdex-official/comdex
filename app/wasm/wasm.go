package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	locker2 "github.com/comdex-official/comdex/app/wasm/bindings/locker"
	tokenMint2 "github.com/comdex-official/comdex/app/wasm/bindings/tokenmint"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	tokenMintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := locker2.NewQueryPlugin(locker)
	tokenMintwasmQueryPlugin := tokenMint2.NewQueryPlugin(tokenMint)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: locker2.CustomQuerier(wasmQueryPlugin),
	})
	tokenMintqueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: tokenMint2.CustomQuerier(tokenMintwasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		locker2.CustomMessageDecorator(*locker),
	)

	return []wasm.Option{
		queryPluginOpt,
		tokenMintqueryPluginOpt,
		messengerDecoratorOpt,
	}
}
