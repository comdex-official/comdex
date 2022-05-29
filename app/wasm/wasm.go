package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	asset2 "github.com/comdex-official/comdex/app/wasm/bindings/asset"
	locker2 "github.com/comdex-official/comdex/app/wasm/bindings/locker"
	tokenMint2 "github.com/comdex-official/comdex/app/wasm/bindings/tokenmint"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	tokenMintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
	asset *assetkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := locker2.NewQueryPlugin(locker)
	tokenMintwasmQueryPlugin := tokenMint2.NewQueryPlugin(tokenMint)
	appDatawasmQueryPlugin := asset2.NewQueryPlugin(asset)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: locker2.CustomQuerier(wasmQueryPlugin),
	})
	tokenMintqueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: tokenMint2.CustomQuerier(tokenMintwasmQueryPlugin),
	})
	appDataqueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: asset2.CustomQuerier(appDatawasmQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		locker2.CustomMessageDecorator(*locker),
	)

	return []wasm.Option{
		queryPluginOpt,
		tokenMintqueryPluginOpt,
		appDataqueryPluginOpt,
		messengerDecoratorOpt,
	}
}
