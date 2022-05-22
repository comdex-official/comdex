package keeper

import (
	"fmt"
	"math"
	"sort"
	"time"

	incentivestypes "github.com/comdex-official/comdex/x/incentives/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetAMMPoolInterfaceObject(ctx sdk.Context, poolId uint64) (*types.Pool, *types.Pair, *amm.BasicPool, error) {
	pool, _ := k.GetPool(ctx, poolId)
	if pool.Disabled {
		return nil, nil, nil, sdkerrors.Wrapf(types.ErrDisabledPool, "pool %d is disabled", poolId)
	}

	pair, _ := k.GetPair(ctx, pool.PairId)
	rx, ry := k.getPoolBalances(ctx, pool, pair)
	ps := k.GetPoolCoinSupply(ctx, pool)
	ammPool := amm.NewBasicPool(rx.Amount, ry.Amount, ps)
	if ammPool.IsDepleted() {
		return nil, nil, nil, sdkerrors.Wrapf(types.ErrDepletedPool, "pool %d is depleted", poolId)
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

func oracle(denom string) (uint64, bool) {
	if denom == "ucmdx" {
		return 2000000, true
	} else if denom == "ucgold" {
		return 1800000000, true
	} else if denom == "ucsilver" {
		return 25000000, true
	} else if denom == "ucoil" {
		return 120000000, true
	}
	return 0, false
}

func (k Keeper) GetOraclePrices(ctx sdk.Context, quoteCoinDenom, baseCoinDenom string) (sdk.Dec, string, error) {
	oraclePrice, found := oracle(quoteCoinDenom) // for quote coin
	denom := quoteCoinDenom
	if !found {
		oraclePrice, found = oracle(baseCoinDenom) // for base coin
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

func (k Keeper) GetAggregatedChildPoolContributions(ctx sdk.Context, poolIds []uint64, masterPoolSupplyAddresses []sdk.AccAddress) map[string]sdk.Dec {
	poolSupplyData := make(map[string]sdk.Dec)

	for _, poolId := range poolIds {
		liquidityProvidersDataForPool, found := k.GetPoolLiquidityProvidersData(ctx, poolId)
		if !found {
			continue
		}

		pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, poolId)
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

func (k Keeper) GetFarmingRewardsData(ctx sdk.Context, coinsToDistribute sdk.Coin, liquidityGaugeData incentivestypes.LiquidtyGaugeMetaData) ([]incentivestypes.RewardDistributionDataCollector, error) {

	liquidityProvidersDataForPool, found := k.GetPoolLiquidityProvidersData(ctx, liquidityGaugeData.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrLPDataNotExistsForPool, "data not found for pool id %d", liquidityGaugeData.PoolId)
	}

	pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, liquidityGaugeData.PoolId)
	if err != nil {
		return nil, err
	}

	oraclePrice, denom, err := k.GetOraclePrices(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom)
	if err != nil {
		return nil, err
	}

	quoteCoinPoolBalance, baseCoinPoolBalance := k.getPoolBalances(ctx, *pool, *pair)

	// Logic for master pool mechanism
	if liquidityGaugeData.IsMasterPool {

		masterPoolSupplyAddresses := []sdk.AccAddress{}
		masterPoolSupplies := []sdk.Dec{}
		childPoolSupplies := []sdk.Dec{}
		minMasterChildPoolSupplies := []sdk.Dec{}

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
					masterPoolSupplyAddresses = append(masterPoolSupplyAddresses, addr)
					masterPoolSupplies = append(masterPoolSupplies, sdk.Dec(value))
				}
			}
		}

		childPoolIds := []uint64{}
		if len(liquidityGaugeData.ChildPoolIds) == 0 {
			pools := k.GetAllPools(ctx)
			for _, pool := range pools {
				if pool.Id != liquidityGaugeData.PoolId {
					childPoolIds = append(childPoolIds, pool.Id)
				}
			}
		} else {
			// sanitization
			for _, poolId := range liquidityGaugeData.ChildPoolIds {
				if poolId != liquidityGaugeData.PoolId {
					childPoolIds = append(childPoolIds, poolId)
				}
			}
		}

		chilPoolSuppliesData := k.GetAggregatedChildPoolContributions(ctx, childPoolIds, masterPoolSupplyAddresses)

		for _, accAddress := range masterPoolSupplyAddresses {
			aggregatedSupplyValue, found := chilPoolSuppliesData[accAddress.String()]
			if !found {
				childPoolSupplies = append(childPoolSupplies, sdk.NewDec(0))
			} else {
				childPoolSupplies = append(childPoolSupplies, aggregatedSupplyValue)
			}
		}

		if len(masterPoolSupplyAddresses) != len(masterPoolSupplies) || len(masterPoolSupplyAddresses) != len(childPoolSupplies) {
			return nil, types.ErrSupplyValueCalculationInvalid
		}

		totalRewardEligibleSupply := sdk.NewDec(0)
		for i := 0; i < len(masterPoolSupplyAddresses); i++ {
			minSupply := sdk.Dec{}
			if masterPoolSupplies[i].LTE(childPoolSupplies[i]) {
				minSupply = masterPoolSupplies[i]
			} else {
				minSupply = childPoolSupplies[i]
			}
			totalRewardEligibleSupply = totalRewardEligibleSupply.Add(minSupply)
			minMasterChildPoolSupplies = append(minMasterChildPoolSupplies, minSupply)
		}

		multiplier := sdk.NewDecFromInt(coinsToDistribute.Amount).Quo(totalRewardEligibleSupply)

		rewardData := []incentivestypes.RewardDistributionDataCollector{}
		for index, address := range masterPoolSupplyAddresses {
			if !minMasterChildPoolSupplies[index].IsZero() {
				calculatedReward := int64(math.Floor(minMasterChildPoolSupplies[index].Mul(multiplier).MustFloat64()))
				newData := new(incentivestypes.RewardDistributionDataCollector)
				newData.RewardReceiver = address
				newData.RewardCoin = sdk.NewCoin(coinsToDistribute.Denom, sdk.NewInt(calculatedReward))
				rewardData = append(rewardData, *newData)
			}
		}

		return rewardData, nil
	} else {
		// Logic for exteranl rewards or non masterpool gauges
		// TODO : write logic for external reward calculation
		fmt.Println() // to avoid editor warnings for empty block
	}

	return []incentivestypes.RewardDistributionDataCollector{}, nil
}

