package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) WhitelistAsset(ctx sdk.Context, appMappingId uint64, assetId []uint64) error {

	internalRewards := types.InternalRewards{
		App_mapping_ID: appMappingId,
		Asset_ID:       assetId,
	}

	k.SetReward(ctx, internalRewards)
	return nil
}

func (k Keeper) RemoveWhitelistAsset(ctx sdk.Context, appMappingId uint64, assetId uint64) error {

	rewards, found := k.GetReward(ctx, appMappingId)
	if found != true {
		return nil
	}
	var newAssetIds []uint64
	fmt.Println(rewards.Asset_ID)
	for i := range rewards.Asset_ID {
		if assetId != rewards.Asset_ID[i] {
			newAssetId := rewards.Asset_ID[i]
			newAssetIds = append(newAssetIds, newAssetId)
		}

	}
	newRewards := types.InternalRewards{
		App_mapping_ID: appMappingId,
		Asset_ID:       newAssetIds,
	}
	k.SetReward(ctx, newRewards)
	return nil
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}
