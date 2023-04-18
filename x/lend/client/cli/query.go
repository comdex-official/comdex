package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/lend/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group lend queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "lend",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "lend"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryLend(),
		queryLends(),
		QueryAllLendsByOwner(),
		QueryAllLendsByOwnerAndPoolID(),
		queryPair(),
		queryPairs(),
		queryPool(),
		queryPools(),
		queryAssetToPairMapping(),
		queryAssetToPairMappings(),
		queryBorrow(),
		queryBorrows(),
		QueryAllBorrowsByOwner(),
		QueryAllBorrowsByOwnerAndPoolID(),
		QueryAssetRatesParam(),
		QueryAssetRatesParams(),
		QueryPoolAssetLBMapping(),
		queryReserveBuybackAssetData(),
		queryAuctionParams(),
		QueryModuleBalance(),
		QueryFundModuleBalance(),
		QueryFundReserveBalance(),
		QueryAllReserveStats(),
		QueryFundModBalByAssetPool(),
		queryLendInterest(),
		queryBorrowInterest(),
		queryUserLendRewards(),
		queryUserBorrowInterest(),
	)

	return cmd
}

func queryLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend [id]",
		Short: "Query a lend position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryLend(
				context.Background(),
				&types.QueryLendRequest{
					Id: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryLends() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lends",
		Short: "Query lends",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLends(
				context.Background(),
				&types.QueryLendsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lends")

	return cmd
}

func QueryAllLendsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lends-by-owner [owner]",
		Short: "lends list for a owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAllLendByOwner(cmd.Context(), &types.QueryAllLendByOwnerRequest{
				Owner:      args[0],
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lends-by-owner")

	return cmd
}

func QueryAllLendsByOwnerAndPoolID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lends-by-owner-pool [owner] [pool-id]",
		Short: "lends list for a owner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAllLendByOwnerAndPool(cmd.Context(), &types.QueryAllLendByOwnerAndPoolRequest{
				Owner:      args[0],
				PoolId:     poolID,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lends-by-owner-pool")
	return cmd
}

func queryPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair [id]",
		Short: "Query a pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryPair(
				context.Background(),
				&types.QueryPairRequest{
					Id: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pairs",
		Short: "Query pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryPairs(
				context.Background(),
				&types.QueryPairsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pairs")

	return cmd
}

func queryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [id]",
		Short: "Query a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryPool(
				context.Background(),
				&types.QueryPoolRequest{
					Id: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Query pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryPools(
				context.Background(),
				&types.QueryPoolsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pools")

	return cmd
}

func queryAssetToPairMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-pair-mapping [asset-id] [pool-id]",
		Short: "Query Asset To Pair Mapping",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssetToPairMapping(
				context.Background(),
				&types.QueryAssetToPairMappingRequest{
					AssetId: assetID,
					PoolId:  poolID,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryAssetToPairMappings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-pair-mappings",
		Short: "Query Asset To Pair Mappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssetToPairMappings(
				context.Background(),
				&types.QueryAssetToPairMappingsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "asset-pair-mappings")

	return cmd
}

func queryBorrow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow [id]",
		Short: "Query a borrow position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryBorrow(
				context.Background(),
				&types.QueryBorrowRequest{
					Id: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryBorrows() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrows",
		Short: "Query borrows",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryBorrows(
				context.Background(),
				&types.QueryBorrowsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "borrows")

	return cmd
}

func QueryAllBorrowsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrows-by-owner [owner]",
		Short: "borrows list for a owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAllBorrowByOwner(cmd.Context(), &types.QueryAllBorrowByOwnerRequest{
				Owner:      args[0],
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "borrows-by-owner")

	return cmd
}

func QueryAllBorrowsByOwnerAndPoolID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrows-by-owner-pool [owner] [pool-id]",
		Short: "borrows list for a owner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueryAllBorrowByOwnerAndPool(cmd.Context(), &types.QueryAllBorrowByOwnerAndPoolRequest{
				Owner:      args[0],
				PoolId:     poolID,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "borrows-by-owner-pool")

	return cmd
}

func QueryAssetRatesParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-rates-stat [asset-id]",
		Short: "Query asset rates stat by asset-id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssetRatesParam(
				context.Background(),
				&types.QueryAssetRatesParamRequest{
					Id: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAssetRatesParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-rates-stats",
		Short: "Query asset rates stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssetRatesParams(
				context.Background(),
				&types.QueryAssetRatesParamsRequest{
					Pagination: pagination,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "asset-rates-stats")

	return cmd
}

func QueryPoolAssetLBMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-stats [pool-id] [asset-id]",
		Short: "Query asset stats for an asset-id and pool-id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueryPoolAssetLBMapping(cmd.Context(), &types.QueryPoolAssetLBMappingRequest{
				AssetId: assetID,
				PoolId:  poolID,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryReserveBuybackAssetData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-back-deposit-stats [id]",
		Short: "Query reserve deposit stats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryReserveBuybackAssetData(
				context.Background(),
				&types.QueryReserveBuybackAssetDataRequest{
					AssetId: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "buy-back-deposit-stats")

	return cmd
}

func queryAuctionParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction-params [id]",
		Short: "Query auction-params",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAuctionParams(
				context.Background(),
				&types.QueryAuctionParamRequest{
					AppId: id,
				},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryModuleBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-balance [pool-id]",
		Short: "queries module balance of a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueryModuleBalance(cmd.Context(), &types.QueryModuleBalanceRequest{
				PoolId: poolID,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryFundModuleBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-module-balance ",
		Short: "queries fund module balance history",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryFundModBal(cmd.Context(), &types.QueryFundModBalRequest{})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryFundReserveBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-reserve-balance ",
		Short: "queries fund reserve balance history",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryFundReserveBal(cmd.Context(), &types.QueryFundReserveBalRequest{})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryAllReserveStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-reserve-stats [id]",
		Short: "queries all reserve stats of an asset id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			assetID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueryAllReserveStats(cmd.Context(), &types.QueryAllReserveStatsRequest{AssetId: assetID})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryFundModBalByAssetPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund_mod_bal_by_asset_pool [asset_id] [pool_id]",
		Short: "queries all reserve stats of an asset id and pool id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			assetID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.QueryFundModBalByAssetPool(cmd.Context(), &types.QueryFundModBalByAssetPoolRequest{AssetId: assetID, PoolId: poolID})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryLendInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend_interest",
		Short: "Query all lend interest",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLendInterest(
				context.Background(),
				&types.QueryLendInterestRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lend_interest")

	return cmd
}

func queryBorrowInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow_interest",
		Short: "Query all borrow interest",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryBorrowInterest(
				context.Background(),
				&types.QueryBorrowInterestRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lend_interest")

	return cmd
}

func queryUserLendRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user_lend_rewards [id]",
		Short: "Query user lend rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			ID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.QueryUserLendRewards(
				context.Background(),
				&types.QueryUserLendRewardsRequest{Id: ID},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "user_lend_rewards")

	return cmd
}

func queryUserBorrowInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user_borrow_interest [id]",
		Short: "Query user borrow interest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			ID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.QueryUserBorrowInterest(
				context.Background(),
				&types.QueryUserBorrowInterestRequest{Id: ID},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "user_borrow_interest")

	return cmd
}
