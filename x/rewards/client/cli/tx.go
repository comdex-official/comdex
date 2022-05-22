package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/rewards/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		txWhitelistAsset(),
		txRemoveWhitelistAsset(),
		txWhitelistAppIdVault(),
		txRemoveWhitelistAppIdVault(),
	)

	return cmd
}

func txWhitelistAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-asset [app_mapping_Id] [asset_Id]",
		Short: "Add Whitelisted assetId for Locker savings rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_mapping_Id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_Id, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			var newAssetId []uint64
			for i := range asset_Id {
				newAssetId = append(newAssetId, asset_Id[i])
			}

			msg := types.NewMsgWhitelistAsset(app_mapping_Id, ctx.GetFromAddress(), newAssetId)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRemoveWhitelistAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-asset [app_mapping_Id] [asset_Id]",
		Short: "Remove Whitelisted assetId for Locker savings rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_mapping_Id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_Id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveWhitelistAsset(app_mapping_Id, ctx.GetFromAddress(), asset_Id)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txWhitelistAppIdVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-app-id-vault-interest [app_mapping_Id]",
		Short: "na",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgWhitelistAppIdVault(
				appMappingId,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRemoveWhitelistAppIdVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-app-id-vault-interest [app_mapping_Id] ",
		Short: "na",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveWhitelistAppIdVault(
				appMappingId,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}
