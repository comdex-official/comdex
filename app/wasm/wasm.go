package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	esmKeeper "github.com/comdex-official/comdex/x/esm/keeper"
	gaslessKeeper "github.com/comdex-official/comdex/x/gasless/keeper"
	lendKeeper "github.com/comdex-official/comdex/x/lend/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidityKeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	marketKeeper "github.com/comdex-official/comdex/x/market/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	tokenfactorykeeper "github.com/comdex-official/comdex/x/tokenfactory/keeper"
	tokenMintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
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
	market *marketKeeper.Keeper,
	bank bankkeeper.Keeper,
	tokenfactory *tokenfactorykeeper.Keeper,
	gasless *gaslessKeeper.Keeper,
) []wasmkeeper.Option {
	comdexQueryPlugin := NewQueryPlugin(asset, locker, tokenMint, rewards, collector, liquidation, esm, vault, lend, liquidity, market, bank, tokenfactory, gasless)

	appDataQueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(comdexQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(*locker, *rewards, *asset, *collector, *liquidation, *auction, *tokenMint, *esm, *vault, *liquidity, bank, *tokenfactory, *gasless),
	)

	return []wasm.Option{
		appDataQueryPluginOpt,
		messengerDecoratorOpt,
	}
}
