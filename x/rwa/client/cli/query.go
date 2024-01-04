package cli

import (
	"context"
	"github.com/comdex-official/comdex/x/rwa/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"strconv"
)

func queryRwaUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rwa-user [account_address]",
		Short: "Query an user data by user id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			acct := args[0]
			if len(acct) < 0 {
				return types.ErrAccountAddressEmpty
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryRwaUser(
				context.Background(),
				&types.RwaUserRequest{
					AccountAddress: acct,
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

func queryCounterParty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counter-party [id]",
		Short: "Query a counterparty data by id",
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

			res, err := queryClient.QueryCounterParty(
				context.Background(),
				&types.CounterPartyRequest{
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

func queryInvoice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice [id]",
		Short: "Query a invoice data by id",
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

			res, err := queryClient.QueryInvoice(
				context.Background(),
				&types.InvoiceRequest{
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

func queryInvoiceAsReceiver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice-receiver [account_address]",
		Short: "Query an invoice as sender",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			acct := args[0]
			if len(acct) < 0 {
				return types.ErrAccountAddressEmpty
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryInvoiceAsReceiver(
				context.Background(),
				&types.InvoiceReceiverRequest{
					Receiver: acct,
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

func queryInvoiceAsSender() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice-sender [account_address]",
		Short: "Query an invoice as sender",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			acct := args[0]
			if len(acct) < 0 {
				return types.ErrAccountAddressEmpty
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryInvoiceAsSender(
				context.Background(),
				&types.InvoiceSenderRequest{
					From: acct,
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
