package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/esm/expected"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		asset      expected.AssetKeeper
		vault      expected.VaultKeeper
		bank       expected.BankKeeper
		market     expected.MarketKeeper
		tokenmint  expected.Tokenmint
		collector  expected.Collector
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	bank expected.BankKeeper,
	market expected.MarketKeeper,
	tokenmint expected.Tokenmint,
	collector expected.Collector,
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
		asset:      asset,
		vault:      vault,
		bank:       bank,
		market:     market,
		tokenmint:  tokenmint,
		collector:  collector,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) DepositESM(ctx sdk.Context, depositorAddr string, AppID uint64, Amount sdk.Coin) error {
	// take deposits from the user
	// send tokens to tokenMint module and burn them
	// update global deposit stats and user deposit stats checking if trigger params reached

	appData, found := k.asset.GetApp(ctx, AppID)
	if !found {
		return types.ErrAppDataNotFound
	}
	var govTokenID uint64
	for _, v := range appData.GenesisToken {
		if v.IsGovToken {
			govTokenID = v.AssetId
		}
	}
	govAsset, found := k.asset.GetAsset(ctx, govTokenID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if Amount.Denom != govAsset.Denom {
		return types.ErrBadOfferCoinType
	}
	esmTriggerParams, found := k.GetESMTriggerParams(ctx, AppID)
	if !found {
		return types.ErrESMTriggerParamsNotFound
	}
	existingDeposit, found := k.GetCurrentDepositStats(ctx, AppID)
	if !found {
		existingDeposit = types.CurrentDepositStats{
			AppId:   AppID,
			Balance: Amount,
		}
	} else {
		if existingDeposit.Balance.Amount.GT(esmTriggerParams.TargetValue.Amount) {
			return types.ErrDepositForAppReached
		}
		existingDeposit.Balance = existingDeposit.Balance.Add(Amount)
	}
	addr, _ := sdk.AccAddressFromBech32(depositorAddr)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, tokenminttypes.ModuleName, sdk.NewCoins(Amount)); err != nil {
		return err
	}
	if err1 := k.tokenmint.BurnTokensForApp(ctx, AppID, govAsset.Id, Amount.Amount); err1 != nil {
		return err1
	}
	k.SetCurrentDepositStats(ctx, existingDeposit)
	// currentDeposit, found := k.GetCurrentDepositStats(ctx, AppID)

	userDeposits, found := k.GetUserDepositByApp(ctx, depositorAddr, AppID)
	if !found {
		userDeposits = types.UsersDepositMapping{
			AppId:     AppID,
			Depositor: depositorAddr,
			Deposits:  Amount,
		}
		k.SetUserDepositByApp(ctx, userDeposits)
	} else {
		userDeposits.Deposits = userDeposits.Deposits.Add(Amount)
		k.SetUserDepositByApp(ctx, userDeposits)
	}

	return nil
}

func (k Keeper) ExecuteESM(ctx sdk.Context, executor string, AppID uint64) error {
	// checking if target deposit amount is reached
	// setting end time by adding coolOff period
	_, found := k.asset.GetApp(ctx, AppID)
	if !found {
		return types.ErrAppDataNotFound
	}
	_, found = k.GetESMStatus(ctx, AppID)
	if found {
		return types.ErrESMAlreadyExecuted
	}

	esmTriggerParams, found := k.GetESMTriggerParams(ctx, AppID)
	if !found {
		return types.ErrESMTriggerParamsNotFound
	}

	currentDeposit, found := k.GetCurrentDepositStats(ctx, AppID)
	if !found {
		return types.ErrDepositForAppNotFound
	}

	if currentDeposit.Balance.Amount.GTE(esmTriggerParams.TargetValue.Amount) {
		ESMStatus := types.ESMStatus{
			AppId:     AppID,
			Executor:  executor,
			Status:    true,
			StartTime: ctx.BlockTime(),
			EndTime:   ctx.BlockTime().Add(time.Duration(esmTriggerParams.CoolOffPeriod) * time.Second),
		}
		k.SetESMStatus(ctx, ESMStatus)
	} else {
		return types.ErrCurrentDepositNotReached
	}
	return nil
}

// func (k Keeper) CalculateCollateral(ctx sdk.Context, appID uint64, amount sdk.Coin, esmDataAfterCoolOff types.DataAfterCoolOff, from string) error {
// 	userAddress, err := sdk.AccAddressFromBech32(from)
// 	if err != nil {
// 		return err
// 	}
// 	assetInID, _ := k.asset.GetAssetForDenom(ctx, amount.Denom)

// 	// initializing userWorth
// 	userWorth := sdk.ZeroDec()
// 	for _, v := range esmDataAfterCoolOff.DebtAsset {
// 		if v.AssetID == assetInID.Id {
// 			userWorth = v.DebtTokenWorth.Mul(sdk.NewDecFromInt(amount.Amount))
// 			break
// 		}
// 	}

// 	for i, data := range esmDataAfterCoolOff.CollateralAsset {
// 		collAsset, _ := k.asset.GetAsset(ctx, data.AssetID)
// 		tokenDValue := data.Share.Mul(userWorth)
// 		price, _ := k.GetSnapshotOfPrices(ctx, appID, data.AssetID) // getting last saved prices
// 		oldTokenQuant := tokenDValue.Quo(sdk.NewDecFromInt(sdk.NewIntFromUint64(price)))
// 		usd := 1000000
// 		tokenQuant := oldTokenQuant.Quo(sdk.NewDec(int64(usd)))
// 		err1 := k.bank.SendCoinsFromAccountToModule(ctx, userAddress, types.ModuleName, sdk.NewCoins(amount))
// 		if err1 != nil {
// 			return err1
// 		}
// 		err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(collAsset.Denom, tokenQuant.TruncateInt())))
// 		if err != nil {
// 			return err
// 		}
// 		err2 := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
// 		if err2 != nil {
// 			return err2
// 		}
// 		data.Amount = data.Amount.Sub(tokenQuant.TruncateInt())
// 		esmDataAfterCoolOff.CollateralAsset = append(esmDataAfterCoolOff.CollateralAsset[:i], esmDataAfterCoolOff.CollateralAsset[i+1:]...)
// 		esmDataAfterCoolOff.CollateralAsset = append(esmDataAfterCoolOff.CollateralAsset[:i+1], esmDataAfterCoolOff.CollateralAsset[i:]...)
// 		esmDataAfterCoolOff.CollateralAsset[i] = data
// 		k.SetDataAfterCoolOff(ctx, esmDataAfterCoolOff)
// 	}

// 	for i, b := range esmDataAfterCoolOff.DebtAsset {
// 		if b.AssetID == assetInID.Id {
// 			b.Amount = b.Amount.Sub(amount.Amount)
// 			esmDataAfterCoolOff.DebtAsset = append(esmDataAfterCoolOff.DebtAsset[:i], esmDataAfterCoolOff.DebtAsset[i+1:]...)
// 			esmDataAfterCoolOff.DebtAsset = append(esmDataAfterCoolOff.DebtAsset, b)
// 			k.SetDataAfterCoolOff(ctx, esmDataAfterCoolOff)
// 		}
// 	}

// 	return nil
// }
