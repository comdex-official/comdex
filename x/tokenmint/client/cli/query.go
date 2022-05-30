package cli

import (
	"fmt"
	// "strings"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group tokenmint queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryAllTokenMintedForAllProducts(),
		QueryTokenMintedByProduct(),
		QueryTokenMintedByProductAndAsset(),
	)
	// this line is used by starport scaffolding # 1

	return cmd
}

func QueryAllTokenMintedForAllProducts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-token-minted-all-products",
		Short: "Token minted tokens data",
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

			res, err := queryClient.QueryAllTokenMintedForAllProducts(cmd.Context(), &types.QueryAllTokenMintedForAllProductsRequest{
				Pagination: pagination,
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

func QueryTokenMintedByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-product [app_id]",
		Short: "Token minted by product",
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
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryTokenMintedByProduct(cmd.Context(), &types.QueryTokenMintedByProductRequest{
				AppId:      appId,
				Pagination: pagination,
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

func QueryTokenMintedByProductAndAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Token-minted-by-Product-asset [app_id] [asset_id]",
		Short: "Token minted by product and asset data",
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
			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryTokenMintedByProductAndAsset(cmd.Context(), &types.QueryTokenMintedByProductAndAssetRequest{
				AppId:      appId,
				AssetId:    assetId,
				Pagination: pagination,
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
