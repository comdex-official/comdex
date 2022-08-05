package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
)

func (k Keeper) HasAssetForDenom(ctx sdk.Context, id string) bool {
	return k.asset.HasAssetForDenom(ctx, id)
}

func (k Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	return k.asset.HasAsset(ctx, id)
}
func (k Keeper) GetAssetForDenom(ctx sdk.Context, id string) (types.Asset, bool) {
	return k.asset.GetAssetForDenom(ctx, id)
}

func (k Keeper) GetApp(ctx sdk.Context, id uint64) (types.AppData, bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (types.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return collectortypes.SendCoinFromModuleToModuleIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (mintData types.MintGenesisToken, found bool) {
	return k.asset.GetMintGenesisTokenData(ctx, appID, assetID)
}

func (k Keeper) GetAuctionParams(ctx sdk.Context, AppID uint64) (asset auctiontypes.AuctionParams, found bool) {
	return k.auction.GetAuctionParams(ctx, AppID)
}
