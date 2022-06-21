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
		QueryCollectorLookupByProductAndAsset(),
		QueryCollectorDataByProductAndAsset(),
		QueryAuctionMappingForAppAndAsset(),
		QueryNetFeeCollectedForAppAndAsset())

	return cmd
}

// QueryCollectorLookupByProduct query collector store by product
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
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCollectorLookupByProduct(cmd.Context(), &types.QueryCollectorLookupByProductRequest{
				AppId: appID,
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

// QueryCollectorLookupByProductAndAsset query collector store by product and asset
func QueryCollectorLookupByProductAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-by-product-and-asset [app-id] [asset-id]",
		Short: "collector lookup for a product by asset id",
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
			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCollectorLookupByProductAndAsset(cmd.Context(), &types.QueryCollectorLookupByProductAndAssetRequest{
				AppId:   appID,
				AssetId: assetID,
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

// QueryCollectorDataByProductAndAsset query collector store by product
func QueryCollectorDataByProductAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-data-by-product-and-asset [app-id] [asset_id]",
		Short: "collector data for a product and asset",
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
			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCollectorDataByProductAndAsset(cmd.Context(), &types.QueryCollectorDataByProductAndAssetRequest{
				AppId:   appID,
				AssetId: assetID,
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

func QueryAuctionMappingForAppAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction-data-by-product-and-asset [app-id] [asset_id]",
		Short: "auction data for a product and asset",
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
			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAuctionMappingForAppAndAsset(cmd.Context(), &types.QueryAuctionMappingForAppAndAssetRequest{
				AppId:   appID,
				AssetId: assetID,
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

func QueryNetFeeCollectedForAppAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "net-fee-data-by-product-by-asset [app-id] [asset-id]",
		Short: "net fee data for a product and asset",
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
			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryNetFeeCollectedForAppAndAsset(cmd.Context(), &types.QueryNetFeeCollectedForAppAndAssetRequest{
				AppId:   appID,
				AssetId: assetID,
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
