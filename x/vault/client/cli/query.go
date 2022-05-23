package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/vault/types"
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
		QueryAllVaults(),
		QueryVault(),
		QueryAllVaultsByAppAndExtendedPair(),
		// QueryVaults(),
		QueryVaultOfOwnerByPair(),
		QueryVaultByProduct(),
		QueryAllVaultByOwner(),
		QueryTokenMintedAllProductsByPair(),
		QueryVaultCountByProduct(),
		QueryVaultCountByProductAndPair(),
		QueryTokenMintedAllProducts(),
		QueryTotalValueLockedByProductExtendedPair(),
		QueryExtendedPairIDByProduct(),
		QueryStableVaultInfo(),
		QueryAllStableVaults(),
		QueryStableVaultByProductExtendedPair(),
		
	)

	return cmd
}

func QueryAllVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults",
		Short: "list of all vaults available",
		RunE: func(cmd *cobra.Command, args []string) error {

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryAllVaults(cmd.Context(), &types.QueryAllVaultsRequest{
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

func QueryVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault [id]",
		Short: "vault's information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id := args[0]

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryVault(cmd.Context(), &types.QueryVaultRequest{
				Id: id,
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

func QueryVaultOfOwnerByPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-of-owner-by-pair [product_id] [owner] [extended_pair_id]",
		Short: "vaults list for an individual account",
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

			extendedPairid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryVaultOfOwnerByPair(cmd.Context(), &types.QueryVaultOfOwnerByPairRequest{
				ProductId: productId,
				Owner:      args[1],
				ExtendedPairId: extendedPairid,
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

func QueryAllVaultsByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-app-and-extended-pair [app_id] [extended_pair_id]",
		Short: "vaults list by app and extended pair",
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
			app_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryAllVaultsByAppAndExtendedPair(cmd.Context(), &types.QueryAllVaultsByAppAndExtendedPairRequest{
				AppId: app_id,
				ExtendedPairId: extended_pair_id,
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

// func QueryVaults() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "vaults [owner]",
// 		Short: "vaults list for an individual account",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {

// 			pagination, err := client.ReadPageRequest(cmd.Flags())
// 			if err != nil {
// 				return err
// 			}

// 			ctx, err := client.GetClientQueryContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			queryClient := types.NewQueryServiceClient(ctx)

// 			res, err := queryClient.QueryVaults(cmd.Context(), &types.QueryVaultsRequest{
// 				Owner:      args[0],
// 				Pagination: pagination,
// 			})

// 			if err != nil {
// 				return err
// 			}
// 			return ctx.PrintProto(res)
// 		},
// 	}

// 	flags.AddQueryFlagsToCmd(cmd)
// 	return cmd
// }

func QueryVaultByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-product [product_id]",
		Short: "vaults list for a product",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryVaultByProduct(cmd.Context(), &types.QueryVaultByProductRequest{
				ProductId: productId,
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

func QueryAllVaultByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-owner [owner]",
		Short: "vaults list for a product by owner",
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

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryAllVaultByOwner(cmd.Context(), &types.QueryAllVaultByOwnerRequest{
				Owner : args[0],
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


func QueryTokenMintedAllProductsByPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-products-extended-pair [product_id] [extended_pair_id]",
		Short: "token minted by products and extended pair",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryTokenMintedAllProductsByPair(cmd.Context(), &types.QueryTokenMintedAllProductsByPairRequest{
				ProductId: productId,
				ExtendedPairId :extended_pair_id,
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

func QueryTokenMintedAllProducts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-products [product_id]",
		Short: "token minted by products",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryTokenMintedAllProducts(cmd.Context(), &types.QueryTokenMintedAllProductsRequest{
				ProductId: productId,
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

func QueryVaultCountByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-count-by-products [product_id]",
		Short: "vault count by products",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryVaultCountByProduct(cmd.Context(), &types.QueryVaultCountByProductRequest{
				ProductId: productId,
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

func QueryVaultCountByProductAndPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-count-by-products-and-pair [product_id] [extended_pair_id]",
		Short: "vault count by products and extended pair",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryVaultCountByProductAndPair(cmd.Context(), &types.QueryVaultCountByProductAndPairRequest{
				ProductId: productId,
				ExtendedPairId: extended_pair_id,
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

func QueryTotalValueLockedByProductExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "value-locked-by-product-extended-pair [product_id] [extended_pair_id]",
		Short: "value locked by product extended pair",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryTotalValueLockedByProductExtendedPair(cmd.Context(), &types.QueryTotalValueLockedByProductExtendedPairRequest{
				ProductId: productId,
				ExtendedPairId: extended_pair_id,
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

func QueryExtendedPairIDByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended-pair-by-product [product_id]",
		Short: "value locked by product in extended pair",
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
			productId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryExtendedPairIDByProduct(cmd.Context(), &types.QueryExtendedPairIDByProductRequest{
				ProductId: productId,
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

func QueryStableVaultInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-id [stable_vault_id]",
		Short: "get stable vault by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryStableVaultInfo(cmd.Context(), &types.QueryStableVaultInfoRequest{
				StableVaultId: args[0],
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

func QueryAllStableVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-product [app_id]",
		Short: "get stable vault by product",
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

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryAllStableVaults(cmd.Context(), &types.QueryAllStableVaultsRequest{
				AppId: app_id,
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

func QueryStableVaultByProductExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-product-extendedPair [app_id] [extended_pair_id]",
		Short: "get stable vault by product and extended pair",
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
			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.QueryStableVaultByProductExtendedPair(cmd.Context(), &types.QueryStableVaultByProductExtendedPairRequest{
				AppId: app_id,
				ExtendedPairId: extended_pair_id,
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