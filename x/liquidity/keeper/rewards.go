package keeper

import (
	"math"
	"sort"
	"time"

	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetAMMPoolInterfaceObject(ctx sdk.Context, appID, poolID uint64) (*types.Pool, *types.Pair, *amm.BasicPool, error) {
	pool, found := k.GetPool(ctx, appID, poolID)
	if !found {
		return nil, nil, nil, sdkerrors.Wrapf(types.ErrInvalidPoolID, "pool %d is invalid", poolID)
	}
	if pool.Disabled {
		return nil, nil, nil, sdkerrors.Wrapf(types.ErrDisabledPool, "pool %d is disabled", poolID)
	}

	pair, _ := k.GetPair(ctx, pool.AppId, pool.PairId)
	rx, ry := k.getPoolBalances(ctx, pool, pair)
	ps := k.GetPoolCoinSupply(ctx, pool)
	ammPool := amm.NewBasicPool(rx.Amount, ry.Amount, ps)
	if ammPool.IsDepleted() {
		return nil, nil, nil, sdkerrors.Wrapf(types.ErrDepletedPool, "pool %d is depleted", poolID)
	}
	return &pool, &pair, ammPool, nil
}

func (k Keeper) CalculateXYFromPoolCoin(ctx sdk.Context, ammPool *amm.BasicPool, poolCoin sdk.Coin) (sdk.Int, sdk.Int, error) {
	// ammPool.Withdraw implemets the actual logic for pool token ratio calculation
	x, y := ammPool.Withdraw(poolCoin.Amount, sdk.NewDec(0))
	if x.IsZero() && y.IsZero() {
		return sdk.NewInt(0), sdk.NewInt(0), types.ErrCalculatedPoolAmountIsZero
	}
	return x, y, nil
}

func (k Keeper) OraclePrice(ctx sdk.Context, denom string) (uint64, bool) {

	asset, found := k.assetKeeper.GetAssetForDenom(ctx, denom)
	if !found {
		return 0, false
	}

	price, found := k.marketKeeper.GetPriceForAsset(ctx, asset.Id)
	if !found {
		return 0, false
	}
	return price, true
}

func (k Keeper) GetOraclePrices(ctx sdk.Context, quoteCoinDenom, baseCoinDenom string) (sdk.Dec, string, error) {
	oraclePrice, found := k.OraclePrice(ctx, quoteCoinDenom)
	denom := quoteCoinDenom
	if !found {
		oraclePrice, found = k.OraclePrice(ctx, baseCoinDenom)
		denom = baseCoinDenom
		if !found {
			return sdk.NewDec(0), "", types.ErrOraclePricesNotFound
		}
	}
	// considering oracle prices in 10^6
	return sdk.NewDec(int64(oraclePrice)).Quo(sdk.NewDec(1000000)), denom, nil
}

func (k Keeper) CalculateLiquidityAddedValue(
	ctx sdk.Context,
	quoteCoinPoolBalance, baseCoinPoolBalance sdk.Coin,
	quoteCoin, baseCoin sdk.Coin,
	oraclePrice sdk.Dec,
	oraclePriceDenom string,
) sdk.Dec {
	poolSupplyX := sdk.NewDecFromInt(quoteCoinPoolBalance.Amount)
	poolSupplyY := sdk.NewDecFromInt(baseCoinPoolBalance.Amount)

	baseCoinPoolPrice := sdk.NewDec(0)
	quoteCoinPoolPrice := sdk.NewDec(0)
	if oraclePriceDenom == quoteCoin.Denom {
		baseCoinPoolPrice = poolSupplyX.Quo(poolSupplyY).Mul(oraclePrice)
		quoteCoinPoolPrice = poolSupplyY.Quo(poolSupplyX).Mul(baseCoinPoolPrice)
	} else {
		quoteCoinPoolPrice = poolSupplyY.Quo(poolSupplyX).Mul(oraclePrice)
		baseCoinPoolPrice = poolSupplyX.Quo(poolSupplyY).Mul(quoteCoinPoolPrice)
	}

	supplyX := sdk.NewDecFromInt(quoteCoin.Amount)
	supplyY := sdk.NewDecFromInt(baseCoin.Amount)

	// returns actual $value of quoteCoin + baseCoin (value returned in 10^-6, i.e 1000000=1$)
	return supplyX.Mul(quoteCoinPoolPrice).Add(supplyY.Mul(baseCoinPoolPrice)).Quo(sdk.NewDec(1000000))
}

func (k Keeper) GetAggregatedChildPoolContributions(ctx sdk.Context, appID uint64, poolIds []uint64, masterPoolSupplyAddresses []sdk.AccAddress) map[string]sdk.Dec {
	poolSupplyData := make(map[string]sdk.Dec)

	for _, poolID := range poolIds {
		liquidityProvidersDataForPool, found := k.GetPoolLiquidityProvidersData(ctx, appID, poolID)
		if !found {
			continue
		}

		pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, appID, poolID)
		if err != nil {
			continue
		}

		oraclePrice, denom, err := k.GetOraclePrices(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom)
		if err != nil {
			continue
		}

		quoteCoinPoolBalance, baseCoinPoolBalance := k.getPoolBalances(ctx, *pool, *pair)

		for _, address := range masterPoolSupplyAddresses {
			supplyData, found := liquidityProvidersDataForPool.LiquidityProviders[address.String()]
			if !found {
				continue
			}
			depositedCoins := supplyData.Coins
			for _, coin := range depositedCoins {
				if coin.Denom == pool.PoolCoinDenom {
					x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, coin)
					if err != nil {
						continue
					}
					quoteCoin := sdk.NewCoin(pair.QuoteCoinDenom, x)
					baseCoin := sdk.NewCoin(pair.BaseCoinDenom, y)
					value := k.CalculateLiquidityAddedValue(ctx, quoteCoinPoolBalance, baseCoinPoolBalance, quoteCoin, baseCoin, oraclePrice, denom)
					_, found := poolSupplyData[address.String()]
					if !found {
						poolSupplyData[address.String()] = value
					} else {
						poolSupplyData[address.String()] = poolSupplyData[address.String()].Add(value)
					}
				}
			}
		}
	}
	return poolSupplyData
}

