package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group locker queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "locker",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "locker"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryParams(),
		queryLockedVault(),
		queryLockersByAppToAssetID(),
		queryLockerByAppID(),
		queryTotalDepositByAppAndAssetID(),
		queryOwnerLockerByAppIDbyOwner(),
		queryOwnerLockerOfAllAppsByOwner(),
		queryOwnerLockerByAppToAssetIDbyOwner(),
		queryTotalLockerByAppID(),
		queryTotalLockerByAppToAssetID(),
		queryWhiteListedAssetIDsByAppID(),
		queryWhiteListedAssetByAllApps(),
		queryLockerLookupTableByApp(),
		queryLockerLookupTableByAppAndAssetID(),
		queryLockerTotalDepositedByApp(),
		queryOwnerTxDetailsLockerOfAppByOwnerByAsset(),
		queryLockerByAppByOwner(),
		queryLockerTotalRewardsByAssetAppWise(),
	)

	return cmd
}

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "locker",
		Short:                      fmt.Sprintf("%s transactions subcommands", "locker"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(

		txCreateLocker(),
		txDepositAssetLocker(),
		txWithdrawAssetLocker(),
		txCloseLocker(),
		txlockerRewardCalc(),
	)

	return cmd
}
