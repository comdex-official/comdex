package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/expected"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		bank       expected.BankKeeper
		account    expected.AccountKeeper
		asset      expected.AssetKeeper
		market     expected.MarketKeeper
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
		bank:       bank,
		account:    account,
		asset:      asset,
		market:     market,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func uint64InAssetData(a uint64, list []types.AssetDataPoolMapping) bool {
	for _, b := range list {
		if b.AssetId == a {
			return true
		}
	}
	return false
}

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr string, AssetID uint64, Amount sdk.Coin, PoolID uint64) error {
	asset, _ := k.GetAsset(ctx, AssetID)
	pool, _ := k.GetPool(ctx, PoolID)

	if Amount.Denom != asset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, Amount.Denom)
	}

	found := uint64InAssetData(AssetID, pool.AssetData)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidAssetIDForPool, strconv.FormatUint(AssetID, 10))
	}

	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if k.HasLendForAddressByAsset(ctx, addr, AssetID, PoolID) {
		return types.ErrorDuplicateLend
	}

	loanTokens := sdk.NewCoins(Amount)

	assetRatesStat, found := k.GetAssetRatesStats(ctx, AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(AssetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
	cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, loanTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err := k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////
	lendID := k.GetUserLendIDHistory(ctx)

	lendPos := types.LendAsset{
		ID:                 lendID + 1,
		AssetId:            AssetID,
		PoolId:             PoolID,
		Owner:              lenderAddr,
		AmountIn:           Amount,
		LendingTime:        ctx.BlockTime(),
		UpdatedAmountIn:    Amount.Amount,
		AvailableToBorrow:  Amount.Amount,
		Reward_Accumulated: sdk.ZeroInt(),
	}
	assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, AssetID, PoolID)
	if !found {
		assetStats.TotalLend = sdk.ZeroInt()
	}
	AssetStats := types.AssetStats{
		PoolId:    PoolID,
		AssetId:   AssetID,
		TotalLend: assetStats.TotalLend.Add(Amount.Amount),
	}
	k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
	k.SetUserLendIDHistory(ctx, lendPos.ID)
	k.SetLend(ctx, lendPos)
	k.SetLendForAddressByAsset(ctx, addr, lendPos.AssetId, lendPos.ID, lendPos.PoolId)
	err = k.UpdateUserLendIDMapping(ctx, lenderAddr, lendPos.ID, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDByOwnerAndPoolMapping(ctx, lenderAddr, lendPos.ID, lendPos.PoolId, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDsMapping(ctx, lendPos.ID, true)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, addr string, lendID uint64, withdrawal sdk.Coin) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}
	getAsset, _ := k.GetAsset(ctx, lendPos.AssetId)
	pool, _ := k.GetPool(ctx, lendPos.PoolId)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	if withdrawal.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, withdrawal.Denom)
	}

	reservedAmount := k.GetReserveFunds(ctx, withdrawal.Denom)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, withdrawal.Denom)

	if withdrawal.Amount.GT(lendPos.AmountIn.Amount) {
		return sdkerrors.Wrap(types.ErrWithdrawalAmountExceeds, withdrawal.String())
	}

	if withdrawal.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, withdrawal.String())
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
	cToken := sdk.NewCoin(cAsset.Denom, withdrawal.Amount)

	tokens := sdk.NewCoins(withdrawal)
	if err != nil {
		return err
	}
	_, found = k.GetLendIDToBorrowIDMapping(ctx, lendID)
	if !found {
		if withdrawal.Amount.LT(lendPos.UpdatedAmountIn) {
			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
				return err
			}

			//burn c/Token
			cTokens := sdk.NewCoins(cToken)
			err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
			if err != nil {
				return err
			}

			if err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
				return err
			}

			lendPos.AmountIn = lendPos.AmountIn.Sub(withdrawal)
			lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(withdrawal.Amount)
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
			assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.AssetId, lendPos.PoolId)
			assetStats.TotalLend = assetStats.TotalLend.Sub(withdrawal.Amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
			k.SetLend(ctx, lendPos)
		} else {
			return nil
		}
	} else {
		if withdrawal.Amount.LT(lendPos.UpdatedAmountIn) {
			// add CR validation
			// lend to borrow mapping

			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
				return err
			}

			//burn c/Token
			cTokens := sdk.NewCoins(cToken)
			err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
			if err != nil {
				return err
			}

			if err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
				return err
			}

			lendPos.AmountIn = lendPos.AmountIn.Sub(withdrawal)
			lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(withdrawal.Amount)
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
			assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.AssetId, lendPos.PoolId)
			assetStats.TotalLend = assetStats.TotalLend.Sub(withdrawal.Amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
			k.SetLend(ctx, lendPos)
		} else {
			return nil
		}
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

	getAsset, _ := k.GetAsset(ctx, lendPos.AssetId)
	pool, _ := k.GetPool(ctx, lendPos.PoolId)

	if deposit.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, deposit.Denom)
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
	cToken := sdk.NewCoin(cAsset.Denom, deposit.Amount)

	cTokens := sdk.NewCoins(cToken)

	if err = k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(deposit)); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, cTokens)
	if err != nil {
		return err
	}

	lendPos.AmountIn = lendPos.AmountIn.Add(deposit)
	lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Add(deposit.Amount)
	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(deposit.Amount)
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.AssetId, lendPos.PoolId)
	assetStats.TotalLend = assetStats.TotalLend.Add(deposit.Amount)
	k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
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
	pool, _ := k.GetPool(ctx, lendPos.PoolId)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	_, found = k.GetLendIDToBorrowIDMapping(ctx, lendID)
	if found {
		return types.ErrBorrowingPositionOpen
	}
	reservedAmount := k.GetReserveFunds(ctx, lendPos.AmountIn.Denom)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, lendPos.AmountIn.Denom)

	if lendPos.UpdatedAmountIn.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, lendPos.UpdatedAmountIn.String())
	}

	tokens := sdk.NewCoins(sdk.NewCoin(lendPos.AmountIn.Denom, lendPos.UpdatedAmountIn))
	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
	cToken := sdk.NewCoin(cAsset.Denom, lendPos.UpdatedAmountIn)

	if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
		return err
	}

	cTokens := sdk.NewCoins(cToken)
	err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
	if err != nil {
		return err
	}

	if err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
		return err
	}

	k.DeleteLendForAddressByAsset(ctx, lenderAddr, lendPos.AssetId, lendPos.PoolId)

	err = k.UpdateUserLendIDMapping(ctx, addr, lendPos.ID, false)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDByOwnerAndPoolMapping(ctx, addr, lendPos.ID, lendPos.PoolId, false)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDsMapping(ctx, lendPos.ID, false)
	if err != nil {
		return err
	}
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.AssetId, lendPos.PoolId)
	assetStats.TotalLend = assetStats.TotalLend.Sub(lendPos.UpdatedAmountIn)
	k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
	k.DeleteLend(ctx, lendPos.ID)
	return nil
}

