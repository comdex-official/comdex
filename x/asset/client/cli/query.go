package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/asset/types"
)

func queryPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair [id]",
		Short: "Query an asset pair",
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

			var (
				qc = types.NewQueryServiceClient(ctx)
			)

			res, err := qc.QueryPair(
				context.Background(),
				&types.QueryPairRequest{
					ID: id,
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
		Short: "Query asset pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				qc = types.NewQueryServiceClient(ctx)
			)

			res, err := qc.QueryPairs(
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
	flags.AddPaginationFlagsToCmd(cmd, "asset pairs")

	return cmd
}
