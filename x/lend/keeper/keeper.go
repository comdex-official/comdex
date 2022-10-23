package keeper

import (
	"fmt"
	"strconv"

	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/lend/expected"
	"github.com/comdex-official/comdex/x/lend/types"
	// liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    sdk.StoreKey
		memKey      sdk.StoreKey
		paramstore  paramtypes.Subspace
		bank        expected.BankKeeper
		account     expected.AccountKeeper
		asset       expected.AssetKeeper
		market      expected.MarketKeeper
		esm         expected.EsmKeeper
		liquidation expected.LiquidationKeeper
		auction     expected.AuctionKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
	account expected.AccountKeeper,
	asset expected.AssetKeeper,
	market expected.MarketKeeper,
	esm expected.EsmKeeper,
	liquidation expected.LiquidationKeeper,
	auction expected.AuctionKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		memKey:      memKey,
		paramstore:  ps,
		bank:        bank,
		account:     account,
		asset:       asset,
		market:      market,
		esm:         esm,
		liquidation: liquidation,
		auction:     auction,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func uint64InAssetData(a uint64, list []*types.AssetDataPoolMapping) bool {
	for _, b := range list {
		if b.AssetID == a {
			return true
		}
	}
	return false
}

func (k Keeper) CheckSupplyCap(ctx sdk.Context, assetID, poolID uint64, amt sdk.Int) (bool, error) {
	// this fn checks if while depositing the supply cap of a specific lend asset doesn't exceed specified value
	var supplyCap uint64
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)

	currentSupply, err := k.market.CalcAssetPrice(ctx, assetID, assetStats.TotalLend.Add(amt))
	if err != nil {
		return false, err
	}
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return false, types.ErrPoolNotFound
	}

	for _, v := range pool.AssetData {
		if assetID == v.AssetID {
			supplyCap = v.SupplyCap
		}
	}
	if currentSupply.LTE(sdk.NewDecFromInt(sdk.NewIntFromUint64(supplyCap))) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr string, AssetID uint64, Amount sdk.Coin, PoolID, AppID uint64) error {
	// this fn IBC Assets fom the user
	// sends the asset to pool's module-acc
	// mints cAsset representative of the lent asset
	// creates a lent Position and updates global lend

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.asset.GetAsset(ctx, AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	appMapping, found := k.asset.GetApp(ctx, AppID)
	if !found {
		return types.ErrorAppMappingDoesNotExist
	}
	if appMapping.Name != types.AppName {
		return types.ErrorAppMappingIDMismatch
	}

	if Amount.Denom != asset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, Amount.Denom)
	}

	found = uint64InAssetData(AssetID, pool.AssetData)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidAssetIDForPool, strconv.FormatUint(AssetID, 10))
	}

	found, err := k.CheckSupplyCap(ctx, AssetID, PoolID, Amount.Amount)
	if err != nil {
		return err
	}
	if !found {
		return types.ErrorSupplyCapExceeds
	}

	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if k.HasLendForAddressByAsset(ctx, lenderAddr, AssetID, PoolID) {
		return types.ErrorDuplicateLend
	}

	loanTokens := sdk.NewCoins(Amount)

	assetRatesStat, found := k.GetAssetRatesParams(ctx, AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(AssetID, 10))
	}
	cAsset, found := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)

	if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, loanTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err = k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
	if err != nil {
		return err
	}

	lendID := k.GetUserLendIDCounter(ctx)

	var globalIndex sdk.Dec
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, PoolID, AssetID)
	if assetStats.LendApr.IsZero() {
		globalIndex = sdk.OneDec()
	} else {
		globalIndex = assetStats.LendApr
	}

	lendPos := types.LendAsset{
		ID:                  lendID + 1,
		AssetID:             AssetID,
		PoolID:              PoolID,
		Owner:               lenderAddr,
		AmountIn:            Amount,
		LendingTime:         ctx.BlockTime(),
		AvailableToBorrow:   Amount.Amount,
		AppID:               AppID,
		GlobalIndex:         globalIndex,
		LastInteractionTime: ctx.BlockTime(),
		CPoolName:           pool.CPoolName,
	}
	k.UpdateLendStats(ctx, AssetID, PoolID, Amount.Amount, true) // update global lend data in poolAssetLBMappingData
	k.SetUserLendIDCounter(ctx, lendPos.ID)
	k.SetLend(ctx, lendPos)

	// making UserAssetLendBorrowMapping for user
	var mappingData types.UserAssetLendBorrowMapping
	mappingData.Owner = lendPos.Owner
	mappingData.LendId = lendPos.ID
	mappingData.PoolId = PoolID
	mappingData.BorrowId = nil
	k.SetUserLendBorrowMapping(ctx, mappingData)

	// Adding Lend ID mapping to poolAssetLBMappingData
	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, PoolID, AssetID)
	poolAssetLBMappingData.LendIds = append(poolAssetLBMappingData.LendIds, lendPos.ID)
	k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, addr string, lendID uint64, withdrawal sdk.Coin) error {
	// this fn is used to withdraw IBC assets from previously lent asset
	// we take the cTokens from the user's account and burn them.
	// equal amount of IBC Asset is sent back to the user.
	// global lend data and user's lend position are updated after successful transaction

	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	// if user wants to withdraw all amount then his position will be closed
	if withdrawal.Amount.Equal(lendPos.AvailableToBorrow) && lendPos.AvailableToBorrow.GTE(lendPos.AmountIn.Amount) {
		err = k.CloseLend(ctx, addr, lendID)
		if err != nil {
			return err
		}
		return nil
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()

	getAsset, _ := k.asset.GetAsset(ctx, lendPos.AssetID)
	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}

	if withdrawal.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrWithdrawAmountLimitExceeds
	}

	if withdrawal.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, withdrawal.Denom)
	}

	reservedAmount := k.GetReserveFunds(ctx, pool)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, withdrawal.Denom)

	if withdrawal.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, withdrawal.String())
	}

	assetRatesStat, found := k.GetAssetRatesParams(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, withdrawal.Amount)
	cTokens := sdk.NewCoins(cToken)

	tokens := sdk.NewCoins(withdrawal)
	if err != nil {
		return err
	}

	// here 2 conditions exists
	// a. if the amount to withdraw is less than equal the AmountIn of lend position then Available to borrow are subtracted by withdraw amt
	// b. if the user has some lend rewards added to his position and the withdrawal amount is greater tha amountIn then in this case the amountIn is made 0
	// as no further lend rewards will be available for this user.

	if withdrawal.Amount.LT(lendPos.AmountIn.Amount) {
		if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(cToken)); err != nil {
			return err
		}
		// burn c/Token
		err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
		if err != nil {
			return err
		}

		if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
			return err
		}
		k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, withdrawal.Amount, false)
		lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(withdrawal.Amount)
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
		k.SetLend(ctx, lendPos)
	} else {
		if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(cToken)); err != nil {
			return err
		}
		// burn c/Token
		err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
		if err != nil {
			return err
		}

		if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
			return err
		}

		k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, withdrawal.Amount, false)
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
		lendPos.AmountIn.Amount = sdk.ZeroInt()
		k.SetLend(ctx, lendPos)
	}
	return nil
}

