package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		QueryParams())

	return cmd
}

func QueryCdp() *cobra.Command {
	return &cobra.Command{
		Use:   "cdp [owner-addr] [collateral-type]",
		Short: "cdp's information",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var (
				_, error = sdk.AccAddressFromBech32(args[0])
			)
			if error != nil {
				return error
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCDP(context.Background(), &types.QueryCDPRequest{

			})

			if err != nil {
				return err
			}
			flags.AddQueryFlagsToCmd(cmd)
			return ctx.PrintProto(res)
		},
	}
}



func QueryParams() *cobra.Command {
	return &cobra.Command{
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

			res, err := queryClient.QueryParams(context.Background(), &types.QueryParamsRequest{

			})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
}