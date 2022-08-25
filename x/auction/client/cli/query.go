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
		Use:   "surplus-auction [appid] [auction mapping id] [auction id] [history]",
		Short: "Query surplus auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QuerySurplusAuction(
				context.Background(),
				&types.QuerySurplusAuctionRequest{
					AppId:            appID,
					AuctionMappingId: auctionMappingID,
					AuctionId:        auctionID,
					History:          history,
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
		Use:   "all-surplus-auctions [appid] [history]",
		Short: "Query all surplus auctions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QuerySurplusAuctions(
				context.Background(),
				&types.QuerySurplusAuctionsRequest{
					AppId:      appID,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "all-surplus-auctions")
	return cmd
}

func querySurplusBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "surplus-biddings [bidder] [app-id] [history]",
		Short: "Query surplus biddings by bidder address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			bidder, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QuerySurplusBiddings(
				context.Background(),
				&types.QuerySurplusBiddingsRequest{
					AppId:      appID,
					History:    history,
					Bidder:     bidder.String(),
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
	flags.AddPaginationFlagsToCmd(cmd, "surplus-biddings")

	return cmd
}

func queryDebtAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debt-auction [appid] [auction mapping id] [auction id] [history]",
		Short: "Query Debt auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDebtAuction(
				context.Background(),
				&types.QueryDebtAuctionRequest{
					AppId:            appID,
					AuctionMappingId: auctionMappingID,
					AuctionId:        auctionID,
					History:          history,
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
		Use:   "all-debt-auctions [appid] [history]",
		Short: "Query Debt auctions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDebtAuctions(
				context.Background(),
				&types.QueryDebtAuctionsRequest{
					AppId:      appID,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "all-debt-auctions")
	return cmd
}

func queryDebtBidding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debt-biddings [bidder] [app-id] [history]",
		Short: "Query surplus Debt by bidder address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			bidder, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDebtBiddings(
				context.Background(),
				&types.QueryDebtBiddingsRequest{
					Bidder:     bidder.String(),
					AppId:      appID,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "debt-biddings")

	return cmd
}

func queryDutchAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-auction [appid] [auction mapping id] [auction id] [history]",
		Short: "Query Dutch auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchAuction(
				context.Background(),
				&types.QueryDutchAuctionRequest{
					AppId:            appID,
					AuctionMappingId: auctionMappingID,
					AuctionId:        auctionID,
					History:          history,
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
		Use:   "all-dutch-auctions [appid] [history]",
		Short: "Query Dutch auctions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchAuctions(
				context.Background(),
				&types.QueryDutchAuctionsRequest{
					AppId:      appID,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "all-dutch-auctions")
	return cmd
}

func queryProtocolStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-protocol-stats [appid]",
		Short: "Query protocol stats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryProtocolStatistics(
				context.Background(),
				&types.QueryProtocolStatisticsRequest{
					AppId:      appID,
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
	flags.AddPaginationFlagsToCmd(cmd, "all-protocol-stats")
	return cmd
}

func queryDutchBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-biddings [bidder] [app-id] [history]",
		Short: "Query Dutch biddings by bidder address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			bidder, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchBiddings(
				context.Background(),
				&types.QueryDutchBiddingsRequest{
					AppId:      appID,
					History:    history,
					Bidder:     bidder.String(),
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
	flags.AddPaginationFlagsToCmd(cmd, "dutch-biddings")

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

			queryClient := types.NewQueryClient(ctx)

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

func queryDutchLendAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-auction-lend [appid] [auction mapping id] [auction id] [history]",
		Short: "Query Dutch auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			auctionMappingID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			auctionID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchLendAuction(
				context.Background(),
				&types.QueryDutchLendAuctionRequest{
					AppId:            appID,
					AuctionMappingId: auctionMappingID,
					AuctionId:        auctionID,
					History:          history,
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

func queryDutchLendAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-dutch-auctions-lend [appid] [history]",
		Short: "Query Dutch auctions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchLendAuctions(
				context.Background(),
				&types.QueryDutchLendAuctionsRequest{
					AppId:      appID,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "all-dutch-auctions-lend")
	return cmd
}

func queryDutchLendBiddings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dutch-biddings-lend [bidder] [app-id] [history]",
		Short: "Query Dutch biddings by bidder address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			bidder, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryDutchLendBiddings(
				context.Background(),
				&types.QueryDutchLendBiddingsRequest{
					AppId:      appID,
					History:    history,
					Bidder:     bidder.String(),
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
	flags.AddPaginationFlagsToCmd(cmd, "dutch-biddings-lend")

	return cmd
}

func queryFilterDutchAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter-dutch-auctions [appid] [denom] [history]",
		Short: "Query Filtered Dutch auctions",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			denom := args[1]
			history, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryFilterDutchAuctions(
				context.Background(),
				&types.QueryFilterDutchAuctionsRequest{
					AppId:      appID,
					Denom:      denom,
					History:    history,
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
	flags.AddPaginationFlagsToCmd(cmd, "filter-dutch-auctions")
	return cmd
}
