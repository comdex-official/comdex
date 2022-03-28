package keeper

import (
	"strings"
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) GetMintingRewardsID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.MintingRewardsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue()
}

func (k *Keeper) SetMintingRewardsID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.MintingRewardsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetMintingRewards(ctx sdk.Context, mintingReward types.MintingRewards) {
	var (
		store = k.Store(ctx)
		key   = types.MintingRewardsKey(mintingReward.Id)
		value = k.cdc.MustMarshal(&mintingReward)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteCollateralAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.MintingRewardsKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetMintingReward(ctx sdk.Context, id uint64) (mintingReward types.MintingRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.MintingRewardsKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return mintingReward, false
	}
	k.cdc.MustUnmarshal(value, &mintingReward)
	return mintingReward, true
}

func (k *Keeper) GetMintingRewards(ctx sdk.Context) (mintingRewards []types.MintingRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.MintingRewardsKeyPrefix)
	)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var mintingReward types.MintingRewards
		k.cdc.MustUnmarshal(iter.Value(), &mintingReward)
		mintingRewards = append(mintingRewards, mintingReward)
	}
	return mintingRewards
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// New minting-rewards from the passed proposal are being added from here.
func (k Keeper) AddNewMintingRewards(ctx sdk.Context, newMintingRewardsData types.MintingRewards) error {
	availableAssets := k.asset.GetAssets(ctx)
	availableAssetsDenoms := []string{}
	for _, asset := range availableAssets {
		availableAssetsDenoms = append(availableAssetsDenoms, asset.Denom)
	}
	assetsDenomInProposal := []string{newMintingRewardsData.AllowedCasset, newMintingRewardsData.AllowedCollateral}
	invalidAssets := []string{}
	for _, asset := range assetsDenomInProposal {
		if !contains(availableAssetsDenoms, asset) {
			invalidAssets = append(invalidAssets, asset)
		}
	}
	if len(invalidAssets) > 0 {
		return sdkerrors.Wrapf(types.ErrorInvalidAssetDenoms, "invalid denoms %s", strings.Join(invalidAssets, ","))
	}
	availableMintingRewards := k.GetMintingRewards(ctx)
	for _, mintingReward := range availableMintingRewards {
		if mintingReward.AllowedCollateral == newMintingRewardsData.AllowedCollateral && mintingReward.AllowedCasset == newMintingRewardsData.AllowedCasset {
			if !mintingReward.IsActive && mintingReward.Depositor == nil {
				return sdkerrors.Wrapf(types.ErrorMintingRewardPairAlreadyExist, "reward pair %s already exists, deposit pending. ", newMintingRewardsData.AllowedCollateral, newMintingRewardsData.AllowedCasset)
			} else if !mintingReward.IsActive && mintingReward.Depositor != nil && mintingReward.StartTimestamp.After(ctx.BlockTime()) {
				return sdkerrors.Wrapf(types.ErrorMintingRewardPairAlreadyExist, "reward pair %s already exists, deposit made and waiting for start time.", newMintingRewardsData.AllowedCollateral, newMintingRewardsData.AllowedCasset)
			} else if mintingReward.IsActive && mintingReward.Depositor != nil && mintingReward.StartTimestamp.Before(ctx.BlockTime()) {
				return sdkerrors.Wrapf(types.ErrorMintingRewardPairAlreadyExist, "reward pair %s already exists, rewards are in progress", newMintingRewardsData.AllowedCollateral, newMintingRewardsData.AllowedCasset)
			}
		}
	}
	newMintingRewardsData.Id = k.GetMintingRewardsID(ctx) + 1
	k.SetMintingRewardsID(ctx, newMintingRewardsData.Id)
	k.SetMintingRewards(ctx, newMintingRewardsData)

	return nil
}

func (k Keeper) TransferDeposits(ctx sdk.Context, mintingRewardsId uint64, from sdk.AccAddress, startTimeStamp time.Time) error {
	mintingReward, found := k.GetMintingReward(ctx, mintingRewardsId)
	if !found {
		return types.ErrorMintingRewardNotFound
	}
	if mintingReward.IsActive {
		return types.ErrorMintingRewardAlreadyActive
	}
	// reward start time should be atleast 10 minutes after the deposit is being made
	if startTimeStamp.Before(ctx.BlockTime().Add(time.Minute * 10)) {
		return types.ErrorInvalidStartTime
	}
	if mintingReward.Depositor != nil {
		if mintingReward.StartTimestamp.After(ctx.BlockTime()) {
			return types.ErrorDepositAlreadyMade
		}
		return types.ErrorMintingRewardExpired
	}
	err := k.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(mintingReward.TotalRewards))
	if err != nil {
		return err
	}
	mintingReward.StartTimestamp = startTimeStamp
	mintingReward.EndTimestamp = startTimeStamp.Add(time.Hour * 24 * time.Duration(mintingReward.DurationDays))
	// mintingReward.EndTimestamp = startTimeStamp.Add(time.Minute * 5)
	mintingReward.AvailableRewards = mintingReward.TotalRewards
	mintingReward.Depositor = from
	k.SetMintingRewards(ctx, mintingReward)
	return nil
}

