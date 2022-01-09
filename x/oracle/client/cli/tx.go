package cli

import (
	"github.com/comdex-official/comdex/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"
)

func txAddMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-market [symbol] [script-id] [asset]",
		Short: "Add a market",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scriptID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			asset, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			var rates uint64
			rates = 0

			msg := types.NewMsgAddMarketRequest(
				ctx.FromAddress,
				args[0],
				scriptID,
				asset,
				rates,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func txUpdateMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-market [symbol]",
		Short: "Update a market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scriptID, err := cmd.Flags().GetUint64(flagScriptID)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateMarketRequest(
				ctx.FromAddress,
				args[0],
				scriptID,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Uint64(flagScriptID, 0, "script identity")

	return cmd
}

/*func txAddMarketForAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-market-for-asset [asset] [symbol]",
		Short: "Add a market for asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddMarketForAssetRequest(
				ctx.FromAddress,
				asset,
				args[1],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}*/

func txRemoveMarketForAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-market-for-asset [asset] [symbol]",
		Short: "Remove a market for asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveMarketForAssetRequest(
				ctx.FromAddress,
				asset,
				args[1],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
