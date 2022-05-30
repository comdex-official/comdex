package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryPoolsCmd(),
		NewQueryPoolCmd(),
		NewQueryPairsCmd(),
		NewQueryPairCmd(),
		NewQueryDepositRequestsCmd(),
		NewQueryDepositRequestCmd(),
		NewQueryWithdrawRequestsCmd(),
		NewQueryWithdrawRequestCmd(),
		NewQueryOrdersCmd(),
		NewQueryOrderCmd(),
		NewQuerySoftLockCmd(),
		NewQueryDeserializePoolCoinCmd(),
		NewQueryPoolIncentivesCmd(),
		NewQueryFarmedPoolCoinCmd(),
	)

	return cmd
}

// NewQueryParamsCmd implements the params query command.
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current liquidity parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as liquidity parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryPairsCmd implements the pairs query command.
func NewQueryPairsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pairs",
		Args:  cobra.NoArgs,
		Short: "Query for all pairs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all existing pairs on a network.
Example:
$ %s query %s pairs
$ %s query %s pairs --denoms=uatom
$ %s query %s pairs --denoms=uatom,stake
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			denoms, _ := cmd.Flags().GetStringSlice(FlagDenoms)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Pairs(cmd.Context(), &types.QueryPairsRequest{
				Denoms:     denoms,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(flagSetPairs())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryPairCmd implements the pair query command.
func NewQueryPairCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair [pair-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of the pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the pair.
Example:
$ %s query %s pair 1
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Pair(cmd.Context(), &types.QueryPairRequest{
				PairId: pairID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryPoolsCmd implements the pools query command.
func NewQueryPoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Args:  cobra.NoArgs,
		Short: "Query for all liquidity pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all existing liquidity pools on a network.
Example:
$ %s query %s pools
$ %s query %s pools --pair-id=1
$ %s query %s pools --disabled=true
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var pairID uint64

			pairIDStr, _ := cmd.Flags().GetString(FlagPairID)
			if pairIDStr != "" {
				var err error
				pairID, err = strconv.ParseUint(pairIDStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id flag: %w", err)
				}
			}
			disabledStr, _ := cmd.Flags().GetString(FlagDisabled)
			if disabledStr != "" {
				if _, err := strconv.ParseBool(disabledStr); err != nil {
					return fmt.Errorf("parse disabled flag: %w", err)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Pools(cmd.Context(), &types.QueryPoolsRequest{
				PairId:     pairID,
				Disabled:   disabledStr,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(flagSetPools())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryPoolCmd implements the pool query command.
func NewQueryPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Query details of the liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the liquidity pool
Example:
$ %s query %s pool 1
$ %s query %s pool --pool-coin-denom=pool1
$ %s query %s pool --reserve-address=cre1...
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

			var poolID *uint64
			if len(args) > 0 {
				id, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("parse pool id: %w", err)
				}
				poolID = &id
			}
			poolCoinDenom, _ := cmd.Flags().GetString(FlagPoolCoinDenom)
			reserveAddr, _ := cmd.Flags().GetString(FlagReserveAddress)

			if !excConditions(poolID != nil, poolCoinDenom != "", reserveAddr != "") {
				return fmt.Errorf("invalid request")
			}

			queryClient := types.NewQueryClient(clientCtx)
			var res *types.QueryPoolResponse
			switch {
			case poolID != nil:
				res, err = queryClient.Pool(cmd.Context(), &types.QueryPoolRequest{
					PoolId: *poolID,
				})
			case poolCoinDenom != "":
				res, err = queryClient.PoolByPoolCoinDenom(
					cmd.Context(),
					&types.QueryPoolByPoolCoinDenomRequest{
						PoolCoinDenom: poolCoinDenom,
					})
			case reserveAddr != "":
				res, err = queryClient.PoolByReserveAddress(
					cmd.Context(),
					&types.QueryPoolByReserveAddressRequest{
						ReserveAddress: reserveAddr,
					})
			}
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(flagSetPool())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryDepositRequestsCmd implements the deposit requests query command.
func NewQueryDepositRequestsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-requests [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all deposit requests in the pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all deposit requests in the pool.
Example:
$ %s query %s deposit-requests 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositRequests(
				cmd.Context(),
				&types.QueryDepositRequestsRequest{
					PoolId:     poolID,
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryDepositRequestCmd implements the deposit request query command.
func NewQueryDepositRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-request [pool-id] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of the specific deposit request",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific deposit request.
Example:
$ %s query %s deposit-requests 1 1
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

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositRequest(
				cmd.Context(),
				&types.QueryDepositRequestRequest{
					PoolId: poolID,
					Id:     id,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryWithdrawRequestsCmd implements the withdraw requests query command.
func NewQueryWithdrawRequestsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-requests [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all withdraw requests in the pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all withdraw requests in the pool.
Example:
$ %s query %s withdraw-requests 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.WithdrawRequests(
				cmd.Context(),
				&types.QueryWithdrawRequestsRequest{
					PoolId:     poolID,
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryWithdrawRequestCmd implements the withdraw request query command.
func NewQueryWithdrawRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-request [pool-id] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of the specific withdraw request",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific withdraw request.
Example:
$ %s query %s withdraw-requests 1 1
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

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.WithdrawRequest(
				cmd.Context(),
				&types.QueryWithdrawRequestRequest{
					PoolId: poolID,
					Id:     id,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryOrdersCmd implements the orders query command.
func NewQueryOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orders [orderer]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Query for all orders in the pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all orders in the pair.
Example:
$ %s query %s orders cre1...
$ %s query %s orders --pair-id=1 cre1...
$ %s query %s orders --pair-id=1
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var orderer *string
			if len(args) > 0 {
				orderer = &args[0]
			}

			var pairID uint64
			pairIDStr, _ := cmd.Flags().GetString(FlagPairID)
			if pairIDStr != "" {
				pairID, err = strconv.ParseUint(pairIDStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
			}
			if orderer == nil && pairID == 0 {
				return fmt.Errorf("either orderer or pair-id must be specified")
			}

			queryClient := types.NewQueryClient(clientCtx)

			var res *types.QueryOrdersResponse
			if orderer == nil {
				res, err = queryClient.Orders(cmd.Context(), &types.QueryOrdersRequest{
					PairId:     pairID,
					Pagination: pageReq,
				})
			} else {
				res, err = queryClient.OrdersByOrderer(
					cmd.Context(),
					&types.QueryOrdersByOrdererRequest{
						Orderer:    *orderer,
						PairId:     pairID,
						Pagination: pageReq,
					})
			}
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(flagSetOrders())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryOrderCmd implements the order query command.
func NewQueryOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order [pair-id] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of the specific order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific order.
Example:
$ %s query %s order 1 1
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

			id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Order(
				cmd.Context(),
				&types.QueryOrderRequest{
					PairId: pairID,
					Id:     id,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQuerySoftLockCmd implements the soft lock query command.
func NewQuerySoftLockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soft-lock [pool-id] [depositor]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of the soft-lock",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the soft-lock in specific pool for adddress.
Example:
$ %s query %s soft-lock 1 comdex1ed6zea6ppj29vkzk8f867rsauu65lq2p75jc3u
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.SoftLock(
				cmd.Context(),
				&types.QuerySoftLockRequest{
					PoolId:    poolID,
					Depositor: args[1],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQuerySoftLockCmd implements the soft lock query command.
func NewQueryDeserializePoolCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deserialize [pool-id] [pool-coin-amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Deserialize pool coins into the pool assets",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deserialize pool coins into pool assets.
Example:
$ %s query %s deserialize 1 123400000
> {coins : [1000000ucmdx, 4000000ucmst]}
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
			poolCoinAmount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DeserializePoolCoin(
				cmd.Context(),
				&types.QueryDeserializePoolCoinRequest{
					PoolId:         poolID,
					PoolCoinAmount: poolCoinAmount,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryPoolIncentivesCmd implements the pool-incentives query command.
func NewQueryPoolIncentivesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-incentives",
		Args:  cobra.NoArgs,
		Short: "Query pool incentives",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query Pool Incentives
Example:
$ %s query %s pool-incentives
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PoolIncentives(
				cmd.Context(),
				&types.QueryPoolsIncentivesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryFarmedPoolCoinCmd implements the farmed-coin query command.
func NewQueryFarmedPoolCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmed-coin [pool-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query total coins being farmed (in soft-lock)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total coins being farmed (in soft-lock).
Example:
$ %s query %s farmed-coin 1
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FarmedPoolCoin(
				cmd.Context(),
				&types.QueryFarmedPoolCoinRequest{
					PoolId: poolID,
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
