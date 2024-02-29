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
		NewCreateGasProviderCmd(),
		NewAuthorizeActorsCmd(),
		NewUpdateGasProviderStatusCmd(),
		NewUpdateGasProviderConfigsCmd(),
		NewBlockConsumerCmd(),
		NewUnblockConsumerCmd(),
	)

	return cmd
}

func NewCreateGasProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gas-provider [fee-denom] [max-fee-usage-per-tx] [max-txs-count-per-consumer] [max-fee-usage-per-consumer] [txs-allowed] [contracts-allowed] [gas-deposit]",
		Args:  cobra.ExactArgs(7),
		Short: "Create a gas provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a gas provider.
Example:
$ %s tx %s create-gas-provider ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder,/comdex.liquidity.v1beta1.MsgMarketOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3,comdex1zh9gzcw3j5jd53ulfjx9lj4088plur7xy3jayndwr7jxrdqhg7jqqsfqzx 10000000000ucmdx --from mykey
$ %s tx %s create-gas-provider ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 10000000000ucmdx --from mykey
$ %s tx %s create-gas-provider ucmdx 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder "" 10000000000ucmdx --from mykey
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

			msg := types.NewMsgCreateGasProvider(
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
		Use:   "update-authorized-actors [gas-provider-id] [actors]",
		Args:  cobra.ExactArgs(2),
		Short: "Update authorized actors of the gas provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update authorized actors of the gas provider.
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

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-provider-id: %w", err)
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
				gasProviderID,
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

func NewUpdateGasProviderStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-provider-status [gas-provider-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update status of the gas provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update status of the gas provider.
Example:
$ %s tx %s update-gas-provider-status 32 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-provider-id: %w", err)
			}

			msg := types.NewMsgUpdateGasProviderStatus(
				gasProviderID,
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

func NewUpdateGasProviderConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-provider-config [gas-provider-id] [max-fee-usage-per-tx] [max-txs-count-per-consumer] [max-fee-usage-per-consumer] [txs-allowed] [contracts-allowed]",
		Args:  cobra.ExactArgs(6),
		Short: "Update configs of the gas provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update configs of the gas provider.
Example:
$ %s tx %s update-gas-provider-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder,/comdex.liquidity.v1beta1.MsgMarketOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3,comdex1zh9gzcw3j5jd53ulfjx9lj4088plur7xy3jayndwr7jxrdqhg7jqqsfqzx --from mykey
$ %s tx %s update-gas-provider-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 --from mykey
$ %s tx %s update-gas-provider-config 1 25000 200 5000000 /comdex.liquidity.v1beta1.MsgLimitOrder "" --from mykey
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

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-provider-id: %w", err)
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

			msg := types.NewMsgUpdateGasProviderConfig(
				gasProviderID,
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
		Use:   "block-consumer [gas-provider-id] [consumer]",
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

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-provider-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBlockConsumer(
				gasProviderID,
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
		Use:   "unblock-consumer [gas-provider-id] [consumer]",
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

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-provider-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnblockConsumer(
				gasProviderID,
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
