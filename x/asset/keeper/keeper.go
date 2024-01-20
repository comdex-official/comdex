package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"

	storetypes "cosmossdk.io/store/types"

	"github.com/comdex-official/comdex/x/asset/expected"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	key        storetypes.StoreKey
	params     paramstypes.Subspace
	account    expected.AccountKeeper
	bank       expected.BankKeeper
	rewards    expected.RewardsKeeper
	vault      expected.VaultKeeper
	bandoracle expected.Bandoraclekeeper
	authority  string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	params paramstypes.Subspace,
	account expected.AccountKeeper,
	bank expected.BankKeeper,
	rewards expected.RewardsKeeper,
	vault expected.VaultKeeper,
	bandoracle expected.Bandoraclekeeper,
	authority string,
) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(assettypes.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		key:        key,
		params:     params,
		account:    account,
		bank:       bank,
		rewards:    rewards,
		vault:      vault,
		bandoracle: bandoracle,
		authority:  authority,
	}
}

func (k Keeper) Store(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.key)
}
