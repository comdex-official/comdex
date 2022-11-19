package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/petrichormoney/petri/x/liquidity/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "liquidity",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "liquidity"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryGenericPairsCmd(),
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
		NewQueryFarmerCmd(),
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

// NewQueryGenericPairsCmd implements the pgeneric-params query command.
func NewQueryGenericPairsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generic-params [app-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current liquidity parameters information for app",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as liquidity parameters for a given app.
Example:
$ %s query %s generic-params 1
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

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GenericParams(
				cmd.Context(),
				&types.QueryGenericParamsRequest{
					AppId: appID,
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

// NewQueryPairsCmd implements the pairs query command.
func NewQueryPairsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pairs [app-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all pairs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all existing pairs on a network.
Example:
$ %s query %s pairs 1
$ %s query %s pairs 1 --denoms=uatom
$ %s query %s pairs 1 --denoms=uatom,stake
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			denoms, _ := cmd.Flags().GetStringSlice(FlagDenoms)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Pairs(cmd.Context(), &types.QueryPairsRequest{
				Denoms:     denoms,
				AppId:      appID,
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
		Use:   "pair [app-id] [pair-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of the pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the pair.
Example:
$ %s query %s pair 1 1
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Pair(cmd.Context(), &types.QueryPairRequest{
				PairId: pairID,
				AppId:  appID,
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
		Use:   "pools [app-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query for all liquidity pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all existing liquidity pools on a network.
Example:
$ %s query %s pools 1
$ %s query %s pools 1 --pair-id=1
$ %s query %s pools 1 --disabled=true
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
				AppId:      appID,
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
		Use:   "pool [app-id] [pool-id]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of the liquidity pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the liquidity pool
Example:
$ %s query %s pool 1 1
$ %s query %s pool 1 --pool-coin-denom=pool1
$ %s query %s pool 1 --reserve-address=petri...
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

			var poolID *uint64
			if len(args) > 1 {
				id, err := strconv.ParseUint(args[1], 10, 64)
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
					AppId:  appID,
					PoolId: *poolID,
				})
			case poolCoinDenom != "":
				res, err = queryClient.PoolByPoolCoinDenom(
					cmd.Context(),
					&types.QueryPoolByPoolCoinDenomRequest{
						AppId:         appID,
						PoolCoinDenom: poolCoinDenom,
					})
			case reserveAddr != "":
				res, err = queryClient.PoolByReserveAddress(
					cmd.Context(),
					&types.QueryPoolByReserveAddressRequest{
						AppId:          appID,
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
		Use:   "deposit-requests [app-id] [pool-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query for all deposit requests in the pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all deposit requests in the pool.
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositRequests(
				cmd.Context(),
				&types.QueryDepositRequestsRequest{
					AppId:      appID,
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
		Use:   "deposit-request [app-id] [pool-id] [id]",
		Args:  cobra.ExactArgs(3),
		Short: "Query details of the specific deposit request",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific deposit request.
Example:
$ %s query %s deposit-request 1 1 1
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

			id, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DepositRequest(
				cmd.Context(),
				&types.QueryDepositRequestRequest{
					AppId:  appID,
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
		Use:   "withdraw-requests [app-id] [pool-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query for all withdraw requests in the pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all withdraw requests in the pool.
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.WithdrawRequests(
				cmd.Context(),
				&types.QueryWithdrawRequestsRequest{
					AppId:      appID,
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
		Use:   "withdraw-request [app-id] [pool-id] [id]",
		Args:  cobra.ExactArgs(3),
		Short: "Query details of the specific withdraw request",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific withdraw request.
Example:
$ %s query %s withdraw-request 1 1 1
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

			id, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.WithdrawRequest(
				cmd.Context(),
				&types.QueryWithdrawRequestRequest{
					AppId:  appID,
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
		Use:   "orders [app-id] [orderer]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query for all orders in the pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all orders in the pair.
Example:
$ %s query %s orders 1 petri...
$ %s query %s orders 1 --pair-id=1 petri...
$ %s query %s orders 1 --pair-id=1
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse app id: %w", err)
			}

			var orderer *string
			if len(args) > 1 {
				orderer = &args[1]
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
					AppId:      appID,
					PairId:     pairID,
					Pagination: pageReq,
				})
			} else {
				res, err = queryClient.OrdersByOrderer(
					cmd.Context(),
					&types.QueryOrdersByOrdererRequest{
						Orderer:    *orderer,
						AppId:      appID,
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
		Use:   "order [app-id] [pair-id] [id]",
		Args:  cobra.ExactArgs(3),
		Short: "Query details of the specific order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the specific order.
Example:
$ %s query %s order 1 1 1
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

			id, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Order(
				cmd.Context(),
				&types.QueryOrderRequest{
					AppId:  appID,
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

// NewQueryFarmerCmd implements the farmer query command.
func NewQueryFarmerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmer [app-id] [pool-id] [farmer]",
		Args:  cobra.ExactArgs(3),
		Short: "Query farmer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query farming status of the farmer.
Example:
$ %s query %s farmer 1 1 petri...
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
				return fmt.Errorf("parse pool id: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Farmer(
				cmd.Context(),
				&types.QueryFarmerRequest{
					AppId:  appID,
					PoolId: poolID,
					Farmer: args[2],
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

// NewQueryDeserializePoolCoinCmd implements the deserialize query command.
func NewQueryDeserializePoolCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deserialize [app-id] [pool-id] [pool-coin-amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Deserialize pool coins into the pool assets",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deserialize pool coins into pool assets.
Example:
$ %s query %s deserialize 1 1 123400000
> {coins : [1000000upetri, 4000000ucmst]}
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
			poolCoinAmount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DeserializePoolCoin(
				cmd.Context(),
				&types.QueryDeserializePoolCoinRequest{
					AppId:          appID,
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
		Use:   "pool-incentives [app-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query pool incentives",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query Pool Incentives
Example:
$ %s query %s pool-incentives 1
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PoolIncentives(
				cmd.Context(),
				&types.QueryPoolsIncentivesRequest{
					AppId: appID,
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

// NewQueryFarmedPoolCoinCmd implements the farmed-coin query command.
func NewQueryFarmedPoolCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmed-coin [app-id] [pool-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query total coins being farmed",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total coins being farmed.
Example:
$ %s query %s farmed-coin 1 1
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FarmedPoolCoin(
				cmd.Context(),
				&types.QueryFarmedPoolCoinRequest{
					AppId:  appID,
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