func (k Keeper) DepositAsset(ctx sdk.Context, addr string, lendID uint64, deposit sdk.Coin) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()
	// checking if the caller is the owner to Lend Position
	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}
	getAsset, _ := k.asset.GetAsset(ctx, lendPos.AssetID)
	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	found, err = k.CheckSupplyCap(ctx, lendPos.AssetID, lendPos.PoolID, deposit.Amount)
	if err != nil {
		return err
	}
	if !found {
		return types.ErrorSupplyCapExceeds
	}

	if deposit.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, deposit.Denom)
	}

	assetRatesStat, found := k.GetAssetRatesParams(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, deposit.Amount)

	cTokens := sdk.NewCoins(cToken)

	if err = k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(deposit)); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, cTokens)
	if err != nil {
		return err
	}

	lendPos.AmountIn = lendPos.AmountIn.Add(deposit)
	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(deposit.Amount)

	k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, deposit.Amount, true)
	k.SetLend(ctx, lendPos)
	return nil
}

func (k Keeper) CloseLend(ctx sdk.Context, addr string, lendID uint64) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()

	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}

	lendIDToBorrowIDMapping, _ := k.GetUserLendBorrowMapping(ctx, lendPos.Owner, lendID)
	if lendIDToBorrowIDMapping.BorrowId != nil {
		return types.ErrBorrowingPositionOpen
	}
	reservedAmount := k.GetReserveFunds(ctx, pool)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, lendPos.AmountIn.Denom)

	if lendPos.AvailableToBorrow.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, lendPos.AvailableToBorrow.String())
	}

	tokens := sdk.NewCoins(sdk.NewCoin(lendPos.AmountIn.Denom, lendPos.AvailableToBorrow))
	assetRatesStat, found := k.GetAssetRatesParams(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, lendPos.AvailableToBorrow)

	if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(cToken)); err != nil {
		return err
	}

	cTokens := sdk.NewCoins(cToken)
	err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
	if err != nil {
		return err
	}

	if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
		return err
	}

	k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, lendPos.AvailableToBorrow, false)
	k.DeleteLendForAddressByAsset(ctx, lenderAddr.String(), lendPos.ID)

	k.DeleteIDFromAssetStatsMapping(ctx, lendPos.PoolID, lendPos.AssetID, lendID, true)

	k.DeleteLend(ctx, lendPos.ID)

	return nil
}

