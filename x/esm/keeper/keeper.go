package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/esm/expected"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/esm/types"
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
		asset      expected.AssetKeeper
		bank       expected.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	asset expected.AssetKeeper,
	bank expected.BankKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		asset:      asset,
		bank:       bank,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) DepositESM(ctx sdk.Context, lenderAddr string, AppID uint64, Amount sdk.Coin) error {
	//TODO:
	// burn from token_mint module

	appData, found := k.GetApp(ctx, AppID)
	if !found {
		return types.ErrAppDataNotFound
	}
	var govTokenID uint64
	for _, v := range appData.GenesisToken {
		if v.IsGovToken {
			govTokenID = v.AssetId
		}
	}
	govAsset, _ := k.GetAsset(ctx, govTokenID)
	if Amount.Denom != govAsset.Denom {
		return types.ErrBadOfferCoinType
	}
	esmTriggerParams, found := k.GetESMTriggerParams(ctx, AppID)
	if !found {
		return types.ErrESMTriggerParamsNotFound
	}
	currentDeposit, found := k.GetCurrentDepositStats(ctx, AppID)

	if currentDeposit.Balance.Amount.Equal(esmTriggerParams.TargetValue.Amount) {
		return types.ErrDepositForAppReached
	}
	if Amount.Amount.GT(esmTriggerParams.TargetValue.Amount) {
		return types.ErrAmtExceedsTargetValue
	}
	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(Amount)); err != nil {
		return err
	}

	if !found {
		newCurrentDeposit := types.CurrentDepositStats{
			AppId:   AppID,
			Balance: Amount,
		}
		k.SetCurrentDepositStats(ctx, newCurrentDeposit)
	} else {
		newCurrentDeposit := types.CurrentDepositStats{
			AppId:   AppID,
			Balance: currentDeposit.Balance.Add(Amount),
		}
		k.SetCurrentDepositStats(ctx, newCurrentDeposit)
	}

	return nil
}

func (k Keeper) ExecuteESM(ctx sdk.Context, lenderAddr string, AppID uint64) error {
	return nil
}
