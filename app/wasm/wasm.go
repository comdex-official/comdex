package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	tokenMintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
	asset *assetkeeper.Keeper,
) []wasmkeeper.Option {

	comdexQueryPlugin := NewQueryPlugin(asset, locker, tokenMint)

	appDataqueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(comdexQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(*locker),
	)

	return []wasm.Option{
		appDataqueryPluginOpt,
		messengerDecoratorOpt,
	}
}
