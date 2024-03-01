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
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		NewQueryGasConsumerCmd(),
		NewQueryGasConsumersCmd(),
		NewQueryTxGpidsCmd(),
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

			gasProviderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas_provider_id: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasProvider(
				cmd.Context(),
				&types.QueryGasProviderRequest{
					GasProviderId: gasProviderID,
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

func NewQueryGasConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasconsumer [consumer]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of the gas consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the gas consumer
Example:
$ %s query %s gasconsumer
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasConsumer(
				cmd.Context(),
				&types.QueryGasConsumerRequest{
					Consumer: sanitizedConsumer.String(),
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

func NewQueryGasConsumersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasconsumers ",
		Args:  cobra.MinimumNArgs(0),
		Short: "Query details of all the gas consumers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of all the gas consumers
Example:
$ %s query %s gasconsumers
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
			resp, err := queryClient.GasConsumers(
				cmd.Context(),
				&types.QueryGasConsumersRequest{
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

// NewQueryTxGpidsCmd implements the tx-gpids query command.
func NewQueryTxGpidsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx-gpids",
		Args:  cobra.NoArgs,
		Short: "Query all the tx type url and contract address along with associcated gas provider ids",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the tx type url and contract address along with associcated gas provider ids
Example:
$ %s query %s tx-gpids
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

			resp, err := queryClient.GasProviderIdsForAllTXC(cmd.Context(), &types.QueryGasProviderIdsForAllTXC{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
