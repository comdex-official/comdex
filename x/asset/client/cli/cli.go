package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "asset",
		Short:                      "Asset module sub-commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryAsset(),
		queryAssets(),
		queryPair(),
		queryPairs(),
		queryApp(),
		queryExtendedPairVault(),
		queryAllExtendedPairVaults(),
		queryApps(),
		queryAllExtendedPairVaultsByApp(),
		queryAllExtendedPairStableVaultsIDByApp(),
		queryGovTokenByApp(),
		queryAllExtendedPairStableVaultsByApp(),
	)

	return cmd
}

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "asset",
		Short:                      "asset module sub-commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand()

	return cmd
}
