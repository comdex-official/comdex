package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/locker/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group locker queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryParams(),
		queryLockedVault(),
		queryLockerByProductAssetID(),
		queryLockerByProductID(),
		queryTotalDepositByProductAssetID(),
		queryOwnerLockerByProductIDbyOwner(),
		queryOwnerLockerOfAllProductbyOwner(),
		queryOwnerLockerByProductToAssetIDbyOwner(),
		queryTotalLockerByProductID(),
		queryTotalLockerByProductToAssetID(),
		queryWhiteListedAssetIDsByProductID(),
		queryWhiteListedAssetByAllProduct(),
		queryLockerLookupTableByApp(),
		queryLockerLookupTableByAppAndAssetId(),
		queryLockerTotalDepositedByApp(),
		queryOwnerTxDetailsLockerOfProductByOwnerByAsset(),
		queryLockerByProductbyOwner(),
		queryState(),
		queryLockerTotalRewardsByAssetAppWise(),
	)
	// this line is used by starport scaffolding # 1

	return cmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(

		txCreateLocker(),
		txDepositAssetLocker(),
		txWithdrawAssetLocker(),
		txAddWhiteListedAssetLocker(),
	)

	return cmd
}
