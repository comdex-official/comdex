package cli

import (
	"fmt"
	"github.com/comdex-official/comdex/x/collector/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "collector",
		Short:                      fmt.Sprintf("%s transactions subcommands", "collector"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		txDeposit(),
	)
	return cmd
}
func txDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [amount] [app-id]",
		Short: "deposit cmst into collector",
		Long: `This transaction is exclusively designed for depositing CMST tokens into the collector module. We strongly advise regular users against using this transaction, 
               as it could lead to potential fund loss and cannot be reversed. Its sole purpose is to facilitate refunds to stAtom vault owners by administrators.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(ctx.GetFromAddress().String(), asset, appID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
