package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"encoding/hex"

	"cosmossdk.io/log"

	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/common/expected"
	"github.com/comdex-official/comdex/x/common/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		conOps     expected.ContractOpsKeeper
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	conOps expected.ContractOpsKeeper,
	authority string,

) Keeper {

	return Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		conOps:     conOps,
		authority:  authority,
	}
}

// GetAuthority returns the x/common module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//nolint:staticcheck
func (k Keeper) SudoContractCall(ctx sdk.Context, contractAddress string, p []byte) error {

	contractAddr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		return errorsmod.Wrapf(err, "contract")
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

func (k Keeper) CheckSecurityAddress(ctx sdk.Context, from string) bool {
	params := k.GetParams(ctx)
	for _, addr := range params.SecurityAddress {
		if addr == from {
			return true
		}
	}
	return false
}

func (k Keeper) Store(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.storeKey)
}
