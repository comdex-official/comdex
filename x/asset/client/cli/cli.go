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

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "asset module sub-commands",
	}

	cmd.AddCommand(
		txAddAsset(),
		txUpdateAsset(),
		txAddMarket(),
		txUpdateMarket(),
		txAddMarketForAsset(),
		txRemoveMarketForAsset(),
		txAddPair(),
		txUpdatePair(),
		txFetchPrice(),
	)

	return cmd
}
