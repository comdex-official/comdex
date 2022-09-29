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
		QueryAllVaultsByApp(),
		QueryVault(),
		QueryAllVaultsByAppAndExtendedPair(),
		QueryVaultInfoByVaultID(),
		QueryVaultIDOfOwnerByExtendedPairAndApp(),
		QueryVaultIdsByAppInAllExtendedPairs(),
		QueryAllVaultIdsByAnOwner(),
		QueryTokenMintedByAppAndExtendedPair(),
		QueryVaultCountByApp(),
		QueryVaultCountByAppAndExtendedPair(),
		QueryTokenMintedAssetWiseByApp(),
		QueryTotalValueLockedByAppAndExtendedPair(),
		QueryExtendedPairIDsByApp(),
		QueryStableVaultByVaultID(),
		QueryStableVaultByApp(),
		QueryStableVaultByAppAndExtendedPair(),
		QueryExtendedPairVaultMappingByApp(),
		QueryExtendedPairVaultMappingByAppAndExtendedPair(),
		QueryVaultInfoOfOwnerByApp(),
		QueryTVLByAppOfAllExtendedPairs(),
		QueryTVLByApp(),
		QueryUserMyPositionByApp(),
		QueryUserExtendedPairTotalData(),
		QueryPairsLockedAndMintedStatisticByApp(),
	)

	return cmd
}

func QueryAllVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults",
		Short: "Query all vaults in all apps",
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
	flags.AddPaginationFlagsToCmd(cmd, "vaults")
	return cmd
}

func QueryAllVaultsByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-app [appID]",
		Short: "Query all vaults by app id",
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

			res, err := queryClient.QueryAllVaultsByApp(cmd.Context(), &types.QueryAllVaultsByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "vaults-by-app")
	return cmd
}

func QueryVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault [id]",
		Short: "Query vault by vault id",
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

func QueryVaultIDOfOwnerByExtendedPairAndApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-id-of-owner-by-extended-pair-and-app [app_id] [owner] [extendedPairID]",
		Short: "Query vault id for an individual account by extended pair id and app",
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

			extendedPairid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultIDOfOwnerByExtendedPairAndApp(cmd.Context(), &types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest{
				AppId:          appID,
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
	flags.AddPaginationFlagsToCmd(cmd, "vault-id-of-owner-by-extended-pair-and-app")
	return cmd
}

func QueryVaultInfoOfOwnerByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultsInfo-of-owner-by-app [appID] [owner]",
		Short: "Query vaultsInfo of an owner by App",
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

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryVaultInfoOfOwnerByApp(cmd.Context(), &types.QueryVaultInfoOfOwnerByAppRequest{
				AppId:      appID,
				Owner:      args[1],
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "vaultsInfo-of-owner-by-app")
	return cmd
}

func QueryAllVaultsByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaults-by-app-and-extended-pair [appID] [extendedPairID]",
		Short: "Query all vaults by app and extended pair",
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
	flags.AddPaginationFlagsToCmd(cmd, "vaults-by-app-and-extended-pair")
	return cmd
}

func QueryVaultInfoByVaultID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultsInfo-by-vault-id [id]",
		Short: "Query vaultsInfo by vault id",
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

			res, err := queryClient.QueryVaultInfoByVaultID(cmd.Context(), &types.QueryVaultInfoByVaultIDRequest{
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

func QueryVaultIdsByAppInAllExtendedPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultIds-by-app-in-all-extendedPairs [app_id]",
		Short: "Query VaultIds ByApp In All ExtendedPairs",
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

			res, err := queryClient.QueryVaultIdsByAppInAllExtendedPairs(cmd.Context(), &types.QueryVaultIdsByAppInAllExtendedPairsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "vaultIds-by-app-in-all-extendedPairs")
	return cmd
}

func QueryAllVaultIdsByAnOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-ids-by-an-owner [owner]",
		Short: "Query all vault ids by an owner",
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

			res, err := queryClient.QueryAllVaultIdsByAnOwner(cmd.Context(), &types.QueryAllVaultIdsByAnOwnerRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "vault-ids-by-an-owner")
	return cmd
}

func QueryTokenMintedByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-by-App-and-extended-pair [app_id] [extendedPairID]",
		Short: "Query token minted by App and extended pair",
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

			res, err := queryClient.QueryTokenMintedByAppAndExtendedPair(cmd.Context(), &types.QueryTokenMintedByAppAndExtendedPairRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "token-minted-by-App-and-extended-pair")
	return cmd
}

func QueryTokenMintedAssetWiseByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minted-asset-wise-by-app [app_id]",
		Short: "Query token minted asset wise in an App",
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

			res, err := queryClient.QueryTokenMintedAssetWiseByApp(cmd.Context(), &types.QueryTokenMintedAssetWiseByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "token-minted-asset-wise-by-app")
	return cmd
}

func QueryVaultCountByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-count-by-an-App [app_id]",
		Short: "Query vault count by an App",
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

			res, err := queryClient.QueryVaultCountByApp(cmd.Context(), &types.QueryVaultCountByAppRequest{
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

func QueryVaultCountByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault-count-by-an-App-and-extended-pair [app_id] [extendedPairID]",
		Short: "Query vault count by an App and extended pair",
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

			res, err := queryClient.QueryVaultCountByAppAndExtendedPair(cmd.Context(), &types.QueryVaultCountByAppAndExtendedPairRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "vault-count-by-an-App-and-extended-pair")
	return cmd
}

func QueryTotalValueLockedByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "value-locked-by-App-and-extended-pair [app_id] [extendedPairID]",
		Short: "Query value locked in an App and extended pair",
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

			res, err := queryClient.QueryTotalValueLockedByAppAndExtendedPair(cmd.Context(), &types.QueryTotalValueLockedByAppAndExtendedPairRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "value-locked-by-App-and-extended-pair")
	return cmd
}

func QueryExtendedPairIDsByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extended-pair-id-by-App [app_id]",
		Short: "Query extended pair ids by an App",
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

			res, err := queryClient.QueryExtendedPairIDsByApp(cmd.Context(), &types.QueryExtendedPairIDsByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extended-pair-id-by-App")
	return cmd
}

func QueryStableVaultByVaultID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-id [stable_vault_id]",
		Short: "Query stable vault by vault id",
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

			res, err := queryClient.QueryStableVaultByVaultID(cmd.Context(), &types.QueryStableVaultByVaultIDRequest{
				StableVaultId: id,
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

func QueryStableVaultByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-App [appID]",
		Short: "Query stable vault by App",
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

			res, err := queryClient.QueryStableVaultByApp(cmd.Context(), &types.QueryStableVaultByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "stable-vault-by-App")
	return cmd
}

func QueryStableVaultByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stable-vault-by-App-and-extendedPair [appID] [extendedPairID]",
		Short: "Query stable vault by App and extended pair",
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

			res, err := queryClient.QueryStableVaultByAppAndExtendedPair(cmd.Context(), &types.QueryStableVaultByAppAndExtendedPairRequest{
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
		Use:   "extendedPairVault-mapping-by-App [appID]",
		Short: "Query ExtendedPair Vault Mapping By App",
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

			res, err := queryClient.QueryExtendedPairVaultMappingByApp(cmd.Context(), &types.QueryExtendedPairVaultMappingByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "extendedPairVault-mapping-by-App")
	return cmd
}

func QueryExtendedPairVaultMappingByAppAndExtendedPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extendedPairVault-mapping-by-App-and-ExtendedPair [appID] [extendedPairID]",
		Short: "Query ExtendedPair Vault Mapping By App And ExtendedPair",
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

			res, err := queryClient.QueryExtendedPairVaultMappingByAppAndExtendedPair(cmd.Context(), &types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest{
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

func QueryTVLByAppOfAllExtendedPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tvl-by-app-all-extended-pairs [appID]",
		Short: "Query tvl By App of all extended pairs",
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

			res, err := queryClient.QueryTVLByAppOfAllExtendedPairs(cmd.Context(), &types.QueryTVLByAppOfAllExtendedPairsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "tvl-by-app-all-extended-pairs")
	return cmd
}

func QueryTVLByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tvl-by-app [appID]",
		Short: "Query total tvl by an App",
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

			res, err := queryClient.QueryTVLByApp(cmd.Context(), &types.QueryTVLByAppRequest{
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
		Short: "Query user my position by app",
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
		Use:   "user-extended-pair-total-data [owner]",
		Short: "Query user Extended Pair Total Data",
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

func QueryPairsLockedAndMintedStatisticByApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pairs-locked-and-minted-statistic-by-app [appID]",
		Short: "Query pairs locked and minted statistic by app",
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

			res, err := queryClient.QueryPairsLockedAndMintedStatisticByApp(cmd.Context(), &types.QueryPairsLockedAndMintedStatisticByAppRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "pairs-locked-and-minted-statistic-by-app")
	return cmd
}
