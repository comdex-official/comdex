package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strconv"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/auctionsV2/types"
)

func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group auctions queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryAuction(),
		queryAuctions(),
	)

	return cmd
}

func queryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction [auction id] [history]",
		Short: "Query auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			auctionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.Auction(
				context.Background(),
				&types.QueryAuctionRequest{
					AuctionId: auctionID,
					History:   history,
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
		Use:   "auctions [history]",
		Short: "Query all auctions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			history, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.Auctions(
				context.Background(),
				&types.QueryAuctionsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "auctions")
	return cmd
}
