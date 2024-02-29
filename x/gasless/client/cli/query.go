package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/comdex-official/comdex/x/gasless/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gasless",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "gasless"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryMessagesAndContractsCmd(),
		NewQueryGasProviderCmd(),
		NewQueryGasProvidersCmd(),
	)

	return cmd
}

// NewQueryParamsCmd implements the params query command.
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current gasless module's parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as gasless module's parameters.
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

// NewQueryMessagesAndContractsCmd implements the messages and contracts query command.
func NewQueryMessagesAndContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mac",
		Args:  cobra.NoArgs,
		Short: "Query all the available messages and contract addresses",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the available messages and contract addresses.
Example:
$ %s query %s mac
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

			resp, err := queryClient.MessagesAndContracts(
				cmd.Context(),
				&types.QueryMessagesAndContractsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasprovider [gas-provider-id]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of the gas provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the gas provider
Example:
$ %s query %s gasprovider
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasProviderId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas_provider_id: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasProvider(
				cmd.Context(),
				&types.QueryGasProviderRequest{
					GasProviderId: gasProviderId,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasproviders ",
		Args:  cobra.MinimumNArgs(0),
		Short: "Query details of all the gas providers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of all the gas providers
Example:
$ %s query %s gasproviders
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasProviders(
				cmd.Context(),
				&types.QueryGasProvidersRequest{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
