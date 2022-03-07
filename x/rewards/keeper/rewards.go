package keeper

import (
	"fmt"

	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddNewMintingRecords(ctx sdk.Context, proposalData types.NewMintRewards) error {
	availableParams := k.GetParams(ctx)
	fmt.Println(availableParams)
	availableParams.MintRewards = append(availableParams.MintRewards, &types.MintingRewardsV1{
		AllowedCollateral: proposalData.AllowedCollateral,
		AllowedCassets:    proposalData.AllowedCassets,
		TotalRewards:      proposalData.TotalRewards,
		CassetMaxCap:      proposalData.CassetMaxCap,
		DurationDays:      proposalData.DurationDays,
		IsActive:          false,
	})
	k.SetParams(ctx, availableParams)
	return nil
}
