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

func querySurplusAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "surplus-auction [appid] [auction mapping id] [auction id]",
		Short: "Query surplus auction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QuerySurplusAuction(
				context.Background(),
				&types.QuerySurplusAuctionRequest{
					appId,
					auctionMappingId,
					auctionId,
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

func querySurplusAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-surplus-auctions",
		Short: "Query all surplus auctions",
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
			res, err := queryClient.QuerySurplusAuctions(
				context.Background(),
				&types.QuerySurplusAuctionsRequest{
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

func querySurplusBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suplus-biddings [bidder]",
		Short: "Query surplus biddings by bidder address",
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
			res, err := queryClient.QuerySurplusBiddings(
				context.Background(),
				&types.QuerySurplusBiddingsRequest{
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

func queryDebtAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debt-auction [appid] [auction mapping id] [auction id]",
		Short: "Query Debt auction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryDebtAuction(
				context.Background(),
				&types.QueryDebtAuctionRequest{
					appId,
					auctionMappingId,
					auctionId,
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

func queryDebtAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-debt-auctions",
		Short: "Query Debt auctions",
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
			res, err := queryClient.QueryDebtAuctions(
				context.Background(),
				&types.QueryDebtAuctionsRequest{
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

func queryDebtBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "surplus-debt-biddings [bidder]",
		Short: "Query surplus Debt by bidder address",
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
			res, err := queryClient.QueryDebtBiddings(
				context.Background(),
				&types.QueryDebtBiddingsRequest{
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

func queryDutchAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-auction [appid] [auction mapping id] [auction id]",
		Short: "Query Dutch auction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryDutchAuction(
				context.Background(),
				&types.QueryDutchAuctionRequest{
					appId,
					auctionMappingId,
					auctionId,
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

func queryDutchAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-dutch-auctions",
		Short: "Query Dutch auctions",
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
			res, err := queryClient.QueryDutchAuctions(
				context.Background(),
				&types.QueryDutchAuctionsRequest{
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

func queryDutchBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-biddings [bidder]",
		Short: "Query Dutch biddings by bidder address",
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
			res, err := queryClient.QueryDutchBiddings(
				context.Background(),
				&types.QueryDutchBiddingsRequest{
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
