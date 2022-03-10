package keeper

import (
	"fmt"

	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddNewMintingRewards(ctx sdk.Context, newMintingRewardsData types.MintingRewards) error {
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	fmt.Println("data received.....", newMintingRewardsData)
	// fmt.Println("Print address...", k.account.GetModuleAddress(types.ModuleName))
	// fmt.Println("Print assets...", k.asset.GetAssets(ctx))
	// fmt.Println("Print account...", k.account.GetModuleAccount(ctx, types.ModuleName))
	return nil
}