//nolint:funlen
func (k Keeper) BorrowAsset(ctx sdk.Context, addr string, lendID, pairID uint64, IsStableBorrow bool, AmountIn, loan sdk.Coin) error {
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}
	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	//check cr ratio
	assetIn, _ := k.GetAsset(ctx, lendPos.AssetId)
	assetOut, _ := k.GetAsset(ctx, pair.AssetOut)
	assetInRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)

	cAsset, _ := k.GetAsset(ctx, assetInRatesStats.CAssetId)
	if AmountIn.Denom != cAsset.Denom {
		return types.ErrBadOfferCoinAmount
	}

	if k.HasBorrowForAddressByPair(ctx, lenderAddr, pairID) {
		return types.ErrorDuplicateBorrow
	}

	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}

	AssetInPool, _ := k.GetPool(ctx, lendPos.PoolId)
	AssetOutPool, _ := k.GetPool(ctx, pair.AssetOutPoolId)

	assetRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetOut)
	if IsStableBorrow && !assetRatesStats.EnableStableBorrow {
		return sdkerrors.Wrap(types.ErrStableBorrowDisabled, loan.String())
	}

	err := k.VerifyCollaterlizationRatio(ctx, AmountIn.Amount, assetIn, loan.Amount, assetOut, assetRatesStats.LiquidationThreshold)
	if err != nil {
		return err
	}
	borrowID := k.GetUserBorrowIDHistory(ctx)

	if !pair.IsInterPool {
		// check sufficient amt in pool to borrow
		reservedAmount := k.GetReserveFunds(ctx, loan.Denom)
		availableAmount := k.ModuleBalance(ctx, AssetOutPool.ModuleName, loan.Denom)

		if loan.Amount.GT(availableAmount.Sub(reservedAmount)) {
			return sdkerrors.Wrap(types.ErrBorrowingPoolInsufficient, loan.String())
		}

		AmountOut := loan
		// take c/Tokens from the user
		if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
			return err
		}

		if err := k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
			return err
		}

		var StableBorrowRate sdk.Dec
		if assetRatesStats.EnableStableBorrow {
			if IsStableBorrow {
				StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolId, pair.AssetOut, IsStableBorrow)
				if err != nil {
					return err
				}
			} else {
				StableBorrowRate = sdk.ZeroDec()
			}
		} else {
			IsStableBorrow = false
			StableBorrowRate = sdk.ZeroDec()
		}

		borrowPos := types.BorrowAsset{
			ID:                   borrowID + 1,
			LendingID:            lendID,
			PairID:               pairID,
			AmountIn:             AmountIn,
			AmountOut:            AmountOut,
			BridgedAssetAmount:   sdk.NewCoin("", sdk.NewInt(0)),
			IsStableBorrow:       IsStableBorrow,
			StableBorrowRate:     StableBorrowRate,
			BorrowingTime:        ctx.BlockTime(),
			UpdatedAmountOut:     AmountOut.Amount,
			Interest_Accumulated: sdk.ZeroInt(),
		}

		assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolId)
		if !found {
			assetStats.TotalBorrowed = sdk.ZeroInt()
		}
		if borrowPos.IsStableBorrow {
			AssetStats := types.AssetStats{
				PoolId:              pair.AssetOutPoolId,
				AssetId:             pair.AssetOut,
				TotalStableBorrowed: assetStats.TotalStableBorrowed.Add(AmountOut.Amount),
			}
			k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
		} else {
			AssetStats := types.AssetStats{
				PoolId:        pair.AssetOutPoolId,
				AssetId:       pair.AssetOut,
				TotalBorrowed: assetStats.TotalBorrowed.Add(AmountOut.Amount),
			}
			k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
		}
		if borrowPos.IsStableBorrow {
			err := k.UpdateStableBorrowIdsMapping(ctx, borrowPos.ID, true)
			if err != nil {
				return err
			}
		}
		err := k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
		if err != nil {
			return err
		}
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)

		k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
		k.SetBorrow(ctx, borrowPos)
		k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
		err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
		if err != nil {
			return err
		}
		err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolId, true)
		if err != nil {
			return err
		}
		err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
		if err != nil {
			return err
		}
	} else {
		reservedAmount := k.GetReserveFunds(ctx, loan.Denom)
		availableAmount := k.ModuleBalance(ctx, AssetOutPool.ModuleName, loan.Denom)

		if loan.Amount.GT(availableAmount.Sub(reservedAmount)) {
			return sdkerrors.Wrap(types.ErrBorrowingPoolInsufficient, loan.String())
		}
		assetIn := lendPos.UpdatedAmountIn
		priceAssetIn, found := k.GetPriceForAsset(ctx, pair.AssetIn)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		amtIn := assetIn.Mul(sdk.NewIntFromUint64(priceAssetIn))

		//priceFirstBridgedAsset, _ := k.GetPriceForAsset(ctx, AssetInPool.FirstBridgedAssetId)
		priceSecondBridgedAsset, found := k.GetPriceForAsset(ctx, AssetInPool.SecondBridgedAssetId)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		firstBridgedAsset, _ := k.GetAsset(ctx, AssetInPool.FirstBridgedAssetId)
		secondBridgedAsset, _ := k.GetAsset(ctx, AssetInPool.SecondBridgedAssetId)

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		firstBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(1000000))
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstBridgedAsset.Denom)
		secondBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset))
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondBridgedAsset.Denom)

		if firstBridgedAssetQty.LT(firstBridgedAssetBal) { // take c/Tokens from the user
			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}
			bridgedAssetAmount := sdk.NewCoin(firstBridgedAsset.Denom, firstBridgedAssetQty)
			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
				return err
			}

			AmountOut := loan

			var StableBorrowRate sdk.Dec
			if assetRatesStats.EnableStableBorrow {
				if IsStableBorrow {
					StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolId, pair.AssetOut, IsStableBorrow)
					if err != nil {
						return err
					}
				} else {
					StableBorrowRate = sdk.ZeroDec()
				}
			} else {
				IsStableBorrow = false
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                   borrowID + 1,
				LendingID:            lendID,
				PairID:               pairID,
				AmountIn:             AmountIn,
				AmountOut:            AmountOut,
				BridgedAssetAmount:   bridgedAssetAmount,
				IsStableBorrow:       IsStableBorrow,
				StableBorrowRate:     StableBorrowRate,
				BorrowingTime:        ctx.BlockTime(),
				UpdatedAmountOut:     AmountOut.Amount,
				Interest_Accumulated: sdk.ZeroInt(),
			}

			assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolId)
			if !found {
				assetStats.TotalBorrowed = sdk.ZeroInt()
			}
			if borrowPos.StableBorrowRate != sdk.ZeroDec() {
				AssetStats := types.AssetStats{
					PoolId:              pair.AssetOutPoolId,
					AssetId:             pair.AssetOut,
					TotalStableBorrowed: assetStats.TotalStableBorrowed.Add(AmountOut.Amount),
				}
				k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
			} else {
				AssetStats := types.AssetStats{
					PoolId:        pair.AssetOutPoolId,
					AssetId:       pair.AssetOut,
					TotalBorrowed: assetStats.TotalBorrowed.Add(AmountOut.Amount),
				}
				k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
			}
			if borrowPos.IsStableBorrow {
				err := k.UpdateStableBorrowIdsMapping(ctx, borrowPos.ID, true)
				if err != nil {
					return err
				}
			}
			err := k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
			if err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)
			k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
			err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
			if err != nil {
				return err
			}
			err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolId, true)
			if err != nil {
				return err
			}
			err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
			if err != nil {
				return err
			}
		} else if secondBridgedAssetQty.LT(secondBridgedAssetBal) {
			// take c/Tokens from the user
			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			bridgedAssetAmount := sdk.NewCoin(secondBridgedAsset.Denom, secondBridgedAssetQty)
			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
				return err
			}

			AmountOut := loan

			var StableBorrowRate sdk.Dec
			if assetRatesStats.EnableStableBorrow {
				if IsStableBorrow {
					StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolId, pair.AssetOut, IsStableBorrow)
					if err != nil {
						return err
					}
				} else {
					StableBorrowRate = sdk.ZeroDec()
				}
			} else {
				IsStableBorrow = false
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                   borrowID + 1,
				LendingID:            lendID,
				PairID:               pairID,
				AmountIn:             AmountIn,
				AmountOut:            AmountOut,
				BridgedAssetAmount:   bridgedAssetAmount,
				IsStableBorrow:       IsStableBorrow,
				StableBorrowRate:     StableBorrowRate,
				BorrowingTime:        ctx.BlockTime(),
				UpdatedAmountOut:     AmountOut.Amount,
				Interest_Accumulated: sdk.ZeroInt(),
			}

			assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolId)
			if !found {
				assetStats.TotalBorrowed = sdk.ZeroInt()
			}
			if borrowPos.IsStableBorrow {
				AssetStats := types.AssetStats{
					PoolId:              pair.AssetOutPoolId,
					AssetId:             pair.AssetOut,
					TotalStableBorrowed: assetStats.TotalStableBorrowed.Add(AmountOut.Amount),
				}
				k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
			} else {
				AssetStats := types.AssetStats{
					PoolId:        pair.AssetOutPoolId,
					AssetId:       pair.AssetOut,
					TotalBorrowed: assetStats.TotalBorrowed.Add(AmountOut.Amount),
				}
				k.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
			}

			if borrowPos.IsStableBorrow {
				err := k.UpdateStableBorrowIdsMapping(ctx, borrowPos.ID, true)
				if err != nil {
					return err
				}
			}
			err := k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
			if err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)
			k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
			err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
			if err != nil {
				return err
			}
			err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolId, true)
			if err != nil {
				return err
			}
			err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
			if err != nil {
				return err
			}
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
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, _ := k.GetLendPair(ctx, borrowPos.PairID)
	pool, _ := k.GetPool(ctx, pair.AssetOutPoolId)
	lendPos, _ := k.GetLend(ctx, borrowPos.LendingID)
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	if borrowPos.AmountOut.Denom != payment.Denom {
		return types.ErrBadOfferCoinAmount
	}

	if payment.Amount.LT(borrowPos.UpdatedAmountOut) {
		if payment.Amount.LTE(borrowPos.Interest_Accumulated) {
			// sending repayment to moduleAcc from borrower
			if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
				return err
			}
			borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Sub(payment.Amount)
			borrowPos.Interest_Accumulated = borrowPos.Interest_Accumulated.Sub(payment.Amount)
			k.SetBorrow(ctx, borrowPos)
		} else {
			// sending repayment to moduleAcc from borrower
			if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
				return err
			}
			borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Sub(payment.Amount)
			borrowPos.AmountOut.Amount = borrowPos.AmountOut.Amount.Sub(payment.Amount).Add(borrowPos.Interest_Accumulated)
			borrowPos.Interest_Accumulated = sdk.ZeroInt()

			reserveRates, _ := k.GetReserveRate(ctx, pair.AssetOutPoolId, pair.AssetOut)
			amtToReservePool := sdk.NewDec(int64(borrowPos.AmountOut.Amount.Uint64())).Mul(reserveRates)
			amount := sdk.NewCoins(sdk.NewCoin(payment.Denom, sdk.NewInt(amtToReservePool.TruncateInt64())))
			if err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, types.ModuleName, amount); err != nil {
				return err
			}

			k.SetBorrow(ctx, borrowPos)
		}
	} else {
		return types.ErrInvalidRepayment
	}
	return nil
}

