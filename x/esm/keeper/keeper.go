package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/esm/expected"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
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

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) DepositESM(ctx sdk.Context, depositorAddr string, AppID uint64, Amount sdk.Coin) error {

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
	govAsset, found := k.GetAsset(ctx, govTokenID)
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
	currentDeposit, found := k.GetCurrentDepositStats(ctx, AppID)
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
	newCurrentDeposit, _ := k.GetCurrentDepositStats(ctx, AppID)

	if newCurrentDeposit.Balance.Amount.GT(esmTriggerParams.TargetValue.Amount) {
		return types.ErrDepositForAppReached
	}
	if Amount.Amount.GT(esmTriggerParams.TargetValue.Amount) {
		return types.ErrAmtExceedsTargetValue
	}
	addr, _ := sdk.AccAddressFromBech32(depositorAddr)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, tokenminttypes.ModuleName, sdk.NewCoins(Amount)); err != nil {
		return err
	}
	if err1 := k.tokenmint.BurnTokensForApp(ctx, AppID, govAsset.Id, Amount.Amount); err1 != nil {
		return err1
	}

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

	_, found := k.GetApp(ctx, AppID)
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

	currentDeposit, _ := k.GetCurrentDepositStats(ctx, AppID)

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

func (k Keeper) CalculateCollateral(ctx sdk.Context, appId uint64, amount sdk.Coin, esmDataAfterCoolOff types.DataAfterCoolOff, from string) error {
	userAddress, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}
	marketData, found := k.GetESMMarketForAsset(ctx, appId)
	if !found {
		return types.ErrMarketDataNotFound
	}

	assetInID, _ := k.GetAssetForDenom(ctx, amount.Denom)
	var (
		assetInPrice uint64
		assetPrice   uint64
	)
	AppValue := sdk.ZeroInt()
	AppAssetValue := sdk.ZeroInt()
	for _, v := range marketData.Market {
		if assetInID.Id == v.AssetID {
			assetInPrice = v.Rates
		} else {
			assetInPrice = 0
		}
	}

	amtInPrice := amount.Amount.Mul(sdk.NewIntFromUint64(assetInPrice))

	for _, u := range esmDataAfterCoolOff.CollateralAsset {
		for _, w := range marketData.Market {
			if u.AssetID == w.AssetID {
				assetPrice = w.Rates
				AppAssetValue = u.Amount.Mul(sdk.NewIntFromUint64(assetPrice))
				assetVal := types.AssetToAmountValue{
					AppId:                     appId,
					AssetID:                   u.AssetID,
					Amount:                    AppAssetValue,
					AssetValueToAppValueRatio: sdk.ZeroDec(),
				}
				k.SetAssetToAmountValue(ctx, assetVal)
			}
			AppValue = AppValue.Add(AppAssetValue)
			appToAmtValue := types.AppToAmountValue{
				AppId:  appId,
				Amount: AppValue,
			}
			k.SetAppToAmtValue(ctx, appToAmtValue)
		}
	}
	for _, a := range esmDataAfterCoolOff.CollateralAsset {
		assetToAmountValue, _ := k.GetAssetToAmountValue(ctx, appId, a.AssetID)
		appToAmountValue, _ := k.GetAppToAmtValue(ctx, appId)
		Ratio := assetToAmountValue.Amount.Quo(appToAmountValue.Amount).ToDec()
		assetToAmountValue.AssetValueToAppValueRatio = Ratio
		k.SetAssetToAmountValue(ctx, assetToAmountValue)
		asset, _ := k.GetAsset(ctx, a.AssetID)
		factor1 := Ratio.Mul(sdk.Dec(amtInPrice)).Mul(sdk.Dec(amtInPrice))
		amountToDispatch := factor1.Quo(sdk.NewDec(int64(assetInPrice))).TruncateInt64()
		collateralTokens := sdk.NewCoin(asset.Denom, sdk.NewInt(amountToDispatch))

		err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(collateralTokens))
		if err != nil {
			return err
		}
		a.Amount = a.Amount.Sub(collateralTokens.Amount)
		esmDataAfterCoolOff.CollateralAsset = append(esmDataAfterCoolOff.CollateralAsset, a)
		k.SetDataAfterCoolOff(ctx, esmDataAfterCoolOff)
	}

	for _, b := range esmDataAfterCoolOff.DebtAsset {
		if b.AssetID == assetInID.Id{
			b.Amount = b.Amount.Sub(amount.Amount)
			esmDataAfterCoolOff.DebtAsset = append(esmDataAfterCoolOff.DebtAsset, b)
			k.SetDataAfterCoolOff(ctx, esmDataAfterCoolOff)
		}
	}

	return nil
}