func (k Keeper) GetFarmingRewardsData(ctx sdk.Context, appID uint64, coinsToDistribute sdk.Coin, liquidityGaugeData rewardstypes.LiquidtyGaugeMetaData) ([]rewardstypes.RewardDistributionDataCollector, error) {
	liquidityProvidersDataForPool, found := k.GetPoolLiquidityProvidersData(ctx, appID, liquidityGaugeData.PoolId)
	if !found {
		return []rewardstypes.RewardDistributionDataCollector{}, nil
	}

	pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, appID, liquidityGaugeData.PoolId)
	if err != nil {
		return nil, err
	}

	oraclePrice, denom, err := k.GetOraclePrices(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom)
	if err != nil {
		return nil, err
	}

	quoteCoinPoolBalance, baseCoinPoolBalance := k.getPoolBalances(ctx, *pool, *pair)

	lpAddresses := []sdk.AccAddress{}
	lpSupplies := []sdk.Dec{}

	for address, depositedCoins := range liquidityProvidersDataForPool.LiquidityProviders {
		addr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			continue
		}
		for _, coin := range depositedCoins.Coins {
			if coin.Denom == pool.PoolCoinDenom {
				x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, coin)
				if err != nil {
					continue
				}
				quoteCoin := sdk.NewCoin(pair.QuoteCoinDenom, x)
				baseCoin := sdk.NewCoin(pair.BaseCoinDenom, y)
				value := k.CalculateLiquidityAddedValue(ctx, quoteCoinPoolBalance, baseCoinPoolBalance, quoteCoin, baseCoin, oraclePrice, denom)
				lpAddresses = append(lpAddresses, addr)
				lpSupplies = append(lpSupplies, value)
			}
		}
	}

	// Logic for master pool mechanism
	if liquidityGaugeData.IsMasterPool {
		childPoolSupplies := []sdk.Dec{}
		minMasterChildPoolSupplies := []sdk.Dec{}

		childPoolIds := []uint64{}
		if len(liquidityGaugeData.ChildPoolIds) == 0 {
			pools := k.GetAllPools(ctx, appID)
			for _, pool := range pools {
				if pool.Id != liquidityGaugeData.PoolId {
					childPoolIds = append(childPoolIds, pool.Id)
				}
			}
		} else {
			// sanitization
			for _, poolID := range liquidityGaugeData.ChildPoolIds {
				if poolID != liquidityGaugeData.PoolId {
					childPoolIds = append(childPoolIds, poolID)
				}
			}
		}

		chilPoolSuppliesData := k.GetAggregatedChildPoolContributions(ctx, appID, childPoolIds, lpAddresses)

		for _, accAddress := range lpAddresses {
			aggregatedSupplyValue, found := chilPoolSuppliesData[accAddress.String()]
			if !found {
				childPoolSupplies = append(childPoolSupplies, sdk.NewDec(0))
			} else {
				childPoolSupplies = append(childPoolSupplies, aggregatedSupplyValue)
			}
		}

		if len(lpAddresses) != len(lpSupplies) || len(lpAddresses) != len(childPoolSupplies) {
			return nil, types.ErrSupplyValueCalculationInvalid
		}

		totalRewardEligibleSupply := sdk.NewDec(0)
		for i := 0; i < len(lpAddresses); i++ {
			var minSupply sdk.Dec
			if lpSupplies[i].LTE(childPoolSupplies[i]) {
				minSupply = lpSupplies[i]
			} else {
				minSupply = childPoolSupplies[i]
			}
			totalRewardEligibleSupply = totalRewardEligibleSupply.Add(minSupply)
			minMasterChildPoolSupplies = append(minMasterChildPoolSupplies, minSupply)
		}

		rewardData := []rewardstypes.RewardDistributionDataCollector{}
		if !totalRewardEligibleSupply.IsZero() {
			multiplier := sdk.NewDecFromInt(coinsToDistribute.Amount).Quo(totalRewardEligibleSupply)
			for index, address := range lpAddresses {
				if !minMasterChildPoolSupplies[index].IsZero() {
					calculatedReward := int64(math.Floor(minMasterChildPoolSupplies[index].Mul(multiplier).MustFloat64()))
					newData := new(rewardstypes.RewardDistributionDataCollector)
					newData.RewardReceiver = address
					newData.RewardCoin = sdk.NewCoin(coinsToDistribute.Denom, sdk.NewInt(calculatedReward))
					rewardData = append(rewardData, *newData)
				}
			}
		}

		return rewardData, nil
	}
	// Logic for non master pool gauges (external rewards)
	totalRewardEligibleSupply := sdk.NewDec(0)
	for _, supply := range lpSupplies {
		totalRewardEligibleSupply = totalRewardEligibleSupply.Add(supply)
	}

	rewardData := []rewardstypes.RewardDistributionDataCollector{}
	if !totalRewardEligibleSupply.IsZero() {
		multiplier := sdk.NewDecFromInt(coinsToDistribute.Amount).Quo(totalRewardEligibleSupply)
		for index, address := range lpAddresses {
			calculatedReward := int64(math.Floor(lpSupplies[index].Mul(multiplier).MustFloat64()))
			newData := new(rewardstypes.RewardDistributionDataCollector)
			newData.RewardReceiver = address
			newData.RewardCoin = sdk.NewCoin(coinsToDistribute.Denom, sdk.NewInt(calculatedReward))
			rewardData = append(rewardData, *newData)
		}
	}

	return rewardData, nil
}

