package keeper

import (
	"math"
	"strconv"
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

	var (
		baseCoinPoolPrice  sdk.Dec
		quoteCoinPoolPrice sdk.Dec
	)
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
			activeFarmer, found := k.GetActiveFarmer(ctx, appID, poolID, address)
			if !found {
				continue
			}
			x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, activeFarmer.FarmedPoolCoin)
			if err != nil {
				continue
			}
			quoteCoin := sdk.NewCoin(pair.QuoteCoinDenom, x)
			baseCoin := sdk.NewCoin(pair.BaseCoinDenom, y)
			value := k.CalculateLiquidityAddedValue(ctx, quoteCoinPoolBalance, baseCoinPoolBalance, quoteCoin, baseCoin, oraclePrice, denom)
			_, found = poolSupplyData[address.String()]
			if !found {
				poolSupplyData[address.String()] = value
			} else {
				poolSupplyData[address.String()] = poolSupplyData[address.String()].Add(value)
			}
		}
	}
	return poolSupplyData
}

func (k Keeper) GetFarmingRewardsData(ctx sdk.Context, appID uint64, coinsToDistribute sdk.Coin, liquidityGaugeData rewardstypes.LiquidtyGaugeMetaData) ([]rewardstypes.RewardDistributionDataCollector, error) {
	pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, appID, liquidityGaugeData.PoolId)
	if err != nil {
		return nil, err
	}

	oraclePrice, denom, err := k.GetOraclePrices(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom)
	if err != nil {
		return nil, err
	}

	quoteCoinPoolBalance, baseCoinPoolBalance := k.getPoolBalances(ctx, *pool, *pair)

	var lpAddresses []sdk.AccAddress
	var lpSupplies []sdk.Dec

	activeFarmers := k.GetAllActiveFarmers(ctx, appID, pool.Id)
	for _, activeFarmer := range activeFarmers {
		addr, err := sdk.AccAddressFromBech32(activeFarmer.Farmer)
		if err != nil {
			continue
		}
		x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, activeFarmer.FarmedPoolCoin)
		if err != nil {
			continue
		}
		quoteCoin := sdk.NewCoin(pair.QuoteCoinDenom, x)
		baseCoin := sdk.NewCoin(pair.BaseCoinDenom, y)
		value := k.CalculateLiquidityAddedValue(ctx, quoteCoinPoolBalance, baseCoinPoolBalance, quoteCoin, baseCoin, oraclePrice, denom)
		lpAddresses = append(lpAddresses, addr)
		lpSupplies = append(lpSupplies, value)
	}

	// Logic for master pool mechanism
	if liquidityGaugeData.IsMasterPool {
		var childPoolSupplies []sdk.Dec
		var minMasterChildPoolSupplies []sdk.Dec

		var childPoolIds []uint64
		if len(liquidityGaugeData.ChildPoolIds) == 0 {
			pools := k.GetAllPools(ctx, appID)
			for _, pool := range pools {
				if pool.Id != liquidityGaugeData.PoolId && !pool.Disabled {
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

		// if no child pools, than use standard mechanism for reward distribution
		if len(childPoolIds) != 0 {
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

			var rewardData []rewardstypes.RewardDistributionDataCollector
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
	}

	// Logic for non master pool gauges (external rewards), (also used for masterpool if no child pool exists)
	totalRewardEligibleSupply := sdk.NewDec(0)
	for _, supply := range lpSupplies {
		totalRewardEligibleSupply = totalRewardEligibleSupply.Add(supply)
	}

	var rewardData []rewardstypes.RewardDistributionDataCollector
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

func (k Keeper) ValidateMsgFarm(ctx sdk.Context, msg *types.MsgFarm) (sdk.AccAddress, error) {
	farmer, err := sdk.AccAddressFromBech32(msg.Farmer)
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

	if msg.FarmingPoolCoin.Denom != pool.PoolCoinDenom {
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", pool.PoolCoinDenom, msg.FarmingPoolCoin.Denom)
	}
	if !msg.FarmingPoolCoin.Amount.IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrorNotPositiveAmont, "pool coin amount should be positive")
	}
	return farmer, nil
}

func (k Keeper) Farm(ctx sdk.Context, msg *types.MsgFarm) error {
	farmer, err := k.ValidateMsgFarm(ctx, msg)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, farmer, types.ModuleName, sdk.NewCoins(msg.FarmingPoolCoin))
	if err != nil {
		return err
	}

	queuedFarmer, found := k.GetQueuedFarmer(ctx, msg.AppId, msg.PoolId, farmer)

	if !found {
		queuedFarmer = types.NewQueuedfarmer(msg.AppId, msg.PoolId, farmer)
	}

	queuedFarmer.QueudCoins = append(
		queuedFarmer.QueudCoins,
		&types.QueuedCoin{
			FarmedPoolCoin: msg.FarmingPoolCoin,
			CreatedAt:      ctx.BlockTime(),
		},
	)
	k.SetQueuedFarmer(ctx, queuedFarmer)

	ctx.GasMeter().ConsumeGas(types.FarmGas, "FarmGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFarm,
			sdk.NewAttribute(types.AttributeKeyFarmer, msg.Farmer),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolCoin, msg.FarmingPoolCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTimeStamp, ctx.BlockTime().String()),
		),
	})

	return nil
}

