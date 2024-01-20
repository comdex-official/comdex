package cli

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "liquidity",
		Short:                      fmt.Sprintf("%s transactions subcommands", "liquidity"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreatePairCmd(),
		NewCreatePoolCmd(),
		NewCreateRangedPoolCmd(),
		NewDepositCmd(),
		NewWithdrawCmd(),
		NewLimitOrderCmd(),
		NewMarketOrderCmd(),
		NewMMOrderCmd(),
		NewCancelOrderCmd(),
		NewCancelAllOrdersCmd(),
		NewCancelMMOrderCmd(),
		NewFarmCmd(),
		NewUnfarmCmd(),
		NewDepositAndFarmCmd(),
		NewUnfarmAndWithdrawCmd(),
	)

	return cmd
}

func NewCreatePairCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pair [app-id] [base-coin-denom] [quote-coin-denom]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a pair(market) for trading",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a pair(market) for trading.
Example:
$ %s tx %s create-pair 1 uatom stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}
			baseCoinDenom := args[1]
			quoteCoinDenom := args[2]

			msg := types.NewMsgCreatePair(appID, clientCtx.GetFromAddress(), baseCoinDenom, quoteCoinDenom)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCreatePoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [app-id] [pair-id] [deposit-coins]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a liquidity pool with coins.
Example:
$ %s tx %s create-pool 1 1 1000000000uatom,50000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid deposit coints: %w", err)
			}

			msg := types.NewMsgCreatePool(appID, clientCtx.GetFromAddress(), pairID, depositCoins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCreateRangedPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-ranged-pool [app-id] [pair-id] [deposit-coins] [min-price] [max-price] [initial-price]",
		Args:  cobra.ExactArgs(6),
		Short: "Create a ranged liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a ranged liquidity pool with coins.

Example:
$ %s tx %s create-ranged-pool 1 1 1000000000uatom,10000000000stake 0.001 100 1.0 --from mykey
$ %s tx %s create-ranged-pool 1 1 1000000000uatom,10000000000stake 0.9 10000 1.0 --from mykey
$ %s tx %s create-ranged-pool 1 1 1000000000uatom,10000000000stake 1.3 2.5 1.5 --from mykey
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid deposit coins: %w", err)
			}

			minPrice, err := sdkmath.LegacyNewDecFromStr(args[3])
			if err != nil {
				return fmt.Errorf("invalid min price: %w", err)
			}

			maxPrice, err := sdkmath.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid max price: %w", err)
			}

			initialPrice, err := sdkmath.LegacyNewDecFromStr(args[5])
			if err != nil {
				return fmt.Errorf("invalid initial price: %w", err)
			}

			msg := types.NewMsgCreateRangedPool(
				appID,
				clientCtx.GetFromAddress(),
				pairID,
				depositCoins,
				minPrice,
				maxPrice,
				initialPrice,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewDepositCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [app-id] [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit coins to a liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to a liquidity pool.
Example:
$ %s tx %s deposit 1 1 1000000000uatom,50000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid deposit coins: %w", err)
			}

			msg := types.NewMsgDeposit(appID, clientCtx.GetFromAddress(), poolID, depositCoins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [app-id] [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(3),
		Short: "Withdraw coins from the specified liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw coins from the specified liquidity pool.
Example:
$ %s tx %s withdraw 1 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			poolCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(
				appID,
				clientCtx.GetFromAddress(),
				poolID,
				poolCoin,
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

func NewLimitOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit-order [app-id] [pair-id] [direction] [offer-coin] [demand-coin-denom] [price] [amount]",
		Args:  cobra.ExactArgs(7),
		Short: "Make a limit order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a limit order.
Example:
$ %s tx %s limit-order 1 1 buy 5000stake uatom 0.5 10000 --from mykey
$ %s tx %s limit-order 1 1 b 5000stake uatom 0.5 10000 --from mykey
$ %s tx %s limit-order 1 1 sell 10000uatom stake 2.0 10000 --order-lifespan=10m --from mykey
$ %s tx %s limit-order 1 1 s 10000uatom stake 2.0 10000 --order-lifespan=10m --from mykey

[app-id]: application id on which transaction to be made
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
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			dir, err := parseOrderDirection(args[2])
			if err != nil {
				return fmt.Errorf("parse order direction: %w", err)
			}

			offerCoin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid offer coin: %w", err)
			}

			demandCoinDenom := args[4]
			if err := sdk.ValidateDenom(demandCoinDenom); err != nil {
				return fmt.Errorf("invalid demand coin denom: %w", err)
			}

			price, err := sdkmath.LegacyNewDecFromStr(args[5])
			if err != nil {
				return fmt.Errorf("invalid price: %w", err)
			}

			amt, ok := sdkmath.NewIntFromString(args[6])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[6])
			}

			orderLifespan, _ := cmd.Flags().GetDuration(FlagOrderLifespan)

			msg := types.NewMsgLimitOrder(
				appID,
				clientCtx.GetFromAddress(),
				pairID,
				dir,
				offerCoin,
				demandCoinDenom,
				price,
				amt,
				orderLifespan,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetOrder())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewMarketOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market-order [app-id] [pair-id] [direction] [offer-coin] [demand-coin-denom] [amount]",
		Args:  cobra.ExactArgs(6),
		Short: "Make a market order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a market order.
Example:
$ %s tx %s market-order 1 1 buy 5000stake uatom 10000 --from mykey
$ %s tx %s market-order 1 1 b 5000stake uatom 10000 --from mykey
$ %s tx %s market-order 1 1 sell 10000uatom stake 10000 --order-lifespan=10m --from mykey
$ %s tx %s market-order 1 1 s 10000uatom stake 10000 --order-lifespan=10m --from mykey

[app-id]: application id on which transaction to be made
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			dir, err := parseOrderDirection(args[2])
			if err != nil {
				return fmt.Errorf("parse order direction: %w", err)
			}

			offerCoin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid offer coin: %w", err)
			}

			demandCoinDenom := args[4]
			if err := sdk.ValidateDenom(demandCoinDenom); err != nil {
				return fmt.Errorf("invalid demand coin denom: %w", err)
			}

			amt, ok := sdkmath.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[5])
			}

			orderLifespan, _ := cmd.Flags().GetDuration(FlagOrderLifespan)

			msg := types.NewMsgMarketOrder(
				appID,
				clientCtx.GetFromAddress(),
				pairID,
				dir,
				offerCoin,
				demandCoinDenom,
				amt,
				orderLifespan,
			)

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetOrder())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewMMOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mm-order [app-id] [pair-id] [max-sell-price] [min-sell-price] [sell-amount] [max-buy-price] [min-buy-price] [buy-amount]",
		Args:  cobra.ExactArgs(8),
		Short: "Make a market making order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a market making order.