func (k Keeper) UpdateMintRewardStartTime(ctx sdk.Context, mintingRewardsId uint64, from sdk.AccAddress, newStartTimeStamp time.Time) error {
	mintingReward, found := k.GetMintingReward(ctx, mintingRewardsId)
	if !found {
		return types.ErrorMintingRewardNotFound
	}
	if !mintingReward.Depositor.Equals(from) {
		return types.ErrorUnauthorized
	}
	if mintingReward.IsActive {
		return types.ErrorMintingRewardAlreadyActive
	}
	if mintingReward.StartTimestamp.Before(ctx.BlockTime()) {
		return types.ErrorMintingRewardExpired
	}
	// reward start time should be atleast 10 minutes after the current block time
	if newStartTimeStamp.Before(ctx.BlockTime().Add(time.Minute * 10)) {
		return types.ErrorInvalidStartTime
	}
	mintingReward.StartTimestamp = newStartTimeStamp
	mintingReward.EndTimestamp = newStartTimeStamp.Add(time.Hour * 24 * time.Duration(mintingReward.DurationDays))
	// mintingReward.EndTimestamp = newStartTimeStamp.Add(time.Minute * 5)
	k.SetMintingRewards(ctx, mintingReward)
	return nil
}

func (k Keeper) DisableMintingReward(ctx sdk.Context, mintingRewardsId uint64) error {
	mintingReward, found := k.GetMintingReward(ctx, mintingRewardsId)
	if !found {
		return types.ErrorMintingRewardNotFound
	}
	if !mintingReward.IsActive {
		return types.ErrorMintingRewardAlreadyDisabled
	}
	// Add Event Emitters
	err := k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintingReward.Depositor, sdk.NewCoins(mintingReward.AvailableRewards))
	if err != nil {
		return err
	}
	mintingReward.IsActive = false
	k.SetMintingRewards(ctx, mintingReward)
	return nil
}

func (k Keeper) EnableMintingRewards(ctx sdk.Context) {
	mintingRewards := k.GetMintingRewards(ctx)
	for _, mintingReward := range mintingRewards {
		if !mintingReward.IsActive && mintingReward.Depositor != nil {
			diff := mintingReward.StartTimestamp.Sub(ctx.BlockTime()).Seconds()
			// if the time difference (starttime - current blocktime) is between +-10 seconds, mark the reward as active.
			if diff >= -10 && diff <= 10 {
				mintingReward.IsActive = true
				k.SetMintingRewards(ctx, mintingReward)
			}
		}
	}
}

// Once the minting reward is expired, it is being disable by my marking isActive flag as false.
// If any rewards are left, they are being transferred to the owner again
func (k Keeper) DisableMintingRewards(ctx sdk.Context) {
	mintingRewards := k.GetMintingRewards(ctx)
	for _, mintingReward := range mintingRewards {
		if mintingReward.IsActive && mintingReward.Depositor != nil && mintingReward.EndTimestamp.Before(ctx.BlockTime()) {
			k.DisableMintingReward(ctx, mintingReward.Id)
		}
	}
}

func (k Keeper) CalculateMintRewards(ctx sdk.Context, vault vaulttypes.Vault, mintingRewards types.MintingRewards) (sdk.Dec, error) {
	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return sdk.NewDec(0), assettypes.ErrorPairDoesNotExist
	}
	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return sdk.NewDec(0), assettypes.ErrorAssetDoesNotExist
	}
	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return sdk.NewDec(0), assettypes.ErrorAssetDoesNotExist
	}
	assetPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
	if !found {
		return sdk.NewDec(0), types.ErrorPriceNotFound
	}

	// total $ vaule of the cAsset minted by the given collateral
	currentTotalCassetMintedValue := k.GetCAssetTotalValueMintedForCollateral(ctx, assetIn)

	// formula
	// if currentTotalCassetMintedValue < maxCap (mentioned in proposal for a specific reward) {
	//
	//		rewardAmount =           $ vaule of cAsset minted by Vault				* Daily Allocated Rewards i.e(totalRewards/no. of days)
	// 					   ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
	// 						    maxCap (mentioned in proposal for reward)
	// }
	// else {
	//
	//		rewardAmount =           $ vaule of cAsset minted by Vault				  * Daily Allocated Rewards i.e(totalRewards/no. of days)
	// 					   ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
	// 					   total $ vaule of the cAsset minted by the given collateral
	// }
	//
	//
	// rewardAmount = rewardAmount*1000000 (i.e 1CMDX = 1000000ucmdx )

	divisor := currentTotalCassetMintedValue
	if currentTotalCassetMintedValue.LT(sdk.NewDec(int64(mintingRewards.CassetMaxCap))) {
		divisor = sdk.NewDec(int64(mintingRewards.CassetMaxCap))
	}
	dailyAllocatedRewards := mintingRewards.TotalRewards.Amount.ToDec().Quo(sdk.NewDec(1000000)).Quo((sdk.NewDec(int64(mintingRewards.DurationDays))))
	mintValue := vault.AmountOut.ToDec().Quo(sdk.NewDec(1000000)).Mul(sdk.NewDec(int64(assetPrice)).Quo(sdk.NewDec(1000000)))
	rewardAmount := mintValue.Quo(divisor).Mul(dailyAllocatedRewards).Mul(sdk.NewDec(1000000))
	return rewardAmount, nil
}

