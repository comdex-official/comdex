package cli

import (
	"fmt"
	// "strings"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/collector/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group collector queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams(),
	QueryCollectorLookupByProduct(),
	QueryCollectorLookupByProductAndAsset(),)
	// this line is used by starport scaffolding # 1

	return cmd 
}

func QueryCollectorLookupByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-by-product [app-id]",
		Short: "collector lookup for a product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCollectorLookupByProduct(cmd.Context(), &types.QueryCollectorLookupByProductRequest{
				AppId: appId,
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


func QueryCollectorLookupByProductAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-by-product-by-asset [app-id] [asset-id]",
		Short: "collector lookup for a product by asset id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryCollectorLookupByProductAndAsset(cmd.Context(), &types.QueryCollectorLookupByProductAndAssetRequest{
				AppId: appId,
				AssetId: assetId,
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