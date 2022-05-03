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

// GetTxCmd returns the transaction commands for the module
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
		NewBondPoolTokens(),
		NewUnbondPoolTokens(),
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

			pairId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid deposit coints: %w", err)
			}

			msg := types.NewMsgCreatePool(clientCtx.GetFromAddress(), pairId, depositCoins)

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

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid deposit coins: %w", err)
			}

			msg := types.NewMsgDeposit(clientCtx.GetFromAddress(), poolId, depositCoins)

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

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			poolCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(
				clientCtx.GetFromAddress(),
				poolId,
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

			pairId, err := strconv.ParseUint(args[0], 10, 64)
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
				pairId,
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

			pairId, err := strconv.ParseUint(args[0], 10, 64)
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
				pairId,
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

			pairId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			orderId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrder(
				clientCtx.GetFromAddress(),
				pairId,
				orderId,
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

			var pairIds []uint64
			for _, pairIdStr := range strings.Split(args[0], ",") {
				pairId, err := strconv.ParseUint(pairIdStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
				pairIds = append(pairIds, pairId)
			}

			msg := types.NewMsgCancelAllOrders(clientCtx.GetFromAddress(), pairIds)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewBondPoolTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Bond pool coin from the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bond pool coin from the specified liquidity pool.


Example:
$ %s tx %s bond 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from mykey

This example request bonds 10000 pool coin from the specified liquidity pool.
The appropriate pool coin must be requested from the specified pool.

[pool-id]: The pool id of the liquidity pool
[pool-coin]: The amount of pool coin to bond from the liquidity pool
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			userAddress := clientCtx.GetFromAddress()

			// Get pool type index
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, input a valid unsigned 32-bit integer for pool-id", args[0])
			}

			// Get pool coin of the target pool
			poolCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			err = poolCoin.Validate()
			if err != nil {
				return err
			}

			msg := types.NewMsgBondPoolTokens(userAddress, poolID, poolCoin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
func NewUnbondPoolTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(2),
		Short: "Unbond pool coin from the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond pool coin from the specified liquidity pool.


Example:
$ %s tx %s unbond 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from mykey

This example request unbonds 10000 pool coin from the specified liquidity pool.
The appropriate pool coin must be requested from the specified pool.

[pool-id]: The pool id of the liquidity pool
[pool-coin]: The amount of pool coin to unbond from the liquidity pool
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			userAddress := clientCtx.GetFromAddress()

			// Get pool type index
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("pool-id %s not a valid uint, input a valid unsigned 32-bit integer for pool-id", args[0])
			}

			// Get pool coin of the target pool
			poolCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			err = poolCoin.Validate()
			if err != nil {
				return err
			}

			msg := types.NewMsgUnbondPoolTokens(userAddress, poolID, poolCoin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
