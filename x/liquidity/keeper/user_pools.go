package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetIndividualUserPoolsData(ctx sdk.Context, address sdk.AccAddress) (userPoolsData types.UserPoolsData, found bool) {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.UserPoolDataKey(address)
		value = store.Get(key)
	)

	if value == nil {
		return userPoolsData, false
	}

	k.cdc.MustUnmarshal(value, &userPoolsData)

	return userPoolsData, true
}

func (k *Keeper) SetIndividualUserPoolsData(ctx sdk.Context, userPoolsData types.UserPoolsData) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.UserPoolDataKey(sdk.AccAddress(userPoolsData.UserAddress))
		value = k.cdc.MustMarshal(&userPoolsData)
	)
	store.Set(key, value)
}

//For a specific user , to check if a certain pool data exists in the data structre
func (k *Keeper) GetUserPoolsContributionData(userPoolsData types.UserPoolsData, poolId uint64) (found bool) {

	for _, pool := range userPoolsData.UserPoolWiseData {
		if pool.PoolId == poolId {
			return true

		} else {
			continue
		}
	}
	return false

}
func (k *Keeper) UpdateUnbondedTokensUserPoolData(userPoolsData types.UserPoolsData, poolId uint64, unbondedTokens sdk.Int) (updatedUserPoolsData types.UserPoolsData) {
	for _, poolData := range userPoolsData.UserPoolWiseData {
		if poolData.PoolId == poolId {
			updatedTokens := poolData.UnbondedPoolCoin.Add(unbondedTokens)
			fmt.Println("Checking updated data----1", unbondedTokens)
			poolData.UnbondedPoolCoin = &updatedTokens
			fmt.Println("Checking updated data----2", poolData.UnbondedPoolCoin)

			// fmt.Println(reflect.TypeOf(poolData.UnbondedPoolCoin));
		}

	}

	updatedUserPoolsData = userPoolsData
	fmt.Println("Checking updated data-------3", updatedUserPoolsData)
	return updatedUserPoolsData

}

func (k *Keeper) CreatePoolForUser(existinguserPoolsData types.UserPoolsData, poolId uint64, unbondedTokens sdk.Int) (updatedUserPoolsData types.UserPoolsData) {
	var userPoolsData types.UserPools
	// var userUnbondingTokens types.UserPoolUnbondingTokens
	bondedPoolToken := sdk.ZeroInt()
	// unbondingPoolToken := sdk.ZeroInt()
	userPoolsData.PoolId = poolId
	userPoolsData.BondedPoolCoin = &bondedPoolToken
	userPoolsData.UnbondedPoolCoin = &unbondedTokens
	// userUnbondingTokens.IsUnbondingPoolCoin = &unbondingPoolToken
	// userUnbondingTokens.UnbondingStartTime = "0"
	// userUnbondingTokens.UnbondingEndTime = "0"
	// userPoolsData.UserPoolUnbondingTokens = append(userPoolsData.UserPoolUnbondingTokens, &userUnbondingTokens)
	existinguserPoolsData.UserPoolWiseData = append(existinguserPoolsData.UserPoolWiseData, &userPoolsData)

	updatedUserPoolsData = existinguserPoolsData
	return updatedUserPoolsData

}

func (k *Keeper) CalculateUnbondingEndTime(currentTime int64) (endTime float64) {

	//Taking Default UNbonding timline from Params

	defaultUnbondingPeriod := types.DefaultPoolUnbondingDuration
	//COnverting it to float
	value := float64(defaultUnbondingPeriod.Int64())
	//Calculating hours in a day & seconds
	hoursInADay, _ := time.ParseDuration("24h")
	secondsInADay := hoursInADay.Seconds()
	//CAlculating the total unbonding time in seconds - float64
	totalUnbondingDuration := secondsInADay * value
	//Endtime - totalunbonding time + current time
	endTime = float64(currentTime) + totalUnbondingDuration
	fmt.Println("Current Time",currentTime)
	fmt.Println("End Time",endTime)

	return endTime
}

// func (k *Keeper) GetAllUserPools(ctx sdk.Context) (userpools []types.UserPoolsData) {

// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.LockedVaultKeyPrefix)
// 	)

// 	defer iter.Close()

// 	for ; iter.Valid(); iter.Next() {
// 		var locked_vault types.LockedVault
// 		k.cdc.MustUnmarshal(iter.Value(), &locked_vault)
// 		locked_vaults = append(locked_vaults, locked_vault)
// 	}

// 	return locked_vaults

// }

//Functions that will be made
//1. Set USERPOOLSDATA - for a individual user
//2. Get USERPOOLSDATA- for a individual user
//3. Get AllUsersPOOLSDATA- for all users -- this will
//4. Bond User Token - for a individual user
//5. Start Unbonding For User Token- for a individual user
//6. BEgin Blocker - Unbond User Token whose bonding is complete - Automatic Execution
