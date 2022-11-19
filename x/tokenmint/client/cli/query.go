package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/petrichormoney/petri/x/tokenmint/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group tokenmint queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "tokenmint",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "tokenmint"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryAllTokenMintedForAllApps(),
		QueryTokenMintedByApp(),
		QueryTokenMintedByAppAndAsset(),
	)

	return cmd
}

// QueryAllTokenMintedForAllApps Queries the total token minted for all the apps on petri.
func QueryAllTokenMintedForAllApps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-token-minted-all-apps",
		Short: "Token minted by all apps",
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

			res, err := queryClient.QueryAllTokenMintedForAllApps(cmd.Context(), &types.QueryAllTokenMintedForAllAppsRequest{
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "total-token-minted-all-apps")
	return cmd
}

// QueryTokenMintedByApp queries token minted per application/product.
func QueryTokenMintedByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-app [app_id]",
		Short: "Token minted by app",
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

			res, err := queryClient.QueryTokenMintedByApp(cmd.Context(), &types.QueryTokenMintedByAppRequest{
				AppId:      appID,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "token-minted-by-app")
	return cmd
}

// QueryTokenMintedByAppAndAsset queries token minted for an application/product and asset.
func QueryTokenMintedByAppAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-app-and-asset [app_id] [asset_id]",
		Short: "Token minted by app and asset data",
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

			res, err := queryClient.QueryTokenMintedByAppAndAsset(cmd.Context(), &types.QueryTokenMintedByAppAndAssetRequest{
				AppId:      appID,
				AssetId:    assetID,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "token-minted-by-app-and-asset")
	return cmd
}
