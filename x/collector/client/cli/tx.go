package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/collector/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "collector",
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	return cmd
}
