package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/gasless/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gasless",
		Short:                      fmt.Sprintf("%s transactions subcommands", "gasless"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateGasTankCmd(),
		NewAuthorizeActorsCmd(),
		NewUpdateGasTankStatusCmd(),
		NewUpdateGasTankConfigsCmd(),
		NewBlockConsumerCmd(),
		NewUnblockConsumerCmd(),
		NewUpdateGasConsumerLimitCmd(),
	)

	return cmd
}

func NewCreateGasTankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gas-tank [fee-denom] [max-fee-usage-per-tx] [max-txs-count-per-consumer] [max-fee-usage-per-consumer] [txs-allowed] [contracts-allowed] [gas-deposit]",
		Args:  cobra.ExactArgs(7),
		Short: "Create a gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a gas tank.
Example:
$ %s tx %s create-gas-tank ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder,/comdex.liquidity.v1beta1.MsgMarketOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3,comdex1zh9gzcw3j5jd53ulfjx9lj4088plur7xy3jayndwr7jxrdqhg7jqqsfqzx 10000000000ucmdx --from mykey
$ %s tx %s create-gas-tank ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 10000000000ucmdx --from mykey
$ %s tx %s create-gas-tank ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder "" 10000000000ucmdx --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			feeDenom := args[0]
			if err := sdk.ValidateDenom(feeDenom); err != nil {
				return fmt.Errorf("invalid fee denom: %w", err)
			}

			maxFeeUsagePerTx, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-tx: %s", args[1])
			}

			maxTxsCountPerConsumer, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("parse max-txs-count-per-consumer: %w", err)
			}

			maxFeeUsagePerConsumer, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-consumer: %s", args[3])
			}

			txsAllowed, err := ParseStringSliceFromString(args[4], ",")
			if err != nil {
				return err
			}

			contractsAllowed, err := ParseStringSliceFromString(args[5], ",")
			if err != nil {
				return err
			}

			gasDeposit, err := sdk.ParseCoinNormalized(args[6])
			if err != nil {
				return fmt.Errorf("invalid gas-deposit: %w", err)
			}

			msg := types.NewMsgCreateGasTank(
				clientCtx.GetFromAddress(),
				feeDenom,
				maxFeeUsagePerTx,
				maxTxsCountPerConsumer,
				maxFeeUsagePerConsumer,
				txsAllowed,
				contractsAllowed,
				gasDeposit,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewAuthorizeActorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-authorized-actors [gas-tank-id] [actors]",
		Args:  cobra.ExactArgs(2),
		Short: "Update authorized actors of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update authorized actors of the gas tank.
Example:
$ %s tx %s update-authorized-actors 1 comdex1...,comdex2... --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			actors, err := ParseStringSliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			sanitizedActors := []sdk.AccAddress{}
			for _, actor := range actors {
				sanitizedActor, err := sdk.AccAddressFromBech32(actor)
				if err != nil {
					return err
				}
				sanitizedActors = append(sanitizedActors, sanitizedActor)
			}

			msg := types.NewMsgAuthorizeActors(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedActors,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasTankStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-tank-status [gas-tank-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update status of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update status of the gas tank.
Example:
$ %s tx %s update-gas-tank-status 32 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			msg := types.NewMsgUpdateGasTankStatus(
				gasTankID,
				clientCtx.GetFromAddress(),
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasTankConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-tank-config [gas-tank-id] [max-fee-usage-per-tx] [max-txs-count-per-consumer] [max-fee-usage-per-consumer] [txs-allowed] [contracts-allowed]",
		Args:  cobra.ExactArgs(6),
		Short: "Update configs of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update configs of the gas tank.
Example:
$ %s tx %s update-gas-tank-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder,/comdex.liquidity.v1beta1.MsgMarketOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3,comdex1zh9gzcw3j5jd53ulfjx9lj4088plur7xy3jayndwr7jxrdqhg7jqqsfqzx --from mykey
$ %s tx %s update-gas-tank-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 --from mykey
$ %s tx %s update-gas-tank-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder "" --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			maxFeeUsagePerTx, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-tx: %s", args[1])
			}

			maxTxsCountPerConsumer, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("parse max-txs-count-per-consumer: %w", err)
			}

			maxFeeUsagePerConsumer, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-consumer: %s", args[3])
			}

			txsAllowed, err := ParseStringSliceFromString(args[4], ",")
			if err != nil {
				return err
			}

			contractsAllowed, err := ParseStringSliceFromString(args[5], ",")
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateGasTankConfig(
				gasTankID,
				clientCtx.GetFromAddress(),
				maxFeeUsagePerTx,
				maxTxsCountPerConsumer,
				maxFeeUsagePerConsumer,
				txsAllowed,
				contractsAllowed,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewBlockConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-consumer [gas-tank-id] [consumer]",
		Args:  cobra.ExactArgs(2),
		Short: "Block consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Block consumer.
Example:
$ %s tx %s block-consumer 1 comdex1.. --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBlockConsumer(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUnblockConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unblock-consumer [gas-tank-id] [consumer]",
		Args:  cobra.ExactArgs(2),
		Short: "Unblock consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unblock consumer.
Example:
$ %s tx %s unblock-consumer 1 comdex1.. --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnblockConsumer(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasConsumerLimitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-consumer-limit [gas-tank-id] [consumer] [total-txs-allowed] [total-fee-consumption-allowed]",
		Args:  cobra.ExactArgs(4),
		Short: "Update consumer consumption limit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update consumer consumption limit.
Example:
$ %s tx %s update-consumer-limit 1 comdex1.. 200 5000000 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			totalTxsAllowed, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("parse total-txs-allowed: %w", err)
			}

			totalFeeConsumptionAllowed, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid total-fee-consumption-allowed: %s", args[3])
			}

			msg := types.NewMsgUpdateGasConsumerLimit(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
				totalTxsAllowed,
				totalFeeConsumptionAllowed,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
