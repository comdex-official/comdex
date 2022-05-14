package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/vault/types"
)






// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txAddWhiteListedAsset(),
		txCreateLocker(),
		txDepositAsset(),
		txWithdrawAsset(),
	)

	return cmd
}

func txCreateLocker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-locker [app_mapping_id] [asset_id] [amount]",
		Short: "create a new locker",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidAmountIn
			}

			
			msg := types.NewMsgCreateRequest(ctx.FromAddress, appMappingId, assetId, amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [id] [amount]",
		Short: "creates a new deposit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDepositRequest(ctx.FromAddress, id, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [id] [amount]",
		Short: "create a new withdraw",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgWithdrawRequest(ctx.FromAddress, id, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

