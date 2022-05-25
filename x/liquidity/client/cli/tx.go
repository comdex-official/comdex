package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreatePairCmd(),
		NewCreatePoolCmd(),
		NewDepositCmd(),
		NewWithdrawCmd(),
		NewLimitOrderCmd(),
		NewMarketOrderCmd(),
		NewCancelOrderCmd(),
		NewCancelAllOrdersCmd(),
		NewSoftLockTokensCmd(),
		NewSoftUnlockTokensCmd(),
	)

	return cmd
}

func NewCreatePairCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pair [base-coin-denom] [quote-coin-denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Create a pair(market) for trading",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a pair(market) for trading.
Example:
$ %s tx %s create-pair uatom stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			baseCoinDenom := args[0]
			quoteCoinDenom := args[1]

			msg := types.NewMsgCreatePair(clientCtx.GetFromAddress(), baseCoinDenom, quoteCoinDenom)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCreatePoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pair-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Create a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a liquidity pool with coins.
Example:
$ %s tx %s create-pool 1 1000000000uatom,50000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pairID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid deposit coints: %w", err)
			}

			msg := types.NewMsgCreatePool(clientCtx.GetFromAddress(), pairID, depositCoins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit coins to a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to a liquidity pool.
Example:
$ %s tx %s deposit 1 1000000000uatom,50000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid deposit coins: %w", err)
			}

			msg := types.NewMsgDeposit(clientCtx.GetFromAddress(), poolID, depositCoins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw coins from the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins from the specified liquidity pool.
Example:
$ %s tx %s withdraw 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			poolCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(
				clientCtx.GetFromAddress(),
				poolID,
				poolCoin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewLimitOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit-order [pair-id] [direction] [offer-coin] [demand-coin-denom] [price] [amount]",
		Args:  cobra.ExactArgs(6),
		Short: "Make a limit order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a limit order.
Example:
$ %s tx %s limit-order 1 buy 5000stake uatom 0.5 10000 --from mykey
$ %s tx %s limit-order 1 b 5000stake uatom 0.5 10000 --from mykey
$ %s tx %s limit-order 1 sell 10000uatom stake 2.0 10000 --order-lifespan=10m --from mykey
$ %s tx %s limit-order 1 s 10000uatom stake 2.0 10000 --order-lifespan=10m --from mykey

[pair-id]: pair id to swap with
[direction]: order direction (one of: buy,b,sell,s)
[offer-coin]: the amount of offer coin to swap
[demand-coin-denom]: the denom to exchange with the offer coin
[price]: the limit order price for the swap; the exchange ratio is X/Y where X is the amount of quote coin and Y is the amount of base coin
[amount]: the amount of base coin to buy or sell
`,
				version.AppName, types.ModuleName,
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

			pairID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			dir, err := parseOrderDirection(args[1])
			if err != nil {
				return fmt.Errorf("parse order direction: %w", err)
			}

			offerCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid offer coin: %w", err)
			}

			demandCoinDenom := args[3]
			if err := sdk.ValidateDenom(demandCoinDenom); err != nil {
				return fmt.Errorf("invalid demand coin denom: %w", err)
			}

			price, err := sdk.NewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid price: %w", err)
			}

			amt, ok := sdk.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[5])
			}

			orderLifespan, _ := cmd.Flags().GetDuration(FlagOrderLifespan)

			msg := types.NewMsgLimitOrder(
				clientCtx.GetFromAddress(),
				pairID,
				dir,
				offerCoin,
				demandCoinDenom,
				price,
				amt,
				orderLifespan,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetOrder())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewMarketOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market-order [pair-id] [direction] [offer-coin] [demand-coin-denom] [amount]",
		Args:  cobra.ExactArgs(5),
		Short: "Make a market order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a market order.
Example:
$ %s tx %s market-order 1 buy 5000stake uatom 10000 --from mykey
$ %s tx %s market-order 1 b 5000stake uatom 10000 --from mykey
$ %s tx %s market-order 1 sell 10000uatom stake 10000 --order-lifespan=10m --from mykey
$ %s tx %s market-order 1 s 10000uatom stake 10000 --order-lifespan=10m --from mykey

[pair-id]: pair id to swap with
[direction]: order direction (one of: buy,b,sell,s)
[offer-coin]: the amount of offer coin to swap
[demand-coin-denom]: the denom to exchange with the offer coin
[amount]: the amount of base coin to buy or sell
`,
				version.AppName, types.ModuleName,
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

			pairID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			dir, err := parseOrderDirection(args[1])
			if err != nil {
				return fmt.Errorf("parse order direction: %w", err)
			}

			offerCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid offer coin: %w", err)
			}

			demandCoinDenom := args[3]
			if err := sdk.ValidateDenom(demandCoinDenom); err != nil {
				return fmt.Errorf("invalid demand coin denom: %w", err)
			}

			amt, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[4])
			}

			orderLifespan, _ := cmd.Flags().GetDuration(FlagOrderLifespan)

			msg := types.NewMsgMarketOrder(
				clientCtx.GetFromAddress(),
				pairID,
				dir,
				offerCoin,
				demandCoinDenom,
				amt,
				orderLifespan,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetOrder())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-order [pair-id] [order-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel an order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel an order.
Example:
$ %s tx %s cancel-order 1 1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pairID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrder(
				clientCtx.GetFromAddress(),
				pairID,
				orderID,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelAllOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-all-orders [pair-ids]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Cancel all orders",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel all orders.
Example:
$ %s tx %s cancel-all-orders --from mykey
$ %s tx %s cancel-all-orders 1,3 --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var pairIDs []uint64
			for _, pairIDStr := range strings.Split(args[0], ",") {
				pairID, err := strconv.ParseUint(pairIDStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
				pairIDs = append(pairIDs, pairID)
			}

			msg := types.NewMsgCancelAllOrders(clientCtx.GetFromAddress(), pairIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSoftLockTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soft-lock [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "soft-lock coins from the specified liquidity pool, to start earning rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`soft-lock coins from the specified liquidity pool,  to start earning rewards
Example:
$ %s tx %s soft-lock 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			softLockCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSoftLock(
				clientCtx.GetFromAddress(),
				poolID,
				softLockCoin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSoftUnlockTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soft-unlock [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "soft-unlock coins from the specified liquidity pool, to stop receiving rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`soft-unlock coins from the specified liquidity pool,  to stop receiving rewards
Example:
$ %s tx %s soft-unlock 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			softUnlockCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSoftUnlock(
				clientCtx.GetFromAddress(),
				poolID,
				softUnlockCoin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
