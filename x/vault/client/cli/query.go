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
		QueryAllVaultsByProduct(),
		QueryVault(),
		QueryAllVaultsByAppAndExtendedPair(),
		QueryVaultInfo(),
		QueryVaultOfOwnerByExtendedPair(),
		QueryVaultByProduct(),
		QueryAllVaultByOwner(),
		QueryTokenMintedAllProductsByPair(),
		QueryVaultCountByProduct(),
		QueryVaultCountByProductAndPair(),
		QueryTokenMintedByProductAssetWise(),
		QueryTotalValueLockedByProductExtendedPair(),
		QueryExtendedPairIDByProduct(),
		QueryStableVaultInfo(),
		QueryAllStableVaults(),
		QueryStableVaultByProductExtendedPair(),
		QueryExtendedPairVaultMappingByApp(),
		QueryExtendedPairVaultMappingByAppAndExtendedPairID(),
		QueryExtendedPairVaultMappingByOwnerAndApp(),
		QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairID(),
		QueryVaultInfoByAppByOwner(),
		QueryTVLLockedByAppOfAllExtendedPairs(),
		QueryTotalTVLByApp(),
		QueryUserMyPositionByApp(),
		QueryUserExtendedPairTotalData(),
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

			queryClient := types.NewQueryClient(ctx)

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

func QueryAllVaultsByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-product [appID]",
		Short: "list of all vaults available in a product",
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

			res, err := queryClient.QueryAllVaultsByProduct(cmd.Context(), &types.QueryAllVaultsByProductRequest{
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

			queryClient := types.NewQueryClient(ctx)

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

func QueryVaultOfOwnerByExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-of-owner-by-extended-pair [product_id] [owner] [extendedPairID]",
		Short: "vaults list for an individual account by extended pair",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultOfOwnerByExtendedPair(cmd.Context(), &types.QueryVaultOfOwnerByExtendedPairRequest{
				ProductId:      productID,
				Owner:          args[1],
				ExtendedPairId: extendedPairid,
				Pagination:     pagination,
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

func QueryVaultInfoByAppByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultInfoByAppByOwner [appID] [owner]",
		Short: "vaults list for an owner by App",
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

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultInfoByAppByOwner(cmd.Context(), &types.QueryVaultInfoByAppByOwnerRequest{
				AppId: appID,
				Owner: args[1],
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
		Use:   "vaults-by-app-and-extended-pair [appID] [extendedPairID]",
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
			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAllVaultsByAppAndExtendedPair(cmd.Context(), &types.QueryAllVaultsByAppAndExtendedPairRequest{
				AppId:          appID,
				ExtendedPairId: extendedPairID,
				Pagination:     pagination,
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

func QueryVaultInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultsInfo [id]",
		Short: "vaults list for an id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultInfo(cmd.Context(), &types.QueryVaultInfoRequest{
				Id: args[0],
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

func QueryVaultByProduct() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairvaults-by-product [product_id]",
		Short: "extended pair vaults list for a product",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultByProduct(cmd.Context(), &types.QueryVaultByProductRequest{
				ProductId:  productID,
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

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAllVaultByOwner(cmd.Context(), &types.QueryAllVaultByOwnerRequest{
				Owner:      args[0],
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
		Use:   "token-minted-by-products-extended-pair [product_id] [extendedPairID]",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryTokenMintedAllProductsByPair(cmd.Context(), &types.QueryTokenMintedAllProductsByPairRequest{
				ProductId:      productID,
				ExtendedPairId: extendedPairID,
				Pagination:     pagination,
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

func QueryTokenMintedByProductAssetWise() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-products-asset-wise [product_id]",
		Short: "token minted by products asset wise",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryTokenMintedByProductAssetWise(cmd.Context(), &types.QueryTokenMintedByProductAssetWiseRequest{
				ProductId:  productID,
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultCountByProduct(cmd.Context(), &types.QueryVaultCountByProductRequest{
				ProductId:  productID,
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
		Use:   "vault-count-by-products-and-pair [product_id] [extendedPairID]",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultCountByProductAndPair(cmd.Context(), &types.QueryVaultCountByProductAndPairRequest{
				ProductId:      productID,
				ExtendedPairId: extendedPairID,
				Pagination:     pagination,
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
		Use:   "value-locked-by-product-extended-pair [product_id] [extendedPairID]",
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryTotalValueLockedByProductExtendedPair(cmd.Context(), &types.QueryTotalValueLockedByProductExtendedPairRequest{
				ProductId:      productID,
				ExtendedPairId: extendedPairID,
				Pagination:     pagination,
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
			productID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryExtendedPairIDByProduct(cmd.Context(), &types.QueryExtendedPairIDByProductRequest{
				ProductId:  productID,
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

			queryClient := types.NewQueryClient(ctx)

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
		Use:   "stable-vault-by-product [appID]",
		Short: "get stable vault by product",
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

			res, err := queryClient.QueryAllStableVaults(cmd.Context(), &types.QueryAllStableVaultsRequest{
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

func QueryStableVaultByProductExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-product-extendedPair [appID] [extendedPairID]",
		Short: "get stable vault by product and extended pair",
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
			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryStableVaultByProductExtendedPair(cmd.Context(), &types.QueryStableVaultByProductExtendedPairRequest{
				AppId:          appID,
				ExtendedPairId: extendedPairID,
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

func QueryExtendedPairVaultMappingByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairVault-by-product [appID]",
		Short: "get ExtendedPair Vault By App",
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

			res, err := queryClient.QueryExtendedPairVaultMappingByApp(cmd.Context(), &types.QueryExtendedPairVaultMappingByAppRequest{
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

func QueryExtendedPairVaultMappingByAppAndExtendedPairID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairVault-by-product-and-ExtendedPairId [appID] [extendedPairID]",
		Short: "get ExtendedPair Vault Mapping By App And ExtendedPairId",
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
			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryExtendedPairVaultMappingByAppAndExtendedPairId(cmd.Context(), &types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdRequest{
				AppId:          appID,
				ExtendedPairId: extendedPairID,
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

func QueryExtendedPairVaultMappingByOwnerAndApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairVault-by-owner-and-product [owner] [appID]",
		Short: "get ExtendedPair Vault Mapping By owner and App",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryExtendedPairVaultMappingByOwnerAndApp(cmd.Context(), &types.QueryExtendedPairVaultMappingByOwnerAndAppRequest{
				Owner: args[0],
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

func QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairVault-by-owner-product-and-extended-pair [owner] [appID] [extendedPair]",
		Short: "get ExtendedPair Vault Mapping By owner App and extended pair",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			appID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			extendedPair, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairID(cmd.Context(), &types.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairIDRequest{
				Owner:        args[0],
				AppId:        appID,
				ExtendedPair: extendedPair,
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

func QueryTVLLockedByAppOfAllExtendedPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tvl-locked-by-app-all-extended-pairs [appID]",
		Short: "get tvl locked By App of all extended pairs",
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

			res, err := queryClient.QueryTVLLockedByAppOfAllExtendedPairs(cmd.Context(), &types.QueryTVLLockedByAppOfAllExtendedPairsRequest{
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

func QueryTotalTVLByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tvl-locked-by-app [appID]",
		Short: "get tvl locked By App",
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

			res, err := queryClient.QueryTotalTVLByApp(cmd.Context(), &types.QueryTotalTVLByAppRequest{
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

func QueryUserMyPositionByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-my-position-by-app [appID] [owner]",
		Short: "user my position by app",
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

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryUserMyPositionByApp(cmd.Context(), &types.QueryUserMyPositionByAppRequest{
				AppId: appID,
				Owner: args[1],
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

func QueryUserExtendedPairTotalData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-extended-total-data [owner]",
		Short: "User Extended Pair Total Data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryUserExtendedPairTotalData(cmd.Context(), &types.QueryUserExtendedPairTotalDataRequest{
				Owner: args[0],
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