func (k Keeper) ValidateMsgTokensSoftLock(ctx sdk.Context, msg *types.MsgTokensSoftLock) (sdk.AccAddress, error) {
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	pool, found := k.GetPool(ctx, msg.AppId, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", msg.PoolId)
	}

	if msg.SoftLockCoin.Denom != pool.PoolCoinDenom {
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", pool.PoolCoinDenom, msg.SoftLockCoin.Denom)
	}
	return depositor, nil
}

func (k Keeper) SoftLockTokens(ctx sdk.Context, msg *types.MsgTokensSoftLock) error {
	depositor, err := k.ValidateMsgTokensSoftLock(ctx, msg)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(msg.SoftLockCoin))
	if err != nil {
		return err
	}

	liquidityProvidersData, found := k.GetPoolLiquidityProvidersData(ctx, msg.AppId, msg.PoolId)

	if !found {
		liquidityProvidersData = types.PoolLiquidityProvidersData{
			AppId:                    msg.AppId,
			PoolId:                   msg.PoolId,
			BondedLockIds:            []uint64{},
			LiquidityProviders:       make(map[string]*types.DepositsMade),
			QueuedLiquidityProviders: []*types.QueuedLiquidityProvider{},
		}
	}

	liquidityProvidersData.QueuedLiquidityProviders = append(
		liquidityProvidersData.QueuedLiquidityProviders,
		&types.QueuedLiquidityProvider{
			Address:        depositor.String(),
			SupplyProvided: []*sdk.Coin{&msg.SoftLockCoin},
			CreatedAt:      ctx.BlockTime(),
		})
	k.SetPoolLiquidityProvidersData(ctx, liquidityProvidersData)

	return nil
}