func (k Keeper) DepositBorrowAsset(ctx sdk.Context, borrowID uint64, addr string, AmountIn sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	lendID := borrowPos.LendingID
	pairID := borrowPos.PairID
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}
	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}
	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}

	if k.HasBorrowForAddressByPair(ctx, lenderAddr, pairID) {
		return types.ErrorDuplicateBorrow
	}

	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	AssetInPool, _ := k.GetPool(ctx, lendPos.PoolId)
	AssetOutPool, _ := k.GetPool(ctx, pair.AssetOutPoolId)

	if !pair.IsInterPool {
		AmountIn := sdk.NewCoin(lendPos.AmountIn.Denom, AmountIn.Amount)
		// take c/Tokens from the user
		assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
		if !found {
			return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
		}
		cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
		if AmountIn.Denom != cAsset.Denom {
			return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
		}

		if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
			return err
		}
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)
		borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
		k.SetBorrow(ctx, borrowPos)
	} else {
		assetIn := lendPos.UpdatedAmountIn
		priceAssetIn, found := k.GetPriceForAsset(ctx, pair.AssetIn)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		amtIn := assetIn.Mul(sdk.NewIntFromUint64(priceAssetIn))

		priceFirstBridgedAsset, found := k.GetPriceForAsset(ctx, AssetInPool.FirstBridgedAssetId)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		priceSecondBridgedAsset, found := k.GetPriceForAsset(ctx, AssetInPool.SecondBridgedAssetId)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		firstBridgedAsset, _ := k.GetAsset(ctx, AssetInPool.FirstBridgedAssetId)
		secondBridgedAsset, _ := k.GetAsset(ctx, AssetInPool.SecondBridgedAssetId)

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		firstBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceFirstBridgedAsset))
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstBridgedAsset.Denom)
		secondBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset))
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondBridgedAsset.Denom)

		if firstBridgedAssetQty.LT(firstBridgedAssetBal) { // take c/Tokens from the user
			assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
			if !found {
				return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
			}
			cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
			if AmountIn.Denom != cAsset.Denom {
				return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
			}

			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(firstBridgedAsset.Denom, firstBridgedAssetQty))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			k.SetBorrow(ctx, borrowPos)
		} else if secondBridgedAssetQty.LT(secondBridgedAssetBal) {
			// take c/Tokens from the user
			assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetId)
			if !found {
				return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetId, 10))
			}
			cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
			if AmountIn.Denom != cAsset.Denom {
				return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
			}

			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(secondBridgedAsset.Denom, secondBridgedAssetQty))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			k.SetBorrow(ctx, borrowPos)
		} else {
			return types.ErrBorrowingPoolInsufficient
		}
	}
	return nil
}

