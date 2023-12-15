package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/esm/expected"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		asset      expected.AssetKeeper
		vault      expected.VaultKeeper
		bank       expected.BankKeeper
		market     expected.MarketKeeper
		tokenmint  expected.Tokenmint
		collector  expected.Collector
		authority  string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	bank expected.BankKeeper,
	market expected.MarketKeeper,
	tokenmint expected.Tokenmint,
	collector expected.Collector,
	authority string,
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
		authority:  authority,
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

// take token input from user
// find the token and its per token dollar value
// calculate percentage of dollar of each collateral token
// calculate token quantity from dollar value
// subtract total dollar value from the end from the struct,
// subtract token quantity from each collateral
func (k Keeper) CalculateCollateral(ctx sdk.Context, appID uint64, amount sdk.Coin, from string) error {
	userAddress, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}
	assetID, _ := k.asset.GetAssetForDenom(ctx, amount.Denom)
	coolOffData, found := k.GetDataAfterCoolOff(ctx, appID)
	if !found {
		return types.ErrAppDataNotFound
	}
	assetESMData, found := k.GetAssetToAmount(ctx, appID, assetID.Id)
	if !found || assetESMData.IsCollateral || assetESMData.Amount.IsZero() || amount.Amount.GT(assetESMData.Amount) {
		return types.ErrInvalidAsset
	}
	unitWorth := assetESMData.DebtTokenWorth
	// Total worth of debt asset brought by user
	// totalDebtAssetWorth := unitWorth.Mul(amount.Amount.ToDec())
	totalDebtAssetWorth := k.CalcDollarValueOfToken(ctx, unitWorth.TruncateInt().Uint64(), amount.Amount, assetID.Decimals)
	err1 := k.bank.SendCoinsFromAccountToModule(ctx, userAddress, types.ModuleName, sdk.NewCoins(amount))
	if err1 != nil {
		return err1
	}
	err2 := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err2 != nil {
		return err2
	}

	// Calculating share of all collateral asset that needs to be paid to the user upto the $ value of totalDebtAssetWorth
	// Eg. If user brought $20 of CMST (totalDebtAssetWorth) and if 2 collateral exist as CMDX(80%) and ATOM(20%) then 20% of tokens that
	// will be paid will be in terms of ATOM tokens , rest CMDX.
	allAssetToAmtData := k.GetAllAssetToAmount(ctx, appID)
	for _, tokenData := range allAssetToAmtData {
		assetData, found := k.asset.GetAsset(ctx, tokenData.AssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}

		if tokenData.IsCollateral && !tokenData.Amount.IsZero() {
			// unitRate, _ := k.GetSnapshotOfPrices(ctx, appID, assetData.Id)
			unitRate := k.GetRateOfAsset(ctx, appID, assetData.Id)
			if unitRate == 0 {
				unitRate, found = k.GetSnapshotOfPrices(ctx, appID, assetData.Id)
				if !found {
					return types.ErrPriceNotFound
				}
			}
			tokenShare := totalDebtAssetWorth.Mul(tokenData.Share) //$CMST Multiplied with Share of collateral give $share of collateral
			// To calculate quantity of collateral token from the $share of tokenShare
			collateralQuantity := tokenShare.Quo(sdk.NewDecFromInt(sdk.NewIntFromUint64(unitRate)))
			collateralQuantity = collateralQuantity.Mul(sdk.NewDecFromInt(assetData.Decimals))
			err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(assetData.Denom, collateralQuantity.TruncateInt())))
			if err != nil {
				return err
			}
			// Reducing quantity of token in collateral
			tokenData.Amount = tokenData.Amount.Sub(collateralQuantity.TruncateInt())
		}
		k.SetAssetToAmount(ctx, tokenData)
	}

	// To reduce collateral token dollar quantity

	// As the worth was calculated against collateral asset - totalDebtAssetWorth so we use it to subtract collateral dollar value
	coolOffData.CollateralTotalAmount = coolOffData.CollateralTotalAmount.Sub(totalDebtAssetWorth)

	// To reduce Debt token dollar quantity
	// We need to now find the debt token worth,to subtract it

	rateOut := k.GetRateOfAsset(ctx, appID, assetID.Id)
	if rateOut == 0 {
		rateOut, found = k.GetSnapshotOfPrices(ctx, appID, assetID.Id)
		if !found {
			return types.ErrPriceNotFound
		}
	}
	totalPrice := k.CalcDollarValueOfToken(ctx, rateOut, amount.Amount, assetID.Decimals)
	coolOffData.DebtTotalAmount = coolOffData.DebtTotalAmount.Sub(totalPrice)

	// To reduce the debt token amount from asset data.

	assetESMData.Amount = assetESMData.Amount.Sub(amount.Amount)
	k.SetAssetToAmount(ctx, assetESMData)

	k.SetDataAfterCoolOff(ctx, coolOffData)

	return nil
}
