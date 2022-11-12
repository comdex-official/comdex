package keeper

import (
	"fmt"

	v5types "github.com/comdex-official/comdex/x/collector/migrations/v5"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate3to4 migrates from version 1 to 2.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	fmt.Println("Migrate3to4Collector")
	return v5types.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
