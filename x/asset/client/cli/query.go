package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/asset/types"
)

func queryAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset [id]",
		Short: "Query an asset data by asset id",
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

			res, err := queryClient.QueryAsset(
				context.Background(),
				&types.QueryAssetRequest{
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

func queryAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Query data of all the assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssets(
				context.Background(),
				&types.QueryAssetsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "assets")

	return cmd
}

func queryPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair [id]",
		Short: "Query a pair data by pair id",
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

			res, err := queryClient.QueryAssetPair(
				context.Background(),
				&types.QueryAssetPairRequest{
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

func queryPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pairs",
		Short: "Query data of all the pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAssetPairs(
				context.Background(),
				&types.QueryAssetPairsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "pairs")

	return cmd
}

func queryApps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Query data for all the apps of the protocol",
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

			res, err := queryClient.QueryApps(
				context.Background(),
				&types.QueryAppsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "apps")

	return cmd
}

func queryApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app [id]",
		Short: "Query data for an app of the protocol by app id",
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

			res, err := queryClient.QueryApp(
				context.Background(),
				&types.QueryAppRequest{
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

func queryExtendedPairVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended_pair_vault [id]",
		Short: "Query an Extended pairVault by id",
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

			res, err := queryClient.QueryExtendedPairVault(
				context.Background(),
				&types.QueryExtendedPairVaultRequest{
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

func queryAllExtendedPairVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended_pair_vaults",
		Short: "Query all Extended pairVaults",
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

			res, err := queryClient.QueryAllExtendedPairVaults(
				context.Background(),
				&types.QueryAllExtendedPairVaultsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extended_pair_vaults")

	return cmd
}

func queryAllExtendedPairVaultsByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app-extended-pair [app_id]",
		Short: "Query all extended pairs in an app",
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

			res, err := queryClient.QueryAllExtendedPairVaultsByApp(
				context.Background(),
				&types.QueryAllExtendedPairVaultsByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "app-extended-pair")

	return cmd
}

func queryAllExtendedPairStableVaultsIDByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended-pair-stable-vault-wise [app_id]",
		Short: "Query all extended pairs stable vault wise in an app",
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

			res, err := queryClient.QueryAllExtendedPairStableVaultsIDByApp(
				context.Background(),
				&types.QueryAllExtendedPairStableVaultsIDByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extended-pair-stable-vault-wise")

	return cmd
}

func queryGovTokenByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check_for_gov_token_for_an_app [app_id]",
		Short: "Query for gov token in an app",
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

			res, err := queryClient.QueryGovTokenByApp(
				context.Background(),
				&types.QueryGovTokenByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "check_for_gov_token_for_an_app")

	return cmd
}

func queryAllExtendedPairStableVaultsByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended-pair-data-stable-vault-wise [app_id]",
		Short: "Query all extended pairs data stable vault wise in an app",
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

			res, err := queryClient.QueryAllExtendedPairStableVaultsByApp(
				context.Background(),
				&types.QueryAllExtendedPairStableVaultsByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extended-pair-data-stable-vault-wise")

	return cmd
}

func queryExtendedPairStableVaultsByAppWithoutStable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended_pair_stable_vault_data_without_stable [app_id]",
		Short: "Query all extended pairs data without stable vault in an app",
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

			res, err := queryClient.QueryExtendedPairVaultsByAppWithoutStable(
				context.Background(),
				&types.QueryExtendedPairVaultsByAppWithoutStableRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extended_pair_stable_vault_data_without_stable")

	return cmd
}
