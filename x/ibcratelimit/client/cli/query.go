package cli

import (
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/osmoutils"
	"github.com/comdex-official/comdex/x/ibcratelimit/client/queryproto"
	"github.com/comdex-official/comdex/x/ibcratelimit/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := osmoutils.QueryIndexCmd(types.ModuleName)

	cmd.AddCommand(
		osmoutils.GetParams[*queryproto.ParamsRequest](
			types.ModuleName, queryproto.NewQueryClient),
	)

	return cmd
}
