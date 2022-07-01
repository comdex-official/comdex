package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	// Group rewards queries under a subcommand
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
		NewQueryGaugeByIDCmd(),
		NewQueryGaugeByDurationCmd(),
		queryReward(),
		queryRewards(),
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

// NewQueryEpochInfoByDurationCmd implements the epoch-by-duration query command.
func NewQueryEpochInfoByDurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epoch-by-duration [seconds]",
		Args:  cobra.ExactArgs(1),
		Short: "Query epoch info by duration",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query epoch by duration.
Example:
$ %s query %s epoch-by-duration 600
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
				return fmt.Errorf("parse seconds: %w", err)
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

// NewQueryAllEpochsInfoCmd implements the epochs query command.
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

// NewQueryAllGaugesCmd implements the gauges query command.
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

// NewQueryGaugeByIDCmd implements the gauge query command.
func NewQueryGaugeByIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gauge [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query gauge by id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query gauge by id.
Example:
$ %s query %s gauge 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			gaugeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse id: %w", err)
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryGaugeByID(
				context.Background(),
				&types.QueryGaugeByIdRequest{
					GaugeId: gaugeID,
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

// NewQueryGaugeByDurationCmd implements the gauges-by-duration query command.
func NewQueryGaugeByDurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gauges-by-duration [seconds]",
		Args:  cobra.ExactArgs(1),
		Short: "Query gauges by duration",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query gauge by duration.
Example:
$ %s query %s gauges-by-duration 600
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
				return fmt.Errorf("parse seconds: %w", err)
			}

			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryGaugeByDuration(
				context.Background(),
				&types.QueryGaugesByDurationRequest{
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

func queryReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "internal-reward [id]",
		Short: "Query reward",
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

			res, err := queryClient.QueryReward(
				context.Background(),
				&types.QueryRewardRequest{
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

func queryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "internal-rewards",
		Short: "Query internal-rewards",
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

			res, err := queryClient.QueryRewards(
				context.Background(),
				&types.QueryRewardsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "lends")

	return cmd
}
