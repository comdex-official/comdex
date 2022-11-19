package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/petrichormoney/petri/x/esm/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "esm",
		Short:                      fmt.Sprintf("%s transactions subcommands", "esm"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		txDepositESM(),
		txExecuteESM(),
		KillSwitch(),
		CollateralRedemption(),
	)
	return cmd
}

func txDepositESM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-esm [app_id] [amount]",
		Short: "deposit into esm",
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

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(ctx.GetFromAddress().String(), appID, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txExecuteESM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-esm [app_id]",
		Short: "execute esm",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgExecute(ctx.GetFromAddress().String(), appID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func KillSwitch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-all-actions [app_id] [breaker_enable]",
		Short: "Stop/Start all actions of an App",
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

			breakerEnable, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			var switchParams types.KillSwitchParams
			switchParams.AppId = appID
			switchParams.BreakerEnable = breakerEnable

			msg := types.NewMsgKillRequest(ctx.FromAddress, switchParams)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CollateralRedemption() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-collateral [app_id] [amount]",
		Short: "redeem collateral of an App",
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
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgCollateralRedemption(appID, amount, ctx.FromAddress)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
