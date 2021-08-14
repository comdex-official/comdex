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
		queryPool(),
		queryPools(),
	)

	return cmd
}
