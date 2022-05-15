package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	// "strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/locking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group locking queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryLockByIdCmd(),
		NewQueryLockByOwnerCmd(),
		NewQueryUnlockingByIdCmd(),
		NewQueryUnlockingByOwnerCmd(),
		NewQueryAllLocksCmd(),
		NewQueryAllUnlockingsCmd(),
	)
	return cmd
}

// NewQueryParamsCmd implements the params query command.
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current locking parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as locking parameters.
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

func NewQueryLockByIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-by-id [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query lock by id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query lock by id.
Example:
$ %s query %s lock-by-id 1
`,
				version.AppName, types.ModuleName,
			),
		),
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
			res, err := queryClient.QueryLockById(
				context.Background(),
				&types.QueryLockByIdRequest{
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

func NewQueryLockByOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locks-by-owner [owner]",
		Args:  cobra.ExactArgs(1),
		Short: "Query locks by owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query locks by owner.
Example:
$ %s query %s locks-by-owner comdex1d0qjudc6wyqv35z3k9n4kyj2mquyc2647gs7nr
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryLocksByOwner(
				context.Background(),
				&types.QueryLocksByOwnerRequest{
					Owner: owner.String(),
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

func NewQueryAllLocksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locks",
		Args:  cobra.NoArgs,
		Short: "Query locks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all available locks.
Example:
$ %s query %s locks
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
			res, err := queryClient.QueryAllLocks(
				context.Background(),
				&types.QueryAllLocksRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "locks")
	return cmd
}

func NewQueryUnlockingByIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlocking-by-id [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query unlocking by id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unlocking by id.
Example:
$ %s query %s unlocking-by-id 1
`,
				version.AppName, types.ModuleName,
			),
		),
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
			res, err := queryClient.QueryUnlockingById(
				context.Background(),
				&types.QueryUnlockingByIdRequest{
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

func NewQueryUnlockingByOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlockings-by-owner [owner]",
		Args:  cobra.ExactArgs(1),
		Short: "Query unlockings by owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unlockings by owner.
Example:
$ %s query %s unlockings-by-owner comdex1d0qjudc6wyqv35z3k9n4kyj2mquyc2647gs7nr
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(ctx)
			res, err := queryClient.QueryUnlockingsByOwner(
				context.Background(),
				&types.QueryUnlockingsByOwnerRequest{
					Owner: owner.String(),
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

func NewQueryAllUnlockingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlockings",
		Args:  cobra.NoArgs,
		Short: "Query unlockings",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all available unlockings.
Example:
$ %s query %s unlockings
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
			res, err := queryClient.QueryAllUnlockings(
				context.Background(),
				&types.QueryAllUnlockingsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "unlockings")
	return cmd
}
