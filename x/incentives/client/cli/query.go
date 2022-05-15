package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/comdex-official/comdex/x/incentives/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group incentives queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryEpochInfoByDurationCmd(),
		NewQueryAllEpochsInfoCmd(),
		NewQueryAllGaugesCmd(),
	)
	return cmd
}

// NewQueryParamsCmd implements the params query command.
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current incentives parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as incentives parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryEpochInfoByDurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epoch-by-duration [seconds]",
		Args:  cobra.ExactArgs(1),
		Short: "Query epoch info by duration",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query epoch by duration.
Example:
$ %s query %s epoch-by-duration 24h
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			seconds, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gauge-type-id: %w", err)
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryEpochInfoByDuration(
				context.Background(),
				&types.QueryEpochInfoByDurationRequest{
					DurationSeconds: seconds,
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

func NewQueryAllEpochsInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epochs",
		Args:  cobra.NoArgs,
		Short: "Query all epochs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all available epochs.
Example:
$ %s query %s gauges
`,
				version.AppName, types.ModuleName,
			),
		),
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
			res, err := queryClient.QueryAllEpochsInfo(
				context.Background(),
				&types.QueryAllEpochsInfoRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "epochs")
	return cmd
}

func NewQueryAllGaugesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gauges",
		Args:  cobra.NoArgs,
		Short: "Query all gauges",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all available gauges.
Example:
$ %s query %s gauges
`,
				version.AppName, types.ModuleName,
			),
		),
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
			res, err := queryClient.QueryAllGauges(
				context.Background(),
				&types.QueryAllGaugesRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "gauges")
	return cmd
}