func (k Keeper) ValidateMsgTokensSoftLock(ctx sdk.Context, msg *types.MsgTokensSoftLock) (sdk.AccAddress, error) {
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolId, "no pool exists with id : %d", msg.PoolId)
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

	liquidityProvidersData, found := k.GetPoolLiquidityProvidersData(ctx, msg.PoolId)

	if !found {
		liquidityProvidersData = types.PoolLiquidityProvidersData{
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

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolId, "no pool exists with id : %d", msg.PoolId)
	}

	if msg.SoftUnlockCoin.Denom != pool.PoolCoinDenom {
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", &pool.PoolCoinDenom, msg.SoftUnlockCoin.Denom)
	}
	return depositor, nil
}

func (k Keeper) SoftUnlockTokens(ctx sdk.Context, msg *types.MsgTokensSoftUnlock) error {
	depositor, err := k.ValidateMsgTokensSoftUnlock(ctx, msg)
	if err != nil {
		return err
	}

	liquidityProvidersData, found := k.GetPoolLiquidityProvidersData(ctx, msg.PoolId)

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

func (k Keeper) GetMinimumEpochDurationFromPoolId(ctx sdk.Context, poolId uint64, gauges []incentivestypes.Gauge) time.Duration {
	var minEpochDuration time.Duration
	for _, gauge := range gauges {
		switch kind := gauge.Kind.(type) {
		case *incentivestypes.Gauge_LiquidityMetaData:
			if kind.LiquidityMetaData.PoolId == poolId {
				if minEpochDuration == time.Duration(0) {
					minEpochDuration = gauge.TriggerDuration
				} else if gauge.TriggerDuration < minEpochDuration {
					minEpochDuration = gauge.TriggerDuration
				}
			}
		}
	}
	if minEpochDuration == time.Duration(0) {
		minEpochDuration = time.Second * 86400 // if no gauge for given pool, making 24h as default vaule
	}
	return minEpochDuration
}

func (k Keeper) ProcessQueuedLiquidityProviders(ctx sdk.Context) {
	availablePools := k.GetAllPools(ctx)
	availableLiquidityGauges := k.incentivesKeeper.GetAllGaugesByGaugeTypeId(ctx, incentivestypes.LiquidityGaugeTypeId)

	for _, pool := range availablePools {
		poolLpData, found := k.GetPoolLiquidityProvidersData(ctx, pool.Id)
		minEpochDuration := k.GetMinimumEpochDurationFromPoolId(ctx, pool.Id, availableLiquidityGauges)

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
