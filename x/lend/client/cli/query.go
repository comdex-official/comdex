package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group lend queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
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
		queryAssetRatesStat(),
		queryAssetRatesStats(),
		QueryAssetStats(),
		QueryModuleBalance(),
		queryDepositStats(),
		queryUserDepositStats(),
		queryReserveDepositStats(),
		queryBuyBackDepositStats(),
		queryBorrowStats(),
		queryAuctionParams(),
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
	flags.AddPaginationFlagsToCmd(cmd, "lends")

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
	flags.AddPaginationFlagsToCmd(cmd, "lends")

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
	return cmd
}

func queryAssetRatesStat() *cobra.Command {
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

			res, err := queryClient.QueryAssetRatesStat(
				context.Background(),
				&types.QueryAssetRatesStatRequest{
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

func queryAssetRatesStats() *cobra.Command {
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

			res, err := queryClient.QueryAssetRatesStats(
				context.Background(),
				&types.QueryAssetRatesStatsRequest{
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

func QueryAssetStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset-stats [asset-id] [pool-id]",
		Short: "Query asset stats for an asset-id and pool-id",
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

			res, err := queryClient.QueryAssetStats(cmd.Context(), &types.QueryAssetStatsRequest{
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

func QueryModuleBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-balance [pool-id]",
		Short: "borrows list for a owner",
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

func queryDepositStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-stats",
		Short: "Query deposit stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryDepositStats(
				context.Background(),
				&types.QueryDepositStatsRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "deposit-stats")

	return cmd
}

func queryUserDepositStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-deposit-stats",
		Short: "Query user deposit stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryUserDepositStats(
				context.Background(),
				&types.QueryUserDepositStatsRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "user-deposit-stats")

	return cmd
}

func queryReserveDepositStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserve-deposit-stats",
		Short: "Query reserve deposit stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryReserveDepositStats(
				context.Background(),
				&types.QueryReserveDepositStatsRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "reserve-deposit-stats")

	return cmd
}

func queryBuyBackDepositStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-back-deposit-stats",
		Short: "Query reserve deposit stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryBuyBackDepositStats(
				context.Background(),
				&types.QueryBuyBackDepositStatsRequest{},
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

func queryBorrowStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow-stats",
		Short: "Query borrow stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryBorrowStats(
				context.Background(),
				&types.QueryBorrowStatsRequest{},
			)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "borrow-stats")

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
