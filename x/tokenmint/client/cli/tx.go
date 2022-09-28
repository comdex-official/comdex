package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

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
			appID, err := sdk.ParseUint(args[0])
			if err != nil {
				return err
			}
			assetID, err := sdk.ParseUint(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgMintNewTokensRequest(ctx.GetFromAddress().String(), appID.Uint64(), assetID.Uint64())

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
