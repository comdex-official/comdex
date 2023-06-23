package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	//  Migrate these 3 for their store: (export fix)
	// 		FundModBal

	store := ctx.KVStore(storeKey)
	err := MigrateExternalRewardLends(store, cdc)
	if err != nil {
		return err
	}
	return err
}

func MigrateExternalRewardLends(store sdk.KVStore, cdc codec.BinaryCodec) error {
	key := types.ExternalRewardsLendMappingKey(3)
	value := store.Get(key)
	var extRew types.LendExternalRewards
	cdc.MustUnmarshal(value, &extRew)

	store.Delete(key)
	SetExternalRewardLends(store, cdc, extRew)

	return nil
}

func SetExternalRewardLends(store sdk.KVStore, cdc codec.BinaryCodec, extRew types.LendExternalRewards) {
	var (
		key   = types.ExternalRewardsLendMappingKey(extRew.AppMappingId)
		value = cdc.MustMarshal(&extRew)
	)

	store.Set(key, value)
}