func uint64InSlice(a uint64, list []uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//nolint:funlen
func (k Keeper) BorrowAsset(ctx sdk.Context, addr string, lendID, pairID uint64, IsStableBorrow bool, AmountIn, loan sdk.Coin) error {
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}

	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	pairMapping, _ := k.GetAssetToPair(ctx, pair.AssetIn, lendPos.PoolID)
	found = uint64InSlice(pairID, pairMapping.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	// check cr ratio
	assetIn, found := k.asset.GetAsset(ctx, lendPos.AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetOut, found := k.asset.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetInRatesStats, found := k.GetAssetRatesParams(ctx, pair.AssetIn)
	if !found {
		return types.ErrAssetStatsNotFound
	}

	cAsset, found := k.asset.GetAsset(ctx, assetInRatesStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if AmountIn.Denom != cAsset.Denom {
		return types.ErrBadOfferCoinType
	}

	if k.HasBorrowForAddressByPair(ctx, addr, pairID) {
		return types.ErrorDuplicateBorrow
	}

	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}

	if loan.Denom != assetOut.Denom {
		return types.ErrInvalidAsset
	}

	AssetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	AssetOutPool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}

	if IsStableBorrow && !assetInRatesStats.EnableStableBorrow {
		return sdkerrors.Wrap(types.ErrStableBorrowDisabled, loan.String())
	}

	err := k.VerifyCollateralizationRatio(ctx, AmountIn.Amount, assetIn, loan.Amount, assetOut, assetInRatesStats.Ltv)
	if err != nil {
		return err
	}
	borrowID := k.GetUserBorrowIDCounter(ctx)

	reservedAmount := k.GetReserveFunds(ctx, AssetOutPool)
	availableAmount := k.ModuleBalance(ctx, AssetOutPool.ModuleName, loan.Denom)
	// check sufficient amt in pool to borrow
	if loan.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrBorrowingPoolInsufficient, loan.String())
	}
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
	reserveGlobalIndex, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		reserveGlobalIndex = sdk.OneDec()
	}
	globalIndex := assetStats.BorrowApr
	AmountOut := loan
	// There are 3 possible cases of borrowing
	// a. When the borrow is done from the same pool in which the user has lent Asset
	// b. When the borrow is from different pool using First Transit Asset
	// c. When the borrow is from different pool using Second Transit Asset
	// else the borrowing pool is not having sufficient tokens to loan
	if !pair.IsInterPool {
		// take c/Tokens from the user
		if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
			return err
		}

		if err = k.bank.SendCoinsFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, sdk.NewCoins(loan)); err != nil {
			return err
		}

		var StableBorrowRate sdk.Dec
		if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
			StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
			if err != nil {
				return err
			}
		} else {
			StableBorrowRate = sdk.ZeroDec()
		}

		borrowPos := types.BorrowAsset{
			ID:                  borrowID + 1,
			LendingID:           lendID,
			PairID:              pairID,
			AmountIn:            AmountIn,
			AmountOut:           AmountOut,
			BridgedAssetAmount:  sdk.NewCoin(loan.Denom, sdk.NewInt(0)),
			IsStableBorrow:      IsStableBorrow,
			StableBorrowRate:    StableBorrowRate,
			BorrowingTime:       ctx.BlockTime(),
			InterestAccumulated: sdk.ZeroDec(),
			GlobalIndex:         globalIndex,
			ReserveGlobalIndex:  reserveGlobalIndex,
			LastInteractionTime: ctx.BlockTime(),
			CPoolName:           AssetOutPool.CPoolName,
			IsLiquidated:        false,
		}
		k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, AmountOut.Amount, true)
		// err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
		// if err != nil {
		// 	return err
		// }
		poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
		poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, borrowPos.ID)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)
		k.SetUserBorrowIDCounter(ctx, borrowPos.ID)
		k.SetBorrow(ctx, borrowPos)

		mappingData, _ := k.GetUserLendBorrowMapping(ctx, addr, lendID)
		mappingData.BorrowId = append(mappingData.BorrowId, borrowPos.ID)
		k.SetUserLendBorrowMapping(ctx, mappingData)
	} else {
		updatedAmtIn := AmountIn.Amount.ToDec().Mul(assetInRatesStats.Ltv)
		updatedAmtInPrice, err := k.market.CalcAssetPrice(ctx, lendPos.AssetID, updatedAmtIn.TruncateInt())
		if err != nil {
			return err
		}
		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range AssetInPool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}
		firstTransitAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
		secondTransitAsset, _ := k.asset.GetAsset(ctx, secondTransitAssetID)

		unitAmountFirstTransitAsset, err := k.market.CalcAssetPrice(ctx, firstTransitAssetID, sdk.OneInt())
		if err != nil {
			return err
		}
		unitAmountSecondTransitAsset, err := k.market.CalcAssetPrice(ctx, secondTransitAssetID, sdk.OneInt())
		if err != nil {
			return err
		}

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		firstBridgedAssetQty := updatedAmtInPrice.Quo(unitAmountFirstTransitAsset)
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstTransitAsset.Denom)
		secondBridgedAssetQty := updatedAmtInPrice.Quo(unitAmountSecondTransitAsset)
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondTransitAsset.Denom)

		firstBridgedAssetRatesStats, found := k.GetAssetRatesParams(ctx, firstTransitAsset.Id)
		if !found {
			return types.ErrAssetStatsNotFound
		}
		secondBridgedAssetRatesStats, found := k.GetAssetRatesParams(ctx, secondTransitAsset.Id)
		if !found {
			return types.ErrAssetStatsNotFound
		}
		if firstBridgedAssetQty.LT(sdk.NewDecFromInt(firstBridgedAssetBal)) {
			err = k.VerifyCollateralizationRatio(ctx, firstBridgedAssetQty.TruncateInt(), firstTransitAsset, loan.Amount, assetOut, firstBridgedAssetRatesStats.Ltv)
			if err != nil {
				return err
			}
			// take c/Tokens from the user
			if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
				return err
			}
			bridgedAssetAmount := sdk.NewCoin(firstTransitAsset.Denom, firstBridgedAssetQty.TruncateInt())
			if err = k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, sdk.NewCoins(loan)); err != nil {
				return err
			}

			var StableBorrowRate sdk.Dec
			if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
				StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
				if err != nil {
					return err
				}
			} else {
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                  borrowID + 1,
				LendingID:           lendID,
				PairID:              pairID,
				AmountIn:            AmountIn,
				AmountOut:           AmountOut,
				BridgedAssetAmount:  bridgedAssetAmount,
				IsStableBorrow:      IsStableBorrow,
				StableBorrowRate:    StableBorrowRate,
				BorrowingTime:       ctx.BlockTime(),
				InterestAccumulated: sdk.ZeroDec(),
				GlobalIndex:         globalIndex,
				ReserveGlobalIndex:  reserveGlobalIndex,
				LastInteractionTime: ctx.BlockTime(),
				CPoolName:           AssetOutPool.CPoolName,
				IsLiquidated:        false,
			}
			k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, AmountOut.Amount, true)

			poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, borrowPos.ID)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDCounter(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)
			mappingData, _ := k.GetUserLendBorrowMapping(ctx, addr, lendID)
			mappingData.BorrowId = append(mappingData.BorrowId, borrowPos.ID)
			k.SetUserLendBorrowMapping(ctx, mappingData)
		} else if secondBridgedAssetQty.LT(sdk.NewDecFromInt(secondBridgedAssetBal)) {
			err = k.VerifyCollateralizationRatio(ctx, secondBridgedAssetQty.TruncateInt(), secondTransitAsset, loan.Amount, assetOut, secondBridgedAssetRatesStats.Ltv)
			if err != nil {
				return err
			}
			// take c/Tokens from the user
			if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
				return err
			}

			bridgedAssetAmount := sdk.NewCoin(secondTransitAsset.Denom, secondBridgedAssetQty.TruncateInt())
			if err = k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, sdk.NewCoins(loan)); err != nil {
				return err
			}

			var StableBorrowRate sdk.Dec
			if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
				StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
				if err != nil {
					return err
				}
			} else {
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                  borrowID + 1,
				LendingID:           lendID,
				PairID:              pairID,
				AmountIn:            AmountIn,
				AmountOut:           AmountOut,
				BridgedAssetAmount:  bridgedAssetAmount,
				IsStableBorrow:      IsStableBorrow,
				StableBorrowRate:    StableBorrowRate,
				BorrowingTime:       ctx.BlockTime(),
				InterestAccumulated: sdk.ZeroDec(),
				GlobalIndex:         globalIndex,
				ReserveGlobalIndex:  reserveGlobalIndex,
				LastInteractionTime: ctx.BlockTime(),
				CPoolName:           AssetOutPool.CPoolName,
				IsLiquidated:        false,
			}
			k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, AmountOut.Amount, true)

			poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, borrowPos.ID)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDCounter(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)

			mappingData, _ := k.GetUserLendBorrowMapping(ctx, addr, lendID)
			mappingData.BorrowId = append(mappingData.BorrowId, borrowPos.ID)
			k.SetUserLendBorrowMapping(ctx, mappingData)
		} else {
			return types.ErrBorrowingPoolInsufficient
		}
	}
	return nil
}

