package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/comdex-official/comdex/x/cdp/types"
	"github.com/cosmos/cosmos-sdk/client"
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
		QueryParams(),
		QueryCdps(),
		QueryCdpById())

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

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCDP(cmd.Context(), &types.QueryCDPRequest{
				Owner:          args[0],
				CollateralType: args[1],
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

func QueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "get the cdp module parameters",
		Long:  "get the current global cdp module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryParams(context.Background(), &types.QueryParamsRequest{})
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
		Short: "cdp's information",
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

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCDPs(cmd.Context(), &types.QueryCDPsRequest{
				Owner:          args[0],
				Pagination:     pagination,
			})

			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "cdps")
	cmd.Flags().String(Owner, "", "owner")
	return cmd
}

func QueryCdpById() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdpsid [id]",
		Short: "cdp by id information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			u64, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.QueryCDPById(cmd.Context(), &types.QueryCDPByIdRequest{
				Id: u64,
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