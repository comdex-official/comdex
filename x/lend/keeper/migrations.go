package keeper

import (
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	v4_4_0_beta "github.com/comdex-official/comdex/x/lend/migrations/v4.4.0.beta"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator codes goes here

type Migrator struct {
	keeper            Keeper
	liquidationkeeper liquidationkeeper.Keeper
	auctionkeeper     auctionkeeper.Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateTo4_4_0beta(ctx sdk.Context) error {

	return v4_4_0_beta.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
