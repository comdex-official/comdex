package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	assetkeeper "github.com/petrichormoney/petri/x/asset/keeper"
	auctionKeeper "github.com/petrichormoney/petri/x/auction/keeper"
	collectorKeeper "github.com/petrichormoney/petri/x/collector/keeper"
	esmKeeper "github.com/petrichormoney/petri/x/esm/keeper"
	lendKeeper "github.com/petrichormoney/petri/x/lend/keeper"
	liquidationKeeper "github.com/petrichormoney/petri/x/liquidation/keeper"
	liquidityKeeper "github.com/petrichormoney/petri/x/liquidity/keeper"
	lockerkeeper "github.com/petrichormoney/petri/x/locker/keeper"
	rewardsKeeper "github.com/petrichormoney/petri/x/rewards/keeper"
	tokenMintkeeper "github.com/petrichormoney/petri/x/tokenmint/keeper"
	vaultKeeper "github.com/petrichormoney/petri/x/vault/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
	asset *assetkeeper.Keeper,
	rewards *rewardsKeeper.Keeper,
	collector *collectorKeeper.Keeper,
	liquidation *liquidationKeeper.Keeper,
	auction *auctionKeeper.Keeper,
	esm *esmKeeper.Keeper,
	vault *vaultKeeper.Keeper,
	lend *lendKeeper.Keeper,
	liquidity *liquidityKeeper.Keeper,
) []wasmkeeper.Option {
	cmd/petriQueryPlugin := NewQueryPlugin(asset, locker, tokenMint, rewards, collector, liquidation, esm, vault, lend, liquidity)

	appDataQueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(cmd/petriQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(*locker, *rewards, *asset, *collector, *liquidation, *auction, *tokenMint, *esm, *vault),
	)

	return []wasm.Option{
		appDataQueryPluginOpt,
		messengerDecoratorOpt,
	}
}
