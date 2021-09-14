package cli

import (
	"context"
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func queryMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market [symbol]",
		Short: "Query a market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryMarket(
				context.Background(),
				&types.QueryMarketRequest{
					Symbol: args[0],
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

func queryMarkets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markets",
		Short: "Query markets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryMarkets(
				context.Background(),
				&types.QueryMarketsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "markets")

	return cmd
}