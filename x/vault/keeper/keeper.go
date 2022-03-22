package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/vault/expected"
	"github.com/comdex-official/comdex/x/vault/types"
)

type Keeper struct {
	cdc     codec.BinaryCodec
	key     sdk.StoreKey
	bank    expected.BankKeeper
	account expected.AccountKeeper
	asset   expected.AssetKeeper
	oracle  expected.OracleKeeper
	bandoracle expected.BandoracleKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	bank expected.BankKeeper,
	account expected.AccountKeeper,
	asset expected.AssetKeeper,
	oracle expected.OracleKeeper,
	bandoracle expected.BandoracleKeeper,
) Keeper {
	return Keeper{
		cdc:     cdc,
		key:     key,
		bank:    bank,
		account: account,
		asset:   asset,
		oracle:  oracle,
		bandoracle: bandoracle,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