func (k Keeper) RepayAsset(ctx sdk.Context, borrowID uint64, borrowerAddr string, payment sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	if borrowPos.IsLiquidated {
		return types.ErrorBorrowPosLiquidated
	}

	if payment.Amount.Equal(borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt())) {
		err := k.CloseBorrow(ctx, borrowerAddr, borrowID)
		if err != nil {
			return err
		}
		return nil
	}

	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	assetStats, found := k.GetAssetRatesParams(ctx, pair.AssetOut)
	if !found {
		return types.ErrAssetStatsNotFound
	}
	cAsset, found := k.asset.GetAsset(ctx, assetStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorized
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	if borrowPos.AmountOut.Denom != payment.Denom {
		return types.ErrBadOfferCoinAmount
	}
	if payment.Amount.GT(borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.Ceil().RoundInt())) {
		return types.ErrInvalidRepayment
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)

	reservePoolRecords, _ := k.GetBorrowInterestTracker(ctx, borrowID)
	amtToReservePool := reservePoolRecords.ReservePoolInterest

	if payment.Amount.LTE(amtToReservePool.TruncateInt()) {
		// subtract payment amount from reservePoolRecords for that borrow ID and set reservePoolRecords
		// subtract payment amount from borrow-> interest accumulated and set borrow
		// send the payment amount to the reserve pool & update stats

		// sending repayment to moduleAcc from borrower
		if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
		reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Sub(sdk.NewDecFromInt(payment.Amount))
		borrowPos.InterestAccumulated = borrowPos.InterestAccumulated.Sub(sdk.NewDecFromInt(payment.Amount))
		err = k.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, payment, true) // this functions sends the amount to reserve pool from the cPool
		if err != nil {
			return err
		}
	} else if payment.Amount.GT(amtToReservePool.TruncateInt()) && payment.Amount.LTE(borrowPos.InterestAccumulated.TruncateInt()) {
		// from reservePoolRecords amount send tokens to reserve pool
		// send remaining payment back to cPool and mint additional tokens for that amount
		// update global lend and interest accumulated for remaining amount
		// subtract payment amount from borrow-> interest accumulated and set borrow
		// make reservePoolRecords 0

		// sending repayment to moduleAcc from borrower
		if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}

		err = k.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, sdk.NewCoin(payment.Denom, amtToReservePool.TruncateInt()), true) // this functions sends the amount to reserve pool from the cPool
		if err != nil {
			return err
		}

		// calculation for tokens to be minted and updated in global lend and interest accumulated parameter
		cTokensAmount := payment.Amount.Sub(amtToReservePool.TruncateInt())
		if cTokensAmount.LT(sdk.ZeroInt()) {
			return types.ErrReserveRatesNotFound
		}
		if cTokensAmount.GT(sdk.ZeroInt()) {
			err = k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, cTokensAmount)))
			if err != nil {
				return err
			}
			poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Add(cTokensAmount)
			poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Add(cTokensAmount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
		}

		reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Sub(sdk.NewDecFromInt(amtToReservePool.TruncateInt())) // the decimal precision is maintained
		borrowPos.InterestAccumulated = borrowPos.InterestAccumulated.Sub(sdk.NewDecFromInt(payment.Amount))
	} else { // if payment amt greater than borrowPos.InterestAccumulated
		// from reservePoolRecords amount send tokens to reserve pool
		// mint (interest Accumulated - reservePoolRecords) amount cTokens and send to pool
		// update global lend and interest accumulated by the above amount
		// send (payment amount -  reservePoolRecords amount ) back to the pool
		// subtract (payment -interest accumulated) from amount out and set borrow
		// update borrow Stats
		// make reservePoolRecords 0
		// make interest accumulated 0

		// sending repayment to moduleAcc from borrower
		if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}

		err = k.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, sdk.NewCoin(payment.Denom, amtToReservePool.TruncateInt()), true) // this functions sends the amount to reserve pool from the cPool
		if err != nil {
			return err
		}

		// calculation for tokens to be minted and updated in global lend and interest accumulated parameter
		cTokensAmount := borrowPos.InterestAccumulated.TruncateInt().Sub(reservePoolRecords.ReservePoolInterest.TruncateInt())
		if cTokensAmount.LT(sdk.ZeroInt()) {
			return types.ErrReserveRatesNotFound
		}
		if cTokensAmount.GT(sdk.ZeroInt()) {
			err = k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, cTokensAmount)))
			if err != nil {
				return err
			}
			poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Add(cTokensAmount)
			poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Add(cTokensAmount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
		}
		reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Sub(sdk.NewDecFromInt(amtToReservePool.TruncateInt())) // the decimal precision is maintained
		amtToSubFromBorrowPos := payment.Amount.Sub(borrowPos.InterestAccumulated.TruncateInt())
		borrowPos.AmountOut.Amount = borrowPos.AmountOut.Amount.Sub(amtToSubFromBorrowPos)
		borrowPos.InterestAccumulated = borrowPos.InterestAccumulated.Sub(sdk.NewDecFromInt(borrowPos.InterestAccumulated.TruncateInt()))
		k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, amtToSubFromBorrowPos, false)
	}

	k.SetBorrow(ctx, borrowPos)
	k.SetBorrowInterestTracker(ctx, reservePoolRecords)

	return nil
}

