package keeper

import (
	"fmt"

	v5types "github.com/petrichormoney/petri/x/asset/migrations/v5"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	fmt.Println("Migrate1to2")
	return v5types.MigrateStore(ctx, m.keeper.key, m.keeper.cdc)
}
