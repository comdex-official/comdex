package cli

import (
	"github.com/spf13/cobra"

	comdexcli "github.com/comdex-official/comdex/types/comdex_cli"
	"github.com/comdex-official/comdex/x/ibc-rate-limit/client/queryproto"
	"github.com/comdex-official/comdex/x/ibc-rate-limit/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := comdexcli.QueryIndexCmd(types.ModuleName)

	cmd.AddCommand(
		comdexcli.GetParams[*queryproto.ParamsRequest](
			types.ModuleName, queryproto.NewQueryClient),
	)

	return cmd
}
