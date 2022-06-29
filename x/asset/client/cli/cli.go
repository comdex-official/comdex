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
		queryParams(),
		queryAppMappings(),
		queryExtendedPairVault(),
		queryAllExtendedPairVaults(),
		queryAppsMappings(),
		queryAllExtendedPairVaultsByApp(),
		queryAllExtendedPairStableVaultsIdByApp(),
		queryGovTokenByApp(),
		queryAllExtendedPairStableVaultsDataByApp(),
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
