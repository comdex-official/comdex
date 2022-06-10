package cli

import (
	"context"
	"strconv"

	"github.com/comdex-official/comdex/x/locker/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func queryLockedVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-info [id]",
		Short: "Query locker info by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id := args[0]
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerInfo(
				context.Background(),
				&types.QueryLockerInfoRequest{
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

func queryLockerByProductAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockers-by-product-asset-id [product_id] [asset_id]",
		Short: "Query all lockers by product and asset id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockersByProductToAssetID(
				context.Background(),
				&types.QueryLockersByProductToAssetIDRequest{
					ProductId: productId,
					AssetId:   assetId,
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

func queryLockerByProductID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-info-product-id [product_id]",
		Short: "Query locker info by product id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerInfoByProductID(
				context.Background(),
				&types.QueryLockerInfoByProductIDRequest{
					ProductId: productId,
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

func queryTotalDepositByProductAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-deposit-per-product-assetid [product_id] [asset_id]",
		Short: "total deposit per product to asset id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryTotalDepositByProductAssetID(
				context.Background(),
				&types.QueryTotalDepositByProductAssetIDRequest{
					ProductId: productId,
					AssetId:   assetId,
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

func queryOwnerLockerByProductIDbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-product-id-and-owner [product_id] [owner]",
		Short: "owner locker by product id by owner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			owner := args[1]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerByProductIDbyOwner(
				context.Background(),
				&types.QueryOwnerLockerByProductIDbyOwnerRequest{
					ProductId: productId,
					Owner:     owner,
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

func queryLockerByProductbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-by-product-by-owner [product_id] [owner]",
		Short: "locker by product by owner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner := args[1]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerByProductByOwner(
				context.Background(),
				&types.QueryLockerByProductByOwnerRequest{
					ProductId: productId,
					Owner:     owner,
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

func queryOwnerLockerOfAllProductbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-all-product-by-owner [owner]",
		Short: "owner locker by all product by owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner := args[0]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerOfAllProductByOwner(
				context.Background(),
				&types.QueryOwnerLockerOfAllProductByOwnerRequest{
					Owner: owner,
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

func queryOwnerLockerByProductToAssetIDbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-product-to-asset-id-owner [product_id] [asset_id] [owner]",
		Short: "owner locker by product to asset id and owner",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			owner := args[2]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerByProductToAssetIDbyOwner(
				context.Background(),
				&types.QueryOwnerLockerByProductToAssetIDbyOwnerRequest{
					ProductId: productId,
					AssetId:   assetId,
					Owner:     owner,
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

func queryOwnerTxDetailsLockerOfProductByOwnerByAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-tx-details-by-product-to-owner-by-asset [product_id] [owner] [asset_id]",
		Short: "owner locker tx details by product to owner by asset",
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

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner := args[1]

			assetId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerTxDetailsLockerOfProductByOwnerByAsset(
				context.Background(),
				&types.QueryOwnerTxDetailsLockerOfProductByOwnerByAssetRequest{
					ProductId: productId,
					Owner:     owner,
					AssetId: assetId,
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

	return cmd
}

func queryTotalLockerByProductID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-locker-by-product-id [product_id]",
		Short: "total locker by product id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerCountByProductID(
				context.Background(),
				&types.QueryLockerCountByProductIDRequest{
					ProductId: productId,
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

func queryTotalLockerByProductToAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-locker-by-product-to-asset-id [product_id] [asset_id]",
		Short: "total locker by product to asset id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerCountByProductToAssetID(
				context.Background(),
				&types.QueryLockerCountByProductToAssetIDRequest{
					ProductId: productId,
					AssetId:   assetId,
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

func queryWhiteListedAssetIDsByProductID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted-assetIds-by-product-id [product_id]",
		Short: "whitelisted asset Ids by product id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryWhiteListedAssetIDsByProductID(
				context.Background(),
				&types.QueryWhiteListedAssetIDsByProductIDRequest{
					ProductId: productId,
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

func queryWhiteListedAssetByAllProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted-asset-all-product",
		Short: "query white listed asset all product",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryWhiteListedAssetByAllProduct(
				context.Background(),
				&types.QueryWhiteListedAssetByAllProductRequest{},
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

func queryLockerLookupTableByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-lookup-by-app [app_id]",
		Short: "locker lookup by app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			app_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerLookupTableByApp(
				context.Background(),
				&types.QueryLockerLookupTableByAppRequest{
					AppId: app_id,
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

func queryLockerLookupTableByAppAndAssetId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-lookup-by-app-and-assetId [app_id] [asset_id]",
		Short: "locker lookup by app and assetId",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			app_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			asset_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerLookupTableByAppAndAssetId(
				context.Background(),
				&types.QueryLockerLookupTableByAppAndAssetIdRequest{
					AppId:   app_id,
					AssetId: asset_id,
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

func queryLockerTotalDepositedByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-deposited-by-app [app_id]",
		Short: "locker deposited amount by app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			app_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerTotalDepositedByApp(
				context.Background(),
				&types.QueryLockerTotalDepositedByAppRequest{
					AppId: app_id,
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

func queryState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state [address] [denom] [blockheight] [target]",
		Short: "state of an account at a blockheight",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			denom := args[1]
			blockheight := args[2]
			target := args[3]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryState(
				context.Background(),
				&types.QueryStateRequest{
					Address: address,
					Denom:   denom,
					Height:  blockheight,
					Target:  target,
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

func queryLockerTotalRewardsByAssetAppWise() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockers-rewards-by-product-asset-id [app_id] [asset_id]",
		Short: "Query all lockers rewards by product and asset id",
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

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerTotalRewardsByAssetAppWise(
				context.Background(),
				&types.QueryLockerTotalRewardsByAssetAppWiseRequest{
					AppId: appId,
					AssetId:   assetId,
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
