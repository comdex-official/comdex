package v510

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CosMints struct {
	Address     string `json:"address"`
	Amountucmdx string `json:"amount"`
}

var (
	cosValidatorAddress = "comdexvaloper1g9wqptyaxlkzaryt8dezq4eed566kkfpreuq9y"
	cosConsensusAddress = "comdexvalcons1y2re6msd4xjva7yg4fd573khz4wh3x3pkp3m9x"
)

func mintLostTokens(
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
) {
	var cosMints []CosMints
	err := json.Unmarshal([]byte(recordsJSONString), &cosMints)
	if err != nil {
		panic(fmt.Sprintf("error reading COS JSON: %+v", err))
	}

	cosValAddress, err := sdk.ValAddressFromBech32(cosValidatorAddress)
	if err != nil {
		panic(fmt.Sprintf("validator address is not valid bech32: %s", cosValAddress))
	}

	for _, mintRecord := range cosMints {
		coinAmount, mintOk := sdk.NewIntFromString(mintRecord.Amountucmdx)
		if !mintOk {
			panic(fmt.Sprintf("error parsing mint of %sucmdx to %s", mintRecord.Amountucmdx, mintRecord.Address))
		}

		coin := sdk.NewCoin("ucmdx", coinAmount)
		coins := sdk.NewCoins(coin)

		err = mintKeeper.MintCoins(ctx, coins)
		if err != nil {
			panic(fmt.Sprintf("error minting %sucmdx to %s: %+v", mintRecord.Amountucmdx, mintRecord.Address, err))
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(mintRecord.Address)
		if err != nil {
			panic(fmt.Sprintf("error converting human address %s to sdk.AccAddress: %+v", mintRecord.Address, err))
		}

		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delegatorAddress, coins)
		if err != nil {
			panic(fmt.Sprintf("error sending minted %sucmdx to %s: %+v", mintRecord.Amountucmdx, mintRecord.Address, err))
		}

		cosValidator, found := stakingKeeper.GetValidator(ctx, cosValAddress)
		if !found {
			panic(fmt.Sprintf("cos validator '%s' not found", cosValAddress))
		}

		_, err = stakingKeeper.Delegate(ctx, delegatorAddress, coin.Amount, stakingtypes.Unbonded, cosValidator, true)
		if err != nil {
			panic(fmt.Sprintf("error delegating minted %sucmdx from %s to %s: %+v", mintRecord.Amountucmdx, mintRecord.Address, cosValidatorAddress, err))
		}
	}
}

func revertTombstone(ctx sdk.Context, slashingKeeper slashingkeeper.Keeper) error {
	cosValAddress, err := sdk.ValAddressFromBech32(cosValidatorAddress)
	if err != nil {
		panic(fmt.Sprintf("validator address is not valid bech32: %s", cosValAddress))
	}

	cosConsAddress, err := sdk.ConsAddressFromBech32(cosConsensusAddress)
	if err != nil {
		panic(fmt.Sprintf("consensus address is not valid bech32: %s", cosValAddress))
	}

	// Revert Tombstone info
	slashingKeeper.RevertTombstone(ctx, cosConsAddress)

	// Set jail until=now, the validator then must unjail manually
	slashingKeeper.JailUntil(ctx, cosConsAddress, ctx.BlockTime())

	return nil
}

func RevertCosTombstoning(
	ctx sdk.Context,
	slashingKeeper slashingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) error {
	err := revertTombstone(ctx, slashingKeeper)
	if err != nil {
		return err
	}

	mintLostTokens(ctx, bankKeeper, stakingKeeper, mintKeeper)

	return nil
}
