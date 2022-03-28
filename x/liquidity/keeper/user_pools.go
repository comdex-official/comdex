package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//1. get all users addresses
func (k *Keeper) GetUserAddresses(ctx sdk.Context) (usersAddresses types.AllUserAddressesArray) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.UsersAddressesArrayKey(1)
		value = store.Get(key)
	)

	if value == nil {
		return usersAddresses
	}

	k.cdc.MustUnmarshal(value, &usersAddresses)
	return usersAddresses
}

//2. set to all useraddresses
func (k *Keeper) SetUserAddresses(ctx sdk.Context, usersAddresses types.AllUserAddressesArray) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.UsersAddressesArrayKey(1)
		value = k.cdc.MustMarshal(&usersAddresses)
	)

	store.Set(key, value)
}

//3. get all users pool data contribution
//To get all users pool contribution data
func (k *Keeper) GetAllUsersPoolsData(ctx sdk.Context) (usersPoolDataArray []types.UserPoolsData) {

	userContribution := k.GetUserAddresses(ctx)
	for _, userAddress := range userContribution.UserAddresses {
		individualUserData, found := k.GetIndividualUserPoolsData(ctx, sdk.AccAddress(userAddress))
		if !found {
			continue
		}
		usersPoolDataArray = append(usersPoolDataArray, individualUserData)

	}

	return usersPoolDataArray
}

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

	bondedPoolToken := sdk.ZeroInt()
	userPoolsData.PoolId = poolId
	userPoolsData.BondedPoolCoin = &bondedPoolToken
	userPoolsData.UnbondedPoolCoin = &unbondedTokens
	existinguserPoolsData.UserPoolWiseData = append(existinguserPoolsData.UserPoolWiseData, &userPoolsData)
	updatedUserPoolsData = existinguserPoolsData
	return updatedUserPoolsData

}

func (k *Keeper) CalculateUnbondingEndTime(currentTime time.Time) (endTime time.Time) {
	start := currentTime
	endTime = start.AddDate(0, 0, 21)
	return endTime
}
func (k *Keeper) CompletingUnbondingProcess(ctx sdk.Context) {
	allUsersData := k.GetAllUsersPoolsData(ctx)
	for _, individualUser := range allUsersData {
		fmt.Println(individualUser)
		fmt.Println(individualUser.UserPoolWiseData)
		for _, userPoolData := range individualUser.UserPoolWiseData {
			fmt.Println(userPoolData)
			lengthOfUnbondingArrray := len(userPoolData.UserPoolUnbondingTokens)
			if userPoolData.UserPoolUnbondingTokens != nil && lengthOfUnbondingArrray > 0 {
				fmt.Println(userPoolData.UserPoolUnbondingTokens)
				fmt.Println("User has got some unbonding tokens")
				//Check all the unbonding tokens that are in progress right now,
				//and remove them from unbonding , and move them to the unbonded field
				for value, userUnbondingPoolCoins := range userPoolData.UserPoolUnbondingTokens {
					currentTime := time.Now().Unix()
					if !userUnbondingPoolCoins.IsUnbondingPoolCoin.IsZero() && sdk.NewInt(currentTime).GTE(sdk.NewInt(userUnbondingPoolCoins.UnbondingEndTime.Unix())) {
						//Move them to unbonded field
						updatedUnbondedCoins := userPoolData.UnbondedPoolCoin.Add(*userUnbondingPoolCoins.IsUnbondingPoolCoin)
						userPoolData.UnbondedPoolCoin = &updatedUnbondedCoins
						//Need to find a better way to pop this array index from the unbonding field
						zeroVal := sdk.ZeroInt()
						userPoolData.UserPoolUnbondingTokens[value].IsUnbondingPoolCoin = &zeroVal
						// userPoolData.UserPoolUnbondingTokens[value].UnbondingStartTime=nil
						// userPoolData.UserPoolUnbondingTokens[value].UnbondingEndTime=time.Time(zeroVal)
						fmt.Println("Updated-Unbonding struct---", userUnbondingPoolCoins)
					} else {
						continue
					}

				}

			} else {
				fmt.Println("User has got no unbonding unbonding tokens")
				continue
			}
		}
		fmt.Println("Final User Struct----", individualUser)
		k.SetIndividualUserPoolsData(ctx, individualUser)

	}

}

//Functions that will be made

//3. Get AllUsersPOOLSDATA- for all users -- this will -done
//5. Start Unbonding For User Token- for a individual user - done

//Current Pending Tasks:
//1. Setting Unbonding Duration in Params- done ***
//2. Calculating Current Time - Setting in  start time-done ***
//3. Calculating End Time- Setting in End Time-done  ***
//4. Write a Begin Blocker  Function that will change the unbonding Tokens to UNbonded Field-done ****
//5. Writing withdraw CHanges in function for unbonded tokens - done ***
//6. Create Pool Changes - Addding to bond unbond- Verify First
//7. Delete Pool Changes- Checking it how it works & aliging it accordingly
// 8.Query Commands - For All Users
//9.Query Commands - USer Wise
//10. TS Proto Generation For all the above mentioned functions
//11. ENd to ENd Testing
//12. Create a new function which get triggered during execute deposits

//----------
//New Protobuf File-done
//For saving the address of all the users-done
//Key for this in kv store will be a integer-done
//Now everytime user interacts, it will append it to the struct-done
//Get and set function to append a vaule to the kv store-done
//This will be user to automate the begin blocker & rewards distribution mechanism