func (k Keeper) ValidateMsgUnfarm(ctx sdk.Context, msg *types.MsgUnfarm) (sdk.AccAddress, error) {
	farmer, err := sdk.AccAddressFromBech32(msg.Farmer)
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

	if msg.UnfarmingPoolCoin.Denom != pool.PoolCoinDenom {
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", pool.PoolCoinDenom, msg.UnfarmingPoolCoin.Denom)
	}
	if !msg.UnfarmingPoolCoin.Amount.IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrorNotPositiveAmont, "pool coin amount should be positive")
	}
	return farmer, nil
}

func (k Keeper) Unfarm(ctx sdk.Context, msg *types.MsgUnfarm) error {
	farmer, err := k.ValidateMsgUnfarm(ctx, msg)
	if err != nil {
		return err
	}

	activeFarmer, afound := k.GetActiveFarmer(ctx, msg.AppId, msg.PoolId, farmer)
	queuedFarmer, qfound := k.GetQueuedFarmer(ctx, msg.AppId, msg.PoolId, farmer)

	if !afound && !qfound {
		return sdkerrors.Wrapf(types.ErrorFarmerNotFound, "no active farm found for given pool id %d", msg.PoolId)
	}

	farmedCoinAmount := sdk.NewInt(0)

	if qfound {
		for _, qCoin := range queuedFarmer.QueudCoins {
			farmedCoinAmount = farmedCoinAmount.Add(qCoin.FarmedPoolCoin.Amount)
		}
	}

	if afound {
		farmedCoinAmount = farmedCoinAmount.Add(activeFarmer.FarmedPoolCoin.Amount)
	}

	if farmedCoinAmount.LT(msg.UnfarmingPoolCoin.Amount) {
		return sdkerrors.Wrapf(types.ErrInvalidUnfarmAmount, "farmed pool coin amount %d%s smaller than requested unfarming pool coin amount %d%s", farmedCoinAmount.Int64(), msg.UnfarmingPoolCoin.Denom, msg.UnfarmingPoolCoin.Amount.Int64(), msg.UnfarmingPoolCoin.Denom)
	}

	unFarmingCoin := msg.UnfarmingPoolCoin
	queuedCoins := queuedFarmer.QueudCoins
	if qfound {
		for i := len(queuedCoins) - 1; i >= 0; i-- {
			if queuedCoins[i].FarmedPoolCoin.Amount.GTE(msg.UnfarmingPoolCoin.Amount) {
				queuedCoins[i].FarmedPoolCoin.Amount = queuedCoins[i].FarmedPoolCoin.Amount.Sub(msg.UnfarmingPoolCoin.Amount)
				msg.UnfarmingPoolCoin.Amount = sdk.NewInt(0)
				break
			} else {
				msg.UnfarmingPoolCoin.Amount = msg.UnfarmingPoolCoin.Amount.Sub(queuedCoins[i].FarmedPoolCoin.Amount)
				queuedCoins[i].FarmedPoolCoin.Amount = sdk.NewInt(0)
			}
		}
	}

	updatedQueuedCoins := []*types.QueuedCoin{}
	for _, object := range queuedCoins {
		if object.FarmedPoolCoin.Amount.IsZero() {
			break
		} else {
			updatedQueuedCoins = append(updatedQueuedCoins, object)
		}
	}

	queuedFarmer.QueudCoins = updatedQueuedCoins

	aFarmerUpdated := false
	if !msg.UnfarmingPoolCoin.Amount.IsZero() {
		aFarmerUpdated = true
		activeFarmer.FarmedPoolCoin.Amount = activeFarmer.FarmedPoolCoin.Amount.Sub(msg.UnfarmingPoolCoin.Amount)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, farmer, sdk.NewCoins(unFarmingCoin))
	if err != nil {
		return err
	}

	if aFarmerUpdated {
		k.SetActiveFarmer(ctx, activeFarmer)
	}
	k.SetQueuedFarmer(ctx, queuedFarmer)

	ctx.GasMeter().ConsumeGas(types.UnfarmGas, "UnfarmGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnfarm,
			sdk.NewAttribute(types.AttributeKeyFarmer, msg.Farmer),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolCoin, msg.UnfarmingPoolCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTimeStamp, ctx.BlockTime().String()),
		),
	})

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
		minEpochDuration = types.DefaultFarmingQueueDuration // if no gauge for given pool, making 24h as default vaule
	}
	return minEpochDuration
}