A market making order is a set of limit orders for each buy/sell side.
You can leave one side(but not both) empty by passing 0 as its arguments.

Example:
$ %s tx %s mm-order 1 1 102 101 10000 100 99 10000 --from mykey
$ %s tx %s mm-order 1 1 0 0 0 100 99 10000 --from mykey
$ %s tx %s mm-order 1 1 102 101 10000 0 0 0 --from mykey

[app-id]: application id on which transaction to be made
[pair-id]: pair id to make order
[max-sell-price]: maximum price of sell orders
[min-sell-price]]: minimum price of sell orders
[sell-amount]: total amount of sell orders
[max-buy-price]: maximum price of buy orders
[min-buy-price]: minimum price of buy orders
[buy-amount]: the total amount of buy orders
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			maxSellPrice, err := sdkmath.LegacyNewDecFromStr(args[2])
			if err != nil {
				return fmt.Errorf("invalid max sell price: %w", err)
			}

			minSellPrice, err := sdkmath.LegacyNewDecFromStr(args[3])
			if err != nil {
				return fmt.Errorf("invalid min sell price: %w", err)
			}

			sellAmt, ok := sdkmath.NewIntFromString(args[4])
			if !ok {
				return fmt.Errorf("invalid sell amount: %s", args[4])
			}

			maxBuyPrice, err := sdkmath.LegacyNewDecFromStr(args[5])
			if err != nil {
				return fmt.Errorf("invalid max buy price: %w", err)
			}

			minBuyPrice, err := sdkmath.LegacyNewDecFromStr(args[6])
			if err != nil {
				return fmt.Errorf("invalid min buy price: %w", err)
			}

			buyAmt, ok := sdkmath.NewIntFromString(args[7])
			if !ok {
				return fmt.Errorf("invalid buy amount: %s", args[7])
			}

			orderLifespan, _ := cmd.Flags().GetDuration(FlagOrderLifespan)

			msg := types.NewMsgMMOrder(
				appID,
				clientCtx.GetFromAddress(),
				pairID,
				maxSellPrice, minSellPrice, sellAmt,
				maxBuyPrice, minBuyPrice, buyAmt,
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
		Use:   "cancel-order [app-id] [pair-id] [order-id]",
		Args:  cobra.ExactArgs(3),
		Short: "Cancel an order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel an order.
Example:
$ %s tx %s cancel-order 1 1 1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrder(
				appID,
				clientCtx.GetFromAddress(),
				pairID,
				orderID,
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

func NewCancelAllOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-all-orders [app-id] [pair-ids]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Cancel all orders",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel all orders.
Example:
$ %s tx %s cancel-all-orders 1 --from mykey
$ %s tx %s cancel-all-orders 1 1,3 --from mykey
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			var pairIDs []uint64
			for _, pairIDStr := range strings.Split(args[1], ",") {
				pairID, err := strconv.ParseUint(pairIDStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
				pairIDs = append(pairIDs, pairID)
			}

			msg := types.NewMsgCancelAllOrders(appID, clientCtx.GetFromAddress(), pairIDs)

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelMMOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-mm-order [app-id] [pair-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel the mm order in a pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel the mm order in a pair.
This will cancel all limit orders in the pair made by the mm order.

Example:
$ %s tx %s cancel-mm-order 1 1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelMMOrder(appID, clientCtx.GetFromAddress(), pairID)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewFarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farm [app-id] [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(3),
		Short: "farm pool coin",
		Long: strings.TrimSpace(
			fmt.Sprintf(`farm pool coins to be eligible for incentivizations 
Example:
$ %s tx %s farm 1 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			farmedPoolCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgFarm(
				appID,
				poolID,
				clientCtx.GetFromAddress(),
				farmedPoolCoin,
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

func NewUnfarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfarm [app-id] [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(3),
		Short: "unfarm pool coin",
		Long: strings.TrimSpace(
			fmt.Sprintf(`unfarm pool coin
Example:
$ %s tx %s unfarm 1 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			unfarmingPoolCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnfarm(
				appID,
				poolID,
				clientCtx.GetFromAddress(),
				unfarmingPoolCoin,
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

func NewCmdUpdateGenericParamsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidity-param-change [app-id] [keys] [values]",
		Args:  cobra.ExactArgs(3),
		Short: "Update params for app in liquidity module",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			keys, err := ParseStringSliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			values, err := ParseStringSliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUpdateGenericParamsProposal(
				title,
				description,
				appID,
				keys,
				values,
			)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func NewCmdCreateNewLiquidityPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-liquidity-pair [app-id] [base-coin-denom] [quote-coin-denom]",
		Args:  cobra.ExactArgs(3),
		Short: "Create new liquidity pair in given app",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			baseCoinDenom := args[1]
			quoteCoinDenom := args[2]

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewCreateLiquidityPairProposal(
				title,
				description,
				from,
				appID,
				baseCoinDenom,
				quoteCoinDenom,
			)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func NewDepositAndFarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-and-farm [app-id] [pool-id] [deposit-coins]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit coins to a liquidity pool  and farm pool tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit coins to a liquidity pool and farm pool tokens.
Example:
$ %s tx %s deposit-and-farm 1 1 1000000000uatom,50000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool id: %w", err)
			}

			depositCoins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid deposit coins: %w", err)
			}

			msg := types.NewMsgDepositAndFarm(appID, clientCtx.GetFromAddress(), poolID, depositCoins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUnfarmAndWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfarm-and-withdraw [app-id] [pool-id] [pool-coin]",
		Args:  cobra.ExactArgs(3),
		Short: "unfarm pool coin and withdraw liquidity from pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`unfarm pool coin and withdraw liquidity from pool
Example:
$ %s tx %s unfarm-and-withdraw 1 1 10000pool1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			unfarmingPoolCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnfarmAndWithdraw(
				appID,
				poolID,
				clientCtx.GetFromAddress(),
				unfarmingPoolCoin,
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