func (k Keeper) DistributeRewards(ctx sdk.Context, mintingReward types.MintingRewards) {

	// Get vault id mapping for all the vaults that are being generated by providing a specific given collateral
	// e.g
	// collateralBasedVaults = {
	// 		"CollateralDenom" : "ucmdx",
	// 		"CassetsVaultIdsMap": {
	// 			"ucgold": {
	// 				"VaultIds": [1,2,3,4]
	// 			},
	// 			"ucsilver": {
	// 				"VaultIds": [5,6,7,8]
	// 			},
	// 			"ucoil": {
	// 				"VaultIds": [9,10,11,12]
	// 			},
	// 		}
	// 	}
	// If the reward if for Pair ucmdx-ucgold, than only [1,2,3,4] vaultIds are considered as eligible
	collateralBasedVaults, found := k.GetCollateralBasedVaults(ctx, mintingReward.AllowedCollateral)

	if found && collateralBasedVaults.CassetsVaultIdsMap[mintingReward.AllowedCasset] != nil {
		eligibleVaultIds := collateralBasedVaults.CassetsVaultIdsMap[mintingReward.AllowedCasset].VaultIds
		for _, vaultId := range eligibleVaultIds {

			// get vault object from the store with given vault id
			vault, found := k.GetVault(ctx, vaultId)

			// the vault is only eligible for the reward if the market cap of cAsset is below the proposal mentioned marketcap at the time of minting,
			// also vault needs to pass the specific threshold time to be eligible for the minting rewards.
			if found && vault.MarketCap.LT(sdk.NewDec(int64(mintingReward.CassetMaxCap))) && vault.CreatedAt.Add(time.Second*time.Duration(mintingReward.MinLockupTimeSeconds)).Before(ctx.BlockTime()) {
				rewardEligible, err := k.CalculateMintRewards(ctx, vault, mintingReward)
				if err != nil {
					continue
				}
				rewardCoin := sdk.NewCoin(mintingReward.TotalRewards.Denom, rewardEligible.RoundInt())
				parsedOwner, err := sdk.AccAddressFromBech32(vault.Owner)
				if err != nil {
					continue
				}
				err = k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, parsedOwner, sdk.NewCoins(rewardCoin))
				if err == nil {
					ctx.EventManager().EmitEvents(sdk.Events{
						sdk.NewEvent(
							types.TypeEvtMintRewardDistribution,
							sdk.NewAttribute(types.AttributeReceiver, vault.Owner),
							sdk.NewAttribute(types.AttributeAmount, rewardCoin.String()),
						),
					})
					// reduce the available rewards by the reward which being sent in above step.
					if !mintingReward.AvailableRewards.Amount.Sub(rewardCoin.Amount).IsZero() {
						mintingReward.AvailableRewards = sdk.NewCoin(mintingReward.AvailableRewards.Denom, mintingReward.AvailableRewards.Amount.Sub(rewardCoin.Amount))
					}
				}
			}
		}
		k.SetMintingRewards(ctx, mintingReward)
	}
}

func (k Keeper) TriggerRewards(ctx sdk.Context) {
	// default time in go to parse string time into time.Time
	const layoutTime = "15:04:05"

	params := k.GetParams(ctx)

	// time on which rewards needs to be distributed
	distributionTimeStamp, _ := time.Parse(layoutTime, strings.TrimSpace(params.MintRewardTimestamp))
	currentTimeStamp, _ := time.Parse(layoutTime, ctx.BlockTime().Format(layoutTime))

	// difference in seconds between reward distribution time and current time
	diff := distributionTimeStamp.Sub(currentTimeStamp).Seconds()

	// if the difference is more than 1/2 (half) hour, mark `IsMintingRewardsTriggered` flag as false.
	// this is added, so that rewards are triggered only once for the given time.
	// i.e this flag will be marked as false before 1/2 an hour of reward distribution, and will be marked as true after reward distribution is done.
	if diff >= 1800 && types.IsMintingRewardsTriggered {
		types.IsMintingRewardsTriggered = false
	}

	// tolerance of -6 & +6 is being added, because the distributionTimeStamp will never be equal to block time,
	// since the logic will be executed only at th new block and on an average block time if 6 seconds.
	if diff >= -6 && diff <= 6 && !types.IsMintingRewardsTriggered {
		mintingRewards := k.GetMintingRewards(ctx)
		// there can be multiple minting rewards, which are being added via proposals.
		for _, mintingReward := range mintingRewards {

			// only trigger rewards for those minting rewards which are being marked as active.
			if mintingReward.IsActive {
				k.DistributeRewards(ctx, mintingReward)
			}
		}
		types.IsMintingRewardsTriggered = true
	}
}
