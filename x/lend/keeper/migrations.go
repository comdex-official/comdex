package keeper

import (
	v5types "github.com/petrichormoney/petri/x/lend/migrations/v5"

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
	return v5types.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
