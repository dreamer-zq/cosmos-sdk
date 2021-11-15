package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v016 "github.com/cosmos/cosmos-sdk/x/nft/example/nft/migrations/v016"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v016.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
