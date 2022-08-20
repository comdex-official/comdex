package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
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

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData lockertypes.LockerLookupTableData, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appID, assetID)
}

func (k Keeper) GetReward(ctx sdk.Context, appId, assetID uint64) (rewards rewardstypes.InternalRewards, found bool) {
	return k.GetReward(ctx, appId, assetID)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerID uint64) (locker lockertypes.Locker, found bool) {
	return k.GetLocker(ctx, lockerID)
}