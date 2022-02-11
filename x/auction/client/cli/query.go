package cli

import (
	"context"
	"strconv"

	"github.com/comdex-official/comdex/x/auction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func queryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction [id]",
		Short: "Query an auction",
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
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryAuction(
				context.Background(),
				&types.QueryAuctionRequest{
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

func queryAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auctions",
		Short: "Query auctions",
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
			res, err := queryClient.QueryAuctions(
				context.Background(),
				&types.QueryAuctionsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "auctions")
	return cmd
}

func queryBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "biddings [bidder]",
		Short: "Query biddings by bidder address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			bidder, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryBiddings(
				context.Background(),
				&types.QueryBiddingsRequest{
					Bidder: bidder.String(),
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

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query module parameters",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryParams(
				context.Background(),
				&types.QueryParamsRequest{},
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
