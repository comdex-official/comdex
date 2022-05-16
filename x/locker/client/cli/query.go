package cli

import (
	"context"

	"github.com/comdex-official/comdex/x/locker/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func queryLockedVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockerinfo [id]",
		Short: "Query locker info by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id := args[0]
			queryClient := types.NewQueryServiceClient(ctx)
			res, err := queryClient.QueryLockerInfo(
				context.Background(),
				&types.QueryLockerInfoRequest{
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