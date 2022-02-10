package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/comdex-official/comdex/x/poolapi/types"
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
		IndividualPoolLiquidity(),
		PoolsLiquidity(),
		TotalCollateral(),
		PoolAPR(),
	)

	return cmd
}

func IndividualPoolLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [id]",
		Short: "liquidity of the given pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.IndividualPoolLiquidity(cmd.Context(), &types.QueryIndividualPoolLiquidityRequest{
				PoolId: id,
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

func PoolsLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "total liquidity of all pools",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.PoolsLiquidity(cmd.Context(), &types.QueryTotalLiquidityRequest{})

			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func TotalCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral",
		Short: "total collateral in all vaults",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.TotalCollateral(cmd.Context(), &types.QueryTotalCollateralRequest{})

			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func PoolAPR() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apr",
		Short: "apr of pools",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(ctx)

			res, err := queryClient.PoolAPR(cmd.Context(), &types.QueryPoolAPRRequest{})

			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