func (k Keeper) DepositBorrowAsset(ctx sdk.Context, borrowID uint64, addr string, AmountIn sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}

	if borrowPos.IsLiquidated {
		return types.ErrorBorrowPosLiquidated
	}

	lendID := borrowPos.LendingID
	pairID := borrowPos.PairID
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	assetRatesStat, found := k.GetAssetRatesParams(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, found := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if AmountIn.Denom != cAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
	}
	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}
	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	AssetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	AssetOutPool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}

	if !pair.IsInterPool {
		if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
			return err
		}
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)
		borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
		k.SetBorrow(ctx, borrowPos)
	} else {
		amtIn, err := k.market.CalcAssetPrice(ctx, pair.AssetIn, (sdk.NewDecFromInt(AmountIn.Amount).Mul(assetRatesStat.Ltv)).TruncateInt())
		if err != nil {
			return err
		}

		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range AssetInPool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}
		firstTransitAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
		secondTransitAsset, _ := k.asset.GetAsset(ctx, secondTransitAssetID)

		unitAmountFirstTransitAsset, err := k.market.CalcAssetPrice(ctx, firstTransitAssetID, sdk.OneInt())
		if err != nil {
			return err
		}
		unitAmountSecondTransitAsset, err := k.market.CalcAssetPrice(ctx, secondTransitAssetID, sdk.OneInt())
		if err != nil {
			return err
		}

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		firstBridgedAssetQty := amtIn.Quo(unitAmountFirstTransitAsset)
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstTransitAsset.Denom)
		secondBridgedAssetQty := amtIn.Quo(unitAmountSecondTransitAsset)
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondTransitAsset.Denom)

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		if borrowPos.BridgedAssetAmount.Denom == firstTransitAsset.Denom && firstBridgedAssetQty.LT(sdk.NewDecFromInt(firstBridgedAssetBal)) {
			// take c/Tokens from the user
			if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(firstTransitAsset.Denom, firstBridgedAssetQty.TruncateInt()))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			borrowPos.BridgedAssetAmount.Amount = borrowPos.BridgedAssetAmount.Amount.Add(firstBridgedAssetQty.TruncateInt())
			k.SetBorrow(ctx, borrowPos)
		} else if secondBridgedAssetQty.LT(sdk.NewDecFromInt(secondBridgedAssetBal)) {
			// take c/Tokens from the user
			if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, sdk.NewCoins(AmountIn)); err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(secondTransitAsset.Denom, secondBridgedAssetQty.TruncateInt()))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			borrowPos.BridgedAssetAmount.Amount = borrowPos.BridgedAssetAmount.Amount.Add(secondBridgedAssetQty.TruncateInt())
			k.SetBorrow(ctx, borrowPos)
		} else {
			return types.ErrBridgeAssetQtyInsufficient
		}
	}
	return nil
}