func (k Keeper) ValidateMsgTokensSoftUnlock(ctx sdk.Context, msg *types.MsgTokensSoftUnlock) (sdk.AccAddress, error) {
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	pool, found := k.GetPool(ctx, msg.AppId, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", msg.PoolId)
	}

	if msg.SoftUnlockCoin.Denom != pool.PoolCoinDenom {
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", pool.PoolCoinDenom, msg.SoftUnlockCoin.Denom)
	}
	return depositor, nil
}

func (k Keeper) SoftUnlockTokens(ctx sdk.Context, msg *types.MsgTokensSoftUnlock) error {
	depositor, err := k.ValidateMsgTokensSoftUnlock(ctx, msg)
	if err != nil {
		return err
	}

	liquidityProvidersData, found := k.GetPoolLiquidityProvidersData(ctx, msg.AppId, msg.PoolId)

	if !found {
		return sdkerrors.Wrapf(types.ErrNoSoftLockPresent, "no soft locks present for given pool id %d", msg.PoolId)
	}

	sortedQueuedLiquidityProviders := liquidityProvidersData.QueuedLiquidityProviders
	sort.Slice(sortedQueuedLiquidityProviders, func(i, j int) bool {
		return sortedQueuedLiquidityProviders[i].CreatedAt.After(sortedQueuedLiquidityProviders[j].CreatedAt)
	})

	updatedQueuedLiquidityProviders := []*types.QueuedLiquidityProvider{}

	refundAmount := sdk.NewCoin(msg.SoftUnlockCoin.Denom, sdk.NewInt(0))

	for queryIndex, queuedLiquidityProvider := range sortedQueuedLiquidityProviders {
		if queuedLiquidityProvider.Address == depositor.String() {
			for qindex, qcoin := range queuedLiquidityProvider.SupplyProvided {
				if msg.SoftUnlockCoin.Denom == qcoin.Denom && !msg.SoftUnlockCoin.Amount.IsZero() {
					if msg.SoftUnlockCoin.Amount.GTE(qcoin.Amount) {
						msg.SoftUnlockCoin.Amount = msg.SoftUnlockCoin.Amount.Sub(qcoin.Amount)
						refundAmount.Amount = refundAmount.Amount.Add(qcoin.Amount)
						sortedQueuedLiquidityProviders[queryIndex].SupplyProvided[qindex].Amount = sdk.NewInt(0)
					} else {
						sortedQueuedLiquidityProviders[queryIndex].SupplyProvided[qindex].Amount = sortedQueuedLiquidityProviders[queryIndex].SupplyProvided[qindex].Amount.Sub(msg.SoftUnlockCoin.Amount)
						refundAmount.Amount = refundAmount.Amount.Add(msg.SoftUnlockCoin.Amount)
						msg.SoftUnlockCoin.Amount = sdk.NewInt(0)
					}
				}
			}

			canRemoveThisQueuedLiquidityProvider := true
			for _, qcoin := range sortedQueuedLiquidityProviders[queryIndex].SupplyProvided {
				if !qcoin.Amount.IsZero() {
					canRemoveThisQueuedLiquidityProvider = false
				}
			}
			if !canRemoveThisQueuedLiquidityProvider {
				updatedQueuedLiquidityProviders = append(updatedQueuedLiquidityProviders, sortedQueuedLiquidityProviders[queryIndex])
			}

			withdrawCoinFullfilled := true
			if !msg.SoftUnlockCoin.Amount.IsZero() {
				withdrawCoinFullfilled = false
			}

			if withdrawCoinFullfilled {
				updatedQueuedLiquidityProviders = append(updatedQueuedLiquidityProviders, sortedQueuedLiquidityProviders[queryIndex+1:]...)
				break
			}
		} else {
			updatedQueuedLiquidityProviders = append(updatedQueuedLiquidityProviders, sortedQueuedLiquidityProviders[queryIndex])
		}
	}

	liquidityProvidersData.QueuedLiquidityProviders = updatedQueuedLiquidityProviders

	coinWithdrawingFullfilled := true
	if !msg.SoftUnlockCoin.Amount.IsZero() {
		coinWithdrawingFullfilled = false
	}

	if !coinWithdrawingFullfilled {
		if liquidityProvidersData.LiquidityProviders[depositor.String()] != nil {
			providedSupply := liquidityProvidersData.LiquidityProviders[depositor.String()].Coins
			for psIndex, psCoin := range providedSupply {
				if psCoin.Denom == msg.SoftUnlockCoin.Denom {
					if msg.SoftUnlockCoin.Amount.GTE(psCoin.Amount) {
						liquidityProvidersData.LiquidityProviders[depositor.String()].Coins[psIndex].Amount = sdk.NewInt(0)
						refundAmount.Amount = refundAmount.Amount.Add(psCoin.Amount)
						msg.SoftUnlockCoin.Amount = msg.SoftUnlockCoin.Amount.Sub(psCoin.Amount)
					} else {
						liquidityProvidersData.LiquidityProviders[depositor.String()].Coins[psIndex].Amount = liquidityProvidersData.LiquidityProviders[depositor.String()].Coins[psIndex].Amount.Sub(msg.SoftUnlockCoin.Amount)
						refundAmount.Amount = refundAmount.Amount.Add(msg.SoftUnlockCoin.Amount)
						msg.SoftUnlockCoin.Amount = sdk.NewInt(0)
					}
				}
			}
		}
	}

	if msg.SoftUnlockCoin.Amount.Add(refundAmount.Amount).GT(refundAmount.Amount) {
		return sdkerrors.Wrapf(types.ErrInvalidUnlockAmount, "available soft locked amount %d%s smaller than requested amount %d%s", refundAmount.Amount.Int64(), refundAmount.Denom, msg.SoftUnlockCoin.Amount.Add(refundAmount.Amount).Int64(), refundAmount.Denom)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(refundAmount))
	if err != nil {
		return err
	}

	k.SetPoolLiquidityProvidersData(ctx, liquidityProvidersData)

	return nil
}

