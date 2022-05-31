package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	tokenMintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
	asset *assetkeeper.Keeper,
	rewards *rewardsKeeper.Keeper,
) []wasmkeeper.Option {

	comdexQueryPlugin := NewQueryPlugin(asset, locker, tokenMint, rewards)

	appDataqueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(comdexQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(*locker, *rewards),
	)

	return []wasm.Option{
		appDataqueryPluginOpt,
		messengerDecoratorOpt,
	}
}