func (k Keeper) DrawAsset(ctx sdk.Context, borrowID uint64, borrowerAddr string, payment sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, _ := k.GetLendPair(ctx, borrowPos.PairID)
	pool, _ := k.GetPool(ctx, pair.AssetOutPoolId)
	lendPos, _ := k.GetLend(ctx, borrowPos.LendingID)
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	if borrowPos.AmountOut.Denom != payment.Denom {
		return types.ErrBadOfferCoinAmount
	}
	assetIn, _ := k.GetAsset(ctx, lendPos.AssetId)
	assetOut, _ := k.GetAsset(ctx, pair.AssetOut)

	assetRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetOut)
	err := k.VerifyCollaterlizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut, assetRatesStats.LiquidationThreshold)
	if err != nil {
		return err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, pool.ModuleName, addr, payment); err != nil {
		return err
	}
	borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Add(payment.Amount)
	borrowPos.AmountOut.Add(payment)
	k.DeleteLendIDToBorrowIDMapping(ctx, borrowPos.LendingID)
	k.SetBorrow(ctx, borrowPos)

	return nil
}

func (k Keeper) CloseBorrow(ctx sdk.Context, borrowerAddr string, borrowID uint64) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, _ := k.GetLendPair(ctx, borrowPos.PairID)
	pool, _ := k.GetPool(ctx, pair.AssetOutPoolId)
	lendPos, _ := k.GetLend(ctx, borrowPos.LendingID)
	assetInPool, _ := k.GetPool(ctx, lendPos.PoolId)

	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	assetOut, _ := k.GetAsset(ctx, pair.AssetOut)

	amt := sdk.NewCoins(sdk.NewCoin(assetOut.Denom, borrowPos.UpdatedAmountOut))

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, amt); err != nil {
		return err
	}

	reserveRates, _ := k.GetReserveRate(ctx, pair.AssetOutPoolId, pair.AssetOut)
	amtToReservePool := sdk.NewDec(int64(borrowPos.AmountOut.Amount.Uint64())).Mul(reserveRates)
	amount := sdk.NewCoins(sdk.NewCoin(assetOut.Denom, sdk.NewInt(amtToReservePool.TruncateInt64())))
	if err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, types.ModuleName, amount); err != nil {
		return err
	}

	err := k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, false)
	if err != nil {
		return err
	}
	err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolId, false)
	if err != nil {
		return err
	}
	if borrowPos.IsStableBorrow {
		err := k.UpdateStableBorrowIdsMapping(ctx, borrowPos.ID, false)
		if err != nil {
			return err
		}
	}
	err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, false)
	if err != nil {
		return err
	}
	k.DeleteBorrowForAddressByPair(ctx, addr, borrowPos.PairID)
	err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, false)
	if err != nil {
		return err
	}

	if pair.IsInterPool {
		if err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, assetInPool.ModuleName, sdk.NewCoins(borrowPos.BridgedAssetAmount)); err != nil {
			return err
		}
	}
	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(borrowPos.AmountIn.Amount)
	k.SetLend(ctx, lendPos)
	k.DeleteBorrow(ctx, borrowID)

	return nil
}

func (k Keeper) FundModAcc(ctx sdk.Context, moduleName string, assetID uint64, lenderAddr sdk.AccAddress, payment sdk.Coin) error {
	loanTokens := sdk.NewCoins(payment)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, moduleName, loanTokens); err != nil {
		return err
	}

	_, found := k.GetAsset(ctx, assetID)
	if !found {
		return types.ErrLendNotFound
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, assetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(assetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetId)
	cToken := sdk.NewCoin(cAsset.Denom, payment.Amount)

	err := k.MintCoin(ctx, moduleName, cToken)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}
