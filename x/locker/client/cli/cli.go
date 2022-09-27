package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/locker/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
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
		txCloseLocker(),
		txlockerRewardCalc(),
	)

	return cmd
}
