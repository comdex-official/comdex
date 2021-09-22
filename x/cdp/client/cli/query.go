package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/cdp/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryCdp(),
		QueryCdps(),
	)

	return cmd
}

func QueryCdp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdp [owner-addr] [collateralType]",
		Short: "cdp's information",
		Args:  cobra.ExactArgs(2),
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

			res, err := queryClient.QueryCDP(cmd.Context(), &types.QueryCDPRequest{
				Id: id,
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

func QueryCdps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdps [owner-addr]",
		Short: "cdps list for an individual account",
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

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCDPs(cmd.Context(), &types.QueryCDPsRequest{
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
