package v6

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CosMints struct {
	Address     string `json:"address"`
	Amountucmdx string `json:"amount"`
}

var (
	bech32OperatorAddress  = "comdexvaloper1ndslxsucavg3eglqe4mzge74tdx67rcnd7dawq"
	bech32ConcensusAddress = "comdexvalcons1uyn6l9vhh2yqaey9m2yuldj8vdhc2xl95g663q"
)

func mintLostTokens(
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	operatorAddress sdk.ValAddress,
) {
	var cosMints []CosMints
	err := json.Unmarshal([]byte(recordsJSONString), &cosMints)
	if err != nil {
		panic(fmt.Sprintf("error reading COS JSON: %+v", err))
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

		validator, found := stakingKeeper.GetValidator(ctx, operatorAddress)
		if !found {
			panic(fmt.Sprintf("cos validator '%s' not found", operatorAddress))
		}

		_, err = stakingKeeper.Delegate(ctx, delegatorAddress, coin.Amount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			panic(fmt.Sprintf("error delegating minted %sucmdx from %s to %s: %+v", mintRecord.Amountucmdx, mintRecord.Address, operatorAddress.String(), err))
		}
	}
}
