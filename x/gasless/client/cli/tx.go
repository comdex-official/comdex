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
