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
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txCreate(),
		txDeposit(),
		txWithdraw(),
		txDrawDebt(),
		txRepayDebt(),
		txClose(),
		txCreateStableMint(),
		txDepositStableMint(),
		txWithdrawStableMint(),
	)

	return cmd
}

func txCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [appMappingID] [extendedPairVaultID] [amount_in] [amount_out]",
		Short: "create a new vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
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

			msg := types.NewMsgCreateRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, amountIn, amountOut)

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
		Use:   "deposit [appMappingID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "creates a new deposit",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			userVaultid := args[2]

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDepositRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, userVaultid, amount)
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
		Use:   "withdraw [appMappingID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "create a new withdraw",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			userVaultid := args[2]

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgWithdrawRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txDrawDebt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw [appMappingID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "draw debt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			userVaultid := args[2]

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgDrawRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txRepayDebt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repay [appMappingID] [extendedPairVaultID] [userVaultid] [amount]",
		Short: "repay debt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			userVaultid := args[2]

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidAmount
			}

			msg := types.NewMsgRepayRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, userVaultid, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close [appMappingID] [extendedPairVaultID] [userVaultid]",
		Short: "close",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairVaultID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			userVaultid := args[2]

			msg := types.NewMsgLiquidateRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, userVaultid)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txCreateStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-stable-mint [appMappingID] [extendedPairVaultID] [amount] ",
		Short: "create a new stable mint vault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
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

			msg := types.NewMsgCreateStableMintRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txDepositStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-stable-mint [appMappingID] [extendedPairVaultID] [amount] [stablemint_id] ",
		Short: "deposit to stable mint vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
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
			stablemintID := args[3]

			msg := types.NewMsgDepositStableMintRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, amount, stablemintID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txWithdrawStableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-stable-mint [appMappingID] [extendedPairVaultID] [amount] [stablemint_id]",
		Short: "withdraw from stable mint vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
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
			stablemintID := args[3]

			msg := types.NewMsgWithdrawStableMintRequest(ctx.FromAddress, appMappingID, extendedPairVaultID, amount, stablemintID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
