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

// GetTxCmd returns the transaction commands for this module .
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "vault",
		Short:                      fmt.Sprintf("%s transactions subcommands", "vault"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		Create(),
		Deposit(),
		Withdraw(),
		DrawDebt(),
		RepayDebt(),
		Close(),
		DepositAndDraw(),
		CreateStableMint(),
		DepositStableMint(),
		WithdrawStableMint(),
		CalculateInterest(),
		StableMintWithdrawControl(),
	)

	return cmd
}

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [appID] [extendedPairVaultID] [amount_in] [amount_out]",
		Short: "create a new vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amountIn, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidAmountIn
			}

			amountOut, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmountOut
			}

			msg := types.NewMsgCreateRequest(ctx.FromAddress, appID, extendedPairVaultID, amountIn, amountOut)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func Deposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [appID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "creates a new deposit",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDepositRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func Withdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [appID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "create a new withdraw",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgWithdrawRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func DrawDebt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw [appID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "draw debt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDrawRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func RepayDebt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repay [appID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "repay debt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgRepayRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func Close() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close [appID] [extendedPairVaultID] [userVaultid]",
		Short: "close",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidateRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func DepositAndDraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-draw [appID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "creates a new deposit and draw",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			userVaultid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDepositAndDrawRequest(ctx.FromAddress, appID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CreateStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-stable-mint [appID] [extendedPairVaultID] [amount] ",
		Short: "create a new stable mint vault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidAmountIn
			}

			msg := types.NewMsgCreateStableMintRequest(ctx.FromAddress, appID, extendedPairVaultID, amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func DepositStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-stable-mint [appID] [extendedPairVaultID] [amount] [stablemint_id] ",
		Short: "deposit to stable mint vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidAmountIn
			}

			stablemintID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositStableMintRequest(ctx.FromAddress, appID, extendedPairVaultID, amount, stablemintID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func WithdrawStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-stable-mint [appID] [extendedPairVaultID] [amount] [stablemint_id]",
		Short: "withdraw from stable mint vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidAmountIn
			}

			stablemintID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawStableMintRequest(ctx.FromAddress, appID, extendedPairVaultID, amount, stablemintID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CalculateInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate-interest [appID] [userVaultID]",
		Short: "calculate interest by app and user vault id",
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

			userVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgVaultInterestCalcRequest(ctx.FromAddress, appID, userVaultID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func StableMintWithdrawControl() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stablemint-control [appID]",
		Short: "control stable mint vault withdrawal",
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

			msg := types.NewMsgWithdrawStableMintControlRequest(ctx.FromAddress, appID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}