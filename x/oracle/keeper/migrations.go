package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
}

// NewMigrator creates a Migrator.
func NewMigrator(keeper *Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// MigrateValidatorSet fixes the validator set being stored as map
// causing non determinism by storing it as a list.
func (m Migrator) MigrateValidatorSet(ctx sdk.Context) {
	m.keeper.SetValidatorRewardSet(ctx)
}
