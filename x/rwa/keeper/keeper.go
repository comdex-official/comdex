package keeper

import (
	"fmt"

	"github.com/comdex-official/comdex/x/rwa/expected"
	"github.com/comdex-official/comdex/x/rwa/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	key        storetypes.StoreKey
	paramstore paramtypes.Subspace
	bank       expected.BankKeeper
	account    expected.AccountKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ps paramtypes.Subspace, bank expected.BankKeeper, account expected.AccountKeeper) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		paramstore: ps,
		bank:       bank,
		account:    account,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
