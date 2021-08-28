package cli

import (
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "Asset module sub-commands",
	}

	cmd.AddCommand(
		queryAsset(),
		queryAssets(),
		queryMarket(),
		queryMarkets(),
		queryPair(),
		queryPairs(),
		queryParams(),
	)

	return cmd
}
