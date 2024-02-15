package keeper

import (
	"fmt"

	"encoding/hex"
	"github.com/cometbft/cometbft/libs/log"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/common/expected"
	"github.com/comdex-official/comdex/x/common/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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

func (k Keeper) CheckSecurityAddress(ctx sdk.Context, from string) bool {
	params := k.GetParams(ctx)
	for _, addr := range params.SecurityAddress {
		if addr == from {
			return true
		}
	}
	return false
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) SinglePlayer(ctx sdk.Context, contractAddress string, ResolveSinglePlayer []byte, gameName string) {
	logger := k.Logger(ctx)
	err := k.SudoContractCall(ctx, contractAddress, ResolveSinglePlayer)
	if err != nil {
		logger.Error(fmt.Sprintf("Game %s contract call error for single-player", gameName))
	} else {
		logger.Info(fmt.Sprintf("Game %s contract call for single-player success", gameName))
	}
}

func (k Keeper) MultiPlayer(ctx sdk.Context, contractAddress string, SetupMultiPlayer []byte, ResolveMultiPlayer []byte, gameName string) {
	logger := k.Logger(ctx)
	err := k.SudoContractCall(ctx, contractAddress, SetupMultiPlayer)
	if err != nil {
		logger.Error(fmt.Sprintf("Game %s contract call error for setup multi-player", gameName))
	} else {
		logger.Info(fmt.Sprintf("Game %s contract call for setup multi-player success", gameName))
	}

	err = k.SudoContractCall(ctx, contractAddress, ResolveMultiPlayer)
	if err != nil {
		logger.Error(fmt.Sprintf("Game %s contract call error for resolve multi-player", gameName))
	} else {
		logger.Info(fmt.Sprintf("Game %s contract call for resolve multi-player success", gameName))
	}
}
