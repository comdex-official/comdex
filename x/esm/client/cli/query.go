package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/petrichormoney/petri/x/esm/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group esm queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "esm",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "esm"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		queryESMTriggerParams(),
		queryESMStatus(),
		queryCurrentDepositStats(),
		queryUsersDepositMapping(),
		queryDataAfterCoolOff(),
		querySnapshotPrice(),
		queryAssetDataAfterCoolOff(),
	)

	return cmd
}

func queryESMTriggerParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "esm-trigger-params [app-id]",
		Short: "Query a esm trigger params",
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

			res, err := queryClient.QueryESMTriggerParams(
				context.Background(),
				&types.QueryESMTriggerParamsRequest{
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

func queryESMStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "esm-status [app-id]",
		Short: "Query esm status by app",
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

			res, err := queryClient.QueryESMStatus(
				context.Background(),
				&types.QueryESMStatusRequest{
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

func queryCurrentDepositStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-deposit-stats [app-id]",
		Short: "Query current deposit stats",
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

			res, err := queryClient.QueryCurrentDepositStats(
				context.Background(),
				&types.QueryCurrentDepositStatsRequest{
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

func queryUsersDepositMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-deposits [app-id] [depositor]",
		Short: "Query user deposits for esm",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			depositor := args[1]

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryUsersDepositMapping(
				context.Background(),
				&types.QueryUsersDepositMappingRequest{
					Id:        id,
					Depositor: depositor,
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

func queryDataAfterCoolOff() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data_after_cool_off [app-id]",
		Short: "Query data after cool off period for esm",
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

			res, err := queryClient.QueryDataAfterCoolOff(
				context.Background(),
				&types.QueryDataAfterCoolOffRequest{
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

func querySnapshotPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price_snapshot [app-id] [asset-id]",
		Short: "Query price snapshot for esm",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QuerySnapshotPrice(
				context.Background(),
				&types.QuerySnapshotPriceRequest{
					AppId:   id,
					AssetId: asset,
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

func queryAssetDataAfterCoolOff() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset_data_after_cool_off [app-id]",
		Short: "Query asset data after cool off period for esm",
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

			res, err := queryClient.QueryAssetDataAfterCoolOff(
				context.Background(),
				&types.QueryAssetDataAfterCoolOffRequest{
					AppId: id,
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
