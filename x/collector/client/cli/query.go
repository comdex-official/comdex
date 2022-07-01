package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/comdex-official/comdex/x/collector/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group collector queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams(),
		QueryCollectorLookupByApp(),
		QueryCollectorLookupByAppAndAsset(),
		QueryCollectorDataByAppAndAsset(),
		QueryAuctionMappingForAppAndAsset(),
		QueryNetFeeCollectedForAppAndAsset())

	return cmd
}

// QueryCollectorLookupByApp query collector store by product.
func QueryCollectorLookupByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-by-app [app-id]",
		Short: "Query collector lookup for a app",
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

			res, err := queryClient.QueryCollectorLookupByApp(cmd.Context(), &types.QueryCollectorLookupByAppRequest{
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

// QueryCollectorLookupByAppAndAsset query collector store by product and asset.
func QueryCollectorLookupByAppAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-by-app-and-asset [app-id] [asset-id]",
		Short: "Query collector lookup for an app by asset id",
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

			res, err := queryClient.QueryCollectorLookupByAppAndAsset(cmd.Context(), &types.QueryCollectorLookupByAppAndAssetRequest{
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

// QueryCollectorDataByAppAndAsset query collector store by product.
func QueryCollectorDataByAppAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-data-by-app-and-asset [app-id] [asset_id]",
		Short: "Query collector data for an app and asset",
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

			res, err := queryClient.QueryCollectorDataByAppAndAsset(cmd.Context(), &types.QueryCollectorDataByAppAndAssetRequest{
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
		Use:   "auction-data-by-app-and-asset [app-id] [asset_id]",
		Short: "Query auction data for an app and asset",
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
		Use:   "net-fee-data-by-app-by-asset [app-id] [asset-id]",
		Short: "Query net fee data for an app and asset",
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
