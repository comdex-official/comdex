package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/comdex-official/comdex/x/collector/expected"
	"github.com/comdex-official/comdex/x/collector/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		asset      expected.AssetKeeper
		auction    expected.AuctionKeeper
		locker     expected.LockerKeeper
		rewards    expected.RewardsKeeper
		paramStore paramtypes.Subspace
		bank       expected.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	asset expected.AssetKeeper,
	auction expected.AuctionKeeper,
	locker expected.LockerKeeper,
	rewards expected.RewardsKeeper,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		asset:      asset,
		auction:    auction,
		locker:     locker,
		rewards:    rewards,
		paramStore: ps,
		bank:       bank,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}
