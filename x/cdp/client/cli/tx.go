package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/cdp/types"
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
		txCreateCdp(),
		txDeposit(),
		txWithdraw(),
		txDrawDebt(),
		txRepayDebt(),
		txLiquidate(),
	)

	return cmd
}

func txCreateCdp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [collateral] [debt] [collateral_type]",
		Short: "create a new cdp",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			debt, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			ctype := args[2]
			msg := types.NewMsgCreateCDPRequest(ctx.FromAddress, collateral, debt, ctype)

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
	return &cobra.Command{
		Use:   "deposit [owner_addr] [collateral] [collateral_type]",
		Short: "creates a new deposit",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCollateralRequest(owner, collateral, args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags())
		},
	}
}

func txWithdraw() *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [owner-addr] [collateral] [collateral-type]",
		Short: "create a new withdraw",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawCollateralRequest(owner, collateral, args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags())
		},
	}
}

func txDrawDebt() *cobra.Command {
	return &cobra.Command{
		Use:   "drawDebt [collateral] [collateral-type]",
		Short: "draw debt",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			principal, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDrawDebtRequest(owner, args[2], principal)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags())
		},
	}
}

func txRepayDebt() *cobra.Command {
	return &cobra.Command{
		Use:   "repay debt [collateral] [collateral-type]",
		Short: "repay debt",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			payment, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRepayDebtRequest(owner, args[2], payment)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags())
		},
	}
}

func txLiquidate() *cobra.Command {
	return &cobra.Command{
		Use:   "liquidate [owner] [collateral-type]",
		Short: "liquidate",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidateCDPRequest(owner, args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags())
		},
	}
}