func (k Keeper) DrawAsset(ctx sdk.Context, borrowID uint64, borrowerAddr string, amount sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}

	if borrowPos.IsLiquidated {
		return types.ErrorBorrowPosLiquidated
	}

	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorized
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	if borrowPos.AmountOut.Denom != amount.Denom {
		return types.ErrBadOfferCoinAmount
	}
	assetIn, found := k.asset.GetAsset(ctx, lendPos.AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetOut, found := k.asset.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetRatesStats, found := k.GetAssetRatesParams(ctx, pair.AssetIn)
	if !found {
		return types.ErrorAssetStatsNotFound
	}
	err = k.VerifyCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()).Add(amount.Amount), assetOut, assetRatesStats.Ltv)
	if err != nil {
		return err
	}
	if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, sdk.NewCoins(amount)); err != nil {
		return err
	}
	borrowPos.AmountOut = borrowPos.AmountOut.Add(amount)
	k.SetBorrow(ctx, borrowPos)
	k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, amount.Amount, true)

	return nil
}

func (k Keeper) CloseBorrow(ctx sdk.Context, borrowerAddr string, borrowID uint64) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	if borrowPos.IsLiquidated {
		return types.ErrorBorrowPosLiquidated
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	assetStats, found := k.GetAssetRatesParams(ctx, pair.AssetOut)
	if !found {
		return types.ErrAssetStatsNotFound
	}
	cAsset, found := k.asset.GetAsset(ctx, assetStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorized
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	assetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	assetOut, found := k.asset.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	lenderAddr, _ := sdk.AccAddressFromBech32(lendPos.Owner)
	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)

	// take final Amt from the user (AmountOut + interest Accumulated) and send to cPool
	// return cToken Back to user
	// send amt from reservePoolRecords to reserve pool
	// mint cTokens from the (InterestAccumulated - Amnt To reserve pool) amt
	// update global lend and interest accumulated from poolAssetLBMappingData
	// delete borrow
	// delete reservePoolRecords for that borrow ID

	amt := sdk.NewCoins(sdk.NewCoin(assetOut.Denom, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt())))
	if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, amt); err != nil {
		return err
	}
	if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, lenderAddr, sdk.NewCoins(borrowPos.AmountIn)); err != nil {
		return err
	}

	reservePoolRecords, _ := k.GetBorrowInterestTracker(ctx, borrowID)
	amtToReservePool := reservePoolRecords.ReservePoolInterest
	if amtToReservePool.TruncateInt().LT(sdk.ZeroInt()) {
		return types.ErrReserveRatesNotFound
	}
	if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
		amount := sdk.NewCoin(assetOut.Denom, amtToReservePool.TruncateInt())
		err = k.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, amount, true)
		if err != nil {
			return err
		}
	}
	amtToMint := borrowPos.InterestAccumulated.TruncateInt().Sub(amtToReservePool.TruncateInt())
	if amtToMint.GT(sdk.ZeroInt()) {
		err = k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, amtToMint)))
		if err != nil {
			return err
		}
		poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Add(amtToMint)
		poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Add(amtToMint)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
	}

	if pair.IsInterPool {
		if err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, assetInPool.ModuleName, sdk.NewCoins(borrowPos.BridgedAssetAmount)); err != nil {
			return err
		}
	}

	k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)

	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(borrowPos.AmountIn.Amount)
	k.SetLend(ctx, lendPos)
	k.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, borrowID, false)
	k.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, borrowID)
	k.DeleteBorrow(ctx, borrowID)
	k.DeleteBorrowInterestTracker(ctx, borrowID)

	return nil
}