func (k Keeper) ProcessQueuedFarmers(ctx sdk.Context, appID uint64) {
	availablePools := k.GetAllPools(ctx, appID)
	availableLiquidityGauges := k.rewardsKeeper.GetAllGaugesByGaugeTypeID(ctx, rewardstypes.LiquidityGaugeTypeID)

	for _, pool := range availablePools {
		minEpochDuration := k.GetMinimumEpochDurationFromPoolID(ctx, pool.Id, availableLiquidityGauges)
		queuedFarmers := k.GetAllQueuedFarmers(ctx, pool.AppId, pool.Id)

		for _, queuedFarmer := range queuedFarmers {

			activeFarmer, found := k.GetActiveFarmer(ctx, queuedFarmer.AppId, queuedFarmer.PoolId, sdk.MustAccAddressFromBech32(queuedFarmer.Farmer))
			if !found {
				activeFarmer = types.NewActivefarmer(queuedFarmer.AppId, queuedFarmer.PoolId, sdk.MustAccAddressFromBech32(queuedFarmer.Farmer), sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(0)))
			}

			updatedQueue := []*types.QueuedCoin{}
			activeFarmUpdated := false
			for _, queuedCoin := range queuedFarmer.QueudCoins {
				if ctx.BlockTime().Before(queuedCoin.CreatedAt.Add(minEpochDuration)) {
					updatedQueue = append(updatedQueue, queuedCoin)
				} else {
					activeFarmUpdated = true
					activeFarmer.FarmedPoolCoin.Amount = activeFarmer.FarmedPoolCoin.Amount.Add(queuedCoin.FarmedPoolCoin.Amount)
				}
			}
			queuedFarmer.QueudCoins = updatedQueue

			if activeFarmUpdated {
				k.SetActiveFarmer(ctx, activeFarmer)
				k.SetQueuedFarmer(ctx, queuedFarmer)
			}
		}
	}
}
