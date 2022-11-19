package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/petrichormoney/petri/x/locker/types"
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
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
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

func queryLockersByAppToAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockers-by-app-asset-id [app_id] [asset_id]",
		Short: "Query all lockers by app and asset id",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockersByAppToAssetID(
				context.Background(),
				&types.QueryLockersByAppToAssetIDRequest{
					AppId:      appID,
					AssetId:    assetID,
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
	flags.AddPaginationFlagsToCmd(cmd, "lockers-by-app-asset-id")

	return cmd
}

func queryLockerByAppID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-info-product-id [app_id]",
		Short: "Query locker info by app id",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerInfoByAppID(
				context.Background(),
				&types.QueryLockerInfoByAppIDRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "locker-info-product-id")

	return cmd
}

func queryTotalDepositByAppAndAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-deposit-per-app-assetid [app_id] [asset_id]",
		Short: "Query total deposit per app to asset id",
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
			res, err := queryClient.QueryTotalDepositByAppAndAssetID(
				context.Background(),
				&types.QueryTotalDepositByAppAndAssetIDRequest{
					AppId:   appID,
					AssetId: assetID,
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

func queryOwnerLockerByAppIDbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-app-id-and-owner [app_id] [owner]",
		Short: "Query owner locker by app id by owner",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			owner := args[1]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerByAppIDbyOwner(
				context.Background(),
				&types.QueryOwnerLockerByAppIDbyOwnerRequest{
					AppId:      appID,
					Owner:      owner,
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
	flags.AddPaginationFlagsToCmd(cmd, "owner-locker-by-app-id-and-owner")

	return cmd
}

func queryLockerByAppByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-by-app-by-owner [app_id] [owner]",
		Short: "Query locker by app by owner",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner := args[1]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerByAppByOwner(
				context.Background(),
				&types.QueryLockerByAppByOwnerRequest{
					AppId:      appID,
					Owner:      owner,
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
	flags.AddPaginationFlagsToCmd(cmd, "locker-by-app-by-owner")

	return cmd
}

func queryOwnerLockerOfAllAppsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-all-apps-by-owner [owner]",
		Short: "Query owner locker by all apps by owner",
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

			owner := args[0]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerOfAllAppsByOwner(
				context.Background(),
				&types.QueryOwnerLockerOfAllAppsByOwnerRequest{
					Owner:      owner,
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

func queryOwnerLockerByAppToAssetIDbyOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-locker-by-app-to-asset-id-owner [app_id] [asset_id] [owner]",
		Short: "Query owner locker by app to asset id and owner",
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

			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			owner := args[2]

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerLockerByAppToAssetIDbyOwner(
				context.Background(),
				&types.QueryOwnerLockerByAppToAssetIDbyOwnerRequest{
					AppId:   appID,
					AssetId: assetID,
					Owner:   owner,
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

func queryOwnerTxDetailsLockerOfAppByOwnerByAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-tx-details-by-app-to-owner-by-asset [app_id] [owner] [asset_id]",
		Short: "Query owner locker tx details by app to owner by asset",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner := args[1]

			assetID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryOwnerTxDetailsLockerOfAppByOwnerByAsset(
				context.Background(),
				&types.QueryOwnerTxDetailsLockerOfAppByOwnerByAssetRequest{
					AppId:      appID,
					Owner:      owner,
					AssetId:    assetID,
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

func queryTotalLockerByAppID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-locker-by-app-id [app_id]",
		Short: "Query total locker by app id",
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
			res, err := queryClient.QueryLockerCountByAppID(
				context.Background(),
				&types.QueryLockerCountByAppIDRequest{
					AppId: appID,
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

func queryTotalLockerByAppToAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-locker-by-app-to-asset-id [app_id] [asset_id]",
		Short: "Query total locker by app to asset id",
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
			res, err := queryClient.QueryLockerCountByAppToAssetID(
				context.Background(),
				&types.QueryLockerCountByAppToAssetIDRequest{
					AppId:   appID,
					AssetId: assetID,
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

func queryWhiteListedAssetIDsByAppID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted-assetIds-by-app-id [app_id]",
		Short: "Query whitelisted asset Ids by app id",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryWhiteListedAssetIDsByAppID(
				context.Background(),
				&types.QueryWhiteListedAssetIDsByAppIDRequest{
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

	return cmd
}

func queryWhiteListedAssetByAllApps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted-asset-all-apps",
		Short: "Query white listed asset all apps",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryWhiteListedAssetByAllApps(
				context.Background(),
				&types.QueryWhiteListedAssetByAllAppsRequest{
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

func queryLockerLookupTableByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-lookup-by-app [app_id]",
		Short: "Query locker lookup by app",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerLookupTableByApp(
				context.Background(),
				&types.QueryLockerLookupTableByAppRequest{
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

	return cmd
}

func queryLockerLookupTableByAppAndAssetID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locker-lookup-by-app-and-assetId [app_id] [asset_id]",
		Short: "Query locker lookup by app and assetId",
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
			res, err := queryClient.QueryLockerLookupTableByAppAndAssetID(
				context.Background(),
				&types.QueryLockerLookupTableByAppAndAssetIDRequest{
					AppId:   appID,
					AssetId: assetID,
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
		Short: "Query locker deposited amount by app",
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

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockerTotalDepositedByApp(
				context.Background(),
				&types.QueryLockerTotalDepositedByAppRequest{
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

	return cmd
}

func queryLockerTotalRewardsByAssetAppWise() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockers-rewards-by-app-asset-id [app_id] [asset_id]",
		Short: "Query all lockers rewards by app and asset id",
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
			res, err := queryClient.QueryLockerTotalRewardsByAssetAppWise(
				context.Background(),
				&types.QueryLockerTotalRewardsByAssetAppWiseRequest{
					AppId:   appID,
					AssetId: assetID,
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