func (k Keeper) BorrowAlternate(ctx sdk.Context, lenderAddr string, AssetID, PoolID uint64, AmountIn sdk.Coin, PairID uint64, IsStableBorrow bool, AmountOut sdk.Coin, AppID uint64) error {
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.asset.GetAsset(ctx, AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	appMapping, found := k.asset.GetApp(ctx, AppID)
	if !found {
		return types.ErrorAppMappingDoesNotExist
	}
	if appMapping.Name != types.AppName {
		return types.ErrorAppMappingIDMismatch
	}

	if AmountIn.Denom != asset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
	}

	found = uint64InAssetData(AssetID, pool.AssetData)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidAssetIDForPool, strconv.FormatUint(AssetID, 10))
	}

	found, err := k.CheckSupplyCap(ctx, AssetID, PoolID, AmountIn.Amount)
	if err != nil {
		return err
	}
	if !found {
		return types.ErrorSupplyCapExceeds
	}

	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if k.HasLendForAddressByAsset(ctx, lenderAddr, AssetID, PoolID) {
		return types.ErrorDuplicateLend
	}

	loanTokens := sdk.NewCoins(AmountIn)

	assetRatesStat, found := k.GetAssetRatesParams(ctx, AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(AssetID, 10))
	}
	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, AmountIn.Amount)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, loanTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err := k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
	if err != nil {
		return err
	}

	lendID := k.GetUserLendIDCounter(ctx)

	var globalIndex sdk.Dec
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, PoolID, AssetID)
	if assetStats.LendApr.IsZero() {
		globalIndex = sdk.OneDec()
	} else {
		globalIndex = assetStats.LendApr
	}

	lendPos := types.LendAsset{
		ID:                  lendID + 1,
		AssetID:             AssetID,
		PoolID:              PoolID,
		Owner:               lenderAddr,
		AmountIn:            AmountIn,
		LendingTime:         ctx.BlockTime(),
		AvailableToBorrow:   AmountIn.Amount,
		AppID:               AppID,
		GlobalIndex:         globalIndex,
		LastInteractionTime: ctx.BlockTime(),
		CPoolName:           pool.CPoolName,
	}
	k.UpdateLendStats(ctx, AssetID, PoolID, AmountIn.Amount, true)
	k.SetUserLendIDCounter(ctx, lendPos.ID)
	k.SetLend(ctx, lendPos)

	var mappingData types.UserAssetLendBorrowMapping
	mappingData.Owner = lendPos.Owner
	mappingData.LendId = lendPos.ID
	mappingData.PoolId = PoolID
	k.SetUserLendBorrowMapping(ctx, mappingData)

	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, PoolID, AssetID)
	poolAssetLBMappingData.LendIds = append(poolAssetLBMappingData.LendIds, lendPos.ID)
	k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

	err = k.BorrowAsset(ctx, lenderAddr, lendPos.ID, PairID, IsStableBorrow, cToken, AmountOut)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) FundModAcc(ctx sdk.Context, poolID, assetID uint64, lenderAddr sdk.AccAddress, payment sdk.Coin) error {
	loanTokens := sdk.NewCoins(payment)
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrPoolNotFound
	}
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, loanTokens); err != nil {
		return err
	}

	asset, found := k.asset.GetAsset(ctx, assetID)
	if !found {
		return types.ErrLendNotFound
	}

	if asset.Denom != payment.Denom {
		return types.ErrBadOfferCoinType
	}

	assetRatesStat, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesParamsNotFound, strconv.FormatUint(assetID, 10))
	}
	cAsset, found := k.asset.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	cToken := sdk.NewCoin(cAsset.Denom, payment.Amount)

	err := k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(cToken))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault) {
	kind := liqBorrow.GetBorrowMetaData()

	pair, _ := k.GetLendPair(ctx, liqBorrow.ExtendedPairId)
	lendPos, _ := k.GetLend(ctx, kind.LendingId)
	borrowPos, _ := k.GetBorrow(ctx, liqBorrow.OriginalVaultId)

	AssetInPool, _ := k.GetPool(ctx, lendPos.PoolID)
	AssetOutPool, _ := k.GetPool(ctx, pair.AssetOutPoolID)
	assetInRatesStats, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)
	borrowPos.IsLiquidated = false

	amoutOutDiff := borrowPos.AmountOut.Amount.Sub(liqBorrow.AmountOut)
	borrowPos.AmountOut.Amount = liqBorrow.AmountOut
	borrowPos.AmountIn.Amount = liqBorrow.AmountIn

	var firstTransitAssetID, secondTransitAssetID uint64
	for _, data := range AssetInPool.AssetData {
		if data.AssetTransitType == 2 {
			firstTransitAssetID = data.AssetID
		}
		if data.AssetTransitType == 3 {
			secondTransitAssetID = data.AssetID
		}
	}

	// Adjusting bridged asset qty after auctions
	if !kind.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
		priceAssetIn, _ := k.market.GetTwa(ctx, pair.AssetIn)
		adjustedBridgedAssetAmt := borrowPos.AmountIn.Amount.ToDec().Mul(assetInRatesStats.Ltv)
		amtIn := adjustedBridgedAssetAmt.TruncateInt().Mul(sdk.NewIntFromUint64(priceAssetIn.Twa))
		priceFirstBridgedAsset, _ := k.market.GetTwa(ctx, firstTransitAssetID)
		priceSecondBridgedAsset, _ := k.market.GetTwa(ctx, secondTransitAssetID)
		firstBridgedAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)

		if kind.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
			firstBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceFirstBridgedAsset.Twa))
			diff := borrowPos.BridgedAssetAmount.Amount.Sub(firstBridgedAssetQty)
			if diff.GT(sdk.ZeroInt()) {
				err := k.bank.SendCoinsFromModuleToModule(ctx, AssetOutPool.ModuleName, AssetInPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, diff)))
				if err != nil {
					return
				}
				borrowPos.BridgedAssetAmount.Amount = firstBridgedAssetQty
			} else {
				newDiff := firstBridgedAssetQty.Sub(borrowPos.BridgedAssetAmount.Amount)
				if newDiff.GT(sdk.ZeroInt()) {
					err := k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, newDiff)))
					if err != nil {
						return
					}
					borrowPos.BridgedAssetAmount.Amount = firstBridgedAssetQty
				}
			}
		} else {
			secondBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset.Twa))
			diff := borrowPos.BridgedAssetAmount.Amount.Sub(secondBridgedAssetQty)
			if diff.GT(sdk.ZeroInt()) {
				err := k.bank.SendCoinsFromModuleToModule(ctx, AssetOutPool.ModuleName, AssetInPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, diff)))
				if err != nil {
					return
				}
				borrowPos.BridgedAssetAmount.Amount = secondBridgedAssetQty
			} else {
				newDiff := secondBridgedAssetQty.Sub(borrowPos.BridgedAssetAmount.Amount)
				if newDiff.GT(sdk.ZeroInt()) {
					err := k.bank.SendCoinsFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, newDiff)))
					if err != nil {
						return
					}
					borrowPos.BridgedAssetAmount.Amount = secondBridgedAssetQty
				}
			}
		}
	}

	k.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, amoutOutDiff, false)
	k.SetBorrow(ctx, borrowPos)
}

