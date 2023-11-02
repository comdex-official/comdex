package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tokenmint",
		Short:                      fmt.Sprintf("%s transactions subcommands", "tokenmint"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txMint(),
	)

	return cmd
}

// Token mint txs cmd.
func txMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokenmint [app_ID] [asset_id]",
		Short: "mint genesis token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			msg := types.NewMsgMintNewTokensRequest(ctx.GetFromAddress().String(), appID, assetID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