func (k Keeper) GetMinimumEpochDurationFromPoolID(ctx sdk.Context, poolID uint64, gauges []rewardstypes.Gauge) time.Duration {
	var minEpochDuration time.Duration
	for _, gauge := range gauges {
		switch kind := gauge.Kind.(type) {
		case *rewardstypes.Gauge_LiquidityMetaData:
			if kind.LiquidityMetaData.PoolId == poolID && gauge.IsActive {
				if minEpochDuration == time.Duration(0) {
					minEpochDuration = gauge.TriggerDuration
				} else if gauge.TriggerDuration < minEpochDuration {
					minEpochDuration = gauge.TriggerDuration
				}
			}
		}
	}
	if minEpochDuration == time.Duration(0) {
		minEpochDuration = time.Hour * 24 // if no gauge for given pool, making 24h as default vaule
	}
	return minEpochDuration
}

func (k Keeper) ProcessQueuedLiquidityProviders(ctx sdk.Context, appID uint64) {
	availablePools := k.GetAllPools(ctx, appID)
	availableLiquidityGauges := k.rewardsKeeper.GetAllGaugesByGaugeTypeID(ctx, rewardstypes.LiquidityGaugeTypeID)

	for _, pool := range availablePools {
		poolLpData, found := k.GetPoolLiquidityProvidersData(ctx, pool.AppId, pool.Id)
		minEpochDuration := k.GetMinimumEpochDurationFromPoolID(ctx, pool.Id, availableLiquidityGauges)

		if !found {
			continue
		}
		queuedDepositRequests := poolLpData.QueuedLiquidityProviders
		updatedQueue := []*types.QueuedLiquidityProvider{}
		for _, queuedLp := range queuedDepositRequests {
			if ctx.BlockTime().Before(queuedLp.CreatedAt.Add(minEpochDuration)) {
				updatedQueue = append(updatedQueue, queuedLp)
			} else {
				suppliedCoins := []sdk.Coin{}
				for _, coin := range queuedLp.SupplyProvided {
					suppliedCoins = append(suppliedCoins, sdk.NewCoin(coin.Denom, coin.Amount))
				}
				if poolLpData.LiquidityProviders[queuedLp.Address] == nil {
					if len(poolLpData.LiquidityProviders) == 0 {
						poolLpData.LiquidityProviders = make(map[string]*types.DepositsMade)
					}
					newDeposit := new(types.DepositsMade)
					newDeposit.Coins = suppliedCoins
					poolLpData.LiquidityProviders[queuedLp.Address] = newDeposit
				} else {
					for _, coin := range suppliedCoins {
						for lpedIndex, lpedCoin := range poolLpData.LiquidityProviders[queuedLp.Address].Coins {
							if coin.Denom == lpedCoin.Denom {
								poolLpData.LiquidityProviders[queuedLp.Address].Coins[lpedIndex].Amount = poolLpData.LiquidityProviders[queuedLp.Address].Coins[lpedIndex].Amount.Add(coin.Amount)
								break
							}
						}
					}
				}
			}
		}
		poolLpData.QueuedLiquidityProviders = updatedQueue
		k.SetPoolLiquidityProvidersData(ctx, poolLpData)
	}
}