func (k Keeper) MsgCalculateBorrowInterest(ctx sdk.Context, borrowerAddr string, borrowID uint64) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	if borrowPos.IsLiquidated {
		return types.ErrorBorrowPosLiquidated
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorized
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	k.SetBorrow(ctx, borrowPos)
	return nil
}

func (k Keeper) MsgCalculateLendRewards(ctx sdk.Context, addr string, lendID uint64) error {
	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()
	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorized
	}
	k.SetLend(ctx, lendPos)
	return nil
}

func (k Keeper) MsgCalculateInterestAndRewards(ctx sdk.Context, addr string) error {
	var (
		lendIDs   []uint64
		borrowIDs []uint64
	)
	fmt.Println(addr)
	mappingData := k.GetUserTotalMappingData(ctx, addr)
	fmt.Println("mappingData", mappingData)

	for _, data := range mappingData {
		lendIDs = append(lendIDs, data.LendId)
	}
	fmt.Println("lendIDs", lendIDs)
	for _, v := range lendIDs {
		lendBorrowMappingData, _ := k.GetUserLendBorrowMapping(ctx, addr, v)
		borrowIDs = append(borrowIDs, lendBorrowMappingData.BorrowId...)
	}
	if len(lendIDs) == 0 {
		return types.ErrLendNotFound
	}
	fmt.Println("borrowIDs", borrowIDs)

	if len(borrowIDs) != 0 {
		for _, borrowID := range borrowIDs {
			err := k.MsgCalculateBorrowInterest(ctx, addr, borrowID)
			if err != nil {
				continue
			}
		}
	}
	fmt.Println("till here")
	for _, lendID := range lendIDs {
		err := k.MsgCalculateLendRewards(ctx, addr, lendID)
		if err != nil {
			return err
		}
	}
	return nil
}
