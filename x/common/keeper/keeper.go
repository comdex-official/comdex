package keeper

import (
	"fmt"

	"encoding/hex"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/common/expected"
	"github.com/comdex-official/comdex/x/common/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		conOps     expected.ContractOpsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	conOps expected.ContractOpsKeeper,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		conOps:     conOps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//nolint:staticcheck
func (k Keeper) SudoContractCall(ctx sdk.Context, contractAddress string, p []byte) error {
	// if err := p.ValidateBasic(); err != nil {
	// 	return err
	// }

	contractAddr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		return sdkerrors.Wrapf(err, "contract")
	}
	data, err := k.conOps.Sudo(ctx, contractAddr, p)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeContractSudoMsg,
		sdk.NewAttribute(types.AttributeKeyResultDataHex, hex.EncodeToString(data)),
	))
	return nil
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}
