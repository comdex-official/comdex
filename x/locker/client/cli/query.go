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
			queryClient := types.NewQueryServiceClient(ctx)
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

			assetId, err := strconv.ParseUint(args[1], 10, 64)

			queryClient := types.NewQueryServiceClient(ctx)
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

			queryClient := types.NewQueryServiceClient(ctx)
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

func queryTotalDepositByProductToAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-deposit-per-product-asset [product_id] [asset_id]",
		Short: "total deposit per product to asset id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)

			assetId, err := strconv.ParseUint(args[1], 10, 64)

			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryTotalDepositByProductAssetID(
				context.Background(),
				&types.QueryTotalDepositByProductAssetIDRequest{
					ProductId: productId,
					AssetId:   assetId,
					Owner:     ctx.GetFromAddress().String(),
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

func queryOwnerLockerByProductID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-product-id [product_id] [owner]",
		Short: "owner locker by product id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)
			owner := args[1]

			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryOwnerLockerByProductID(
				context.Background(),
				&types.QueryOwnerLockerByProductIDRequest{
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

func queryOwnerLockerOfAllProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-product-id [owner]",
		Short: "owner locker by all product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner := args[0]

			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryOwnerLockerOfAllProduct(
				context.Background(),
				&types.QueryOwnerLockerOfAllProductRequest{
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

func queryOwnerLockerByProductToAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-product-to-asset-id [product_id] [asset_id] [owner]",
		Short: "owner locker by product to asset id",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			productId, err := strconv.ParseUint(args[0], 10, 64)

			assetId, err := strconv.ParseUint(args[1], 10, 64)

			owner := args[2]

			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryOwnerLockerByProductToAssetID(
				context.Background(),
				&types.QueryOwnerLockerByProductToAssetIDRequest{
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

			queryClient := types.NewQueryServiceClient(ctx)
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

			assetId, err := strconv.ParseUint(args[1], 10, 64)

			queryClient := types.NewQueryServiceClient(ctx)
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

			queryClient := types.NewQueryServiceClient(ctx)
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

			queryClient := types.NewQueryServiceClient(ctx)
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
