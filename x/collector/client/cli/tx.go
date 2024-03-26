package cli

import (
	"fmt"
	"strconv"
	"github.com/comdex-official/comdex/x/collector/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		txRefund(),
		txUpdateDebtParams(),
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

func txRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund",
		Short: "refund cmst from collector to stAtom vault owners",
		Long:  `This transaction is exclusively designed for refunding CMST tokens from the collector module. Its sole purpose is to facilitate refunds to stAtom vault owners`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRefund(ctx.GetFromAddress().String())

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txUpdateDebtParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-debt-params [app_id] [asset_id] [slots] [debt_threshold] [lot_size] [debt_lot_size] [is_debt_auction]",
		Short: "update debt params",
		Long:  `This transaction is exclusively designed for updating debt params`,
		Args:  cobra.ExactArgs(7),
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

			slots, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			debtThreshold, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid debt threshold")
			}

			lotSize, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return fmt.Errorf("invalid lot size")
			}

			debtLotSize, ok := sdk.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid debt lot size")
			}

			isDebtAuction, err := strconv.ParseBool(args[6])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDebtparams(ctx.GetFromAddress().String(), appID, assetID, slots, debtThreshold, lotSize, debtLotSize, isDebtAuction)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
