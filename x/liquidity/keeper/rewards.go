package keeper

import (
	"sort"
	"time"

	incentivestypes "github.com/comdex-official/comdex/x/incentives/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetFarmingRewardsData(ctx sdk.Context, liquidityGaugeData incentivestypes.LiquidtyGaugeMetaData) []incentivestypes.RewardDistributionDataCollector {
	return []incentivestypes.RewardDistributionDataCollector{}
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
		return nil, sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found %s", &pool.PoolCoinDenom, msg.SoftLockCoin.Denom)
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
