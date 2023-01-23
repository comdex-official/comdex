package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/liquidation/types"
)

func queryLockedVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vault [app-id] [id]",
		Short: "Query locked-vault",
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
			id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLockedVault(
				context.Background(),
				&types.QueryLockedVaultRequest{
					AppId: appID,
					Id:    id,
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

func queryLockedVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vaults",
		Short: "Query locked-vaults",
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
			res, err := queryClient.QueryLockedVaults(
				context.Background(),
				&types.QueryLockedVaultsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "locked-vaults")

	return cmd
}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query module parameters",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryParams(
				context.Background(),
				&types.QueryParamsRequest{},
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

func queryLockedVaultsHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vaults-history",
		Short: "Query locked-vaults history",
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
			res, err := queryClient.QueryLockedVaultsHistory(
				context.Background(),
				&types.QueryLockedVaultsHistoryRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "locked-vaults-history")

	return cmd
}

func queryUserLockedVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vaults-by-user [user_address]",
		Short: "locked vaults list for an individual account",
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

			res, err := queryClient.QueryUserLockedVaults(cmd.Context(), &types.QueryUserLockedVaultsRequest{
				UserAddress: args[0],
				Pagination:  pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "locked-vaults-by-user")

	return cmd
}

func queryUserLockedVaultsHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vaults-history-by-user [user_address]",
		Short: "historical locked vaults list for an individual account",
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

			res, err := queryClient.QueryUserLockedVaultsHistory(cmd.Context(), &types.QueryUserLockedVaultsHistoryRequest{
				UserAddress: args[0],
				Pagination:  pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "locked-vaults-history-by-user")

	return cmd
}

func queryLockedVaultsPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vaults-pair [pair_id]",
		Short: "locked vaults list With Pair Id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			pairID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryLockedVaultsPair(cmd.Context(), &types.QueryLockedVaultsPairRequest{
				PairId:     pairID,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "locked-vaults-pair")

	return cmd
}

func queryAppIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted-app-id",
		Short: "Query whitelisted app id",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryAppIds(cmd.Context(), &types.QueryAppIdsRequest{})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
