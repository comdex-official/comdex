package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strconv"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidationsV2/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group liquidationsV2 queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		queryLockedVault(),
		queryLockedVaults(),
		queryLiquidationWhitelisting(),
		queryLiquidationWhitelistings(),
		queryLockedVaultsHistory(),
		queryAppReserveFundsTxData(),
	)

	return cmd
}

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

func queryLiquidationWhitelisting() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidation-whitelisting [app-id]",
		Short: "Query liquidation-whitelisting",
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
			res, err := queryClient.QueryLiquidationWhiteListing(
				context.Background(),
				&types.QueryLiquidationWhiteListingRequest{
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

func queryLiquidationWhitelistings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidation-whitelistings",
		Short: "Query liquidation-whitelistings",
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
			res, err := queryClient.QueryLiquidationWhiteListings(
				context.Background(),
				&types.QueryLiquidationWhiteListingsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "liquidation-whitelistings")

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

func queryAppReserveFundsTxData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app-reserve-funds [app-id]",
		Short: "Query app-reserve-funds",
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
			res, err := queryClient.QueryAppReserveFundsTxData(
				context.Background(),
				&types.QueryAppReserveFundsTxDataRequest{
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
