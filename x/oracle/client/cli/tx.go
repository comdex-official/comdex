package cli

import (
	"github.com/comdex-official/comdex/x/oracle/types"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	ibcchannelclientutils "github.com/cosmos/ibc-go/modules/core/04-channel/client/utils"
	"github.com/spf13/cobra"
)

func txAddMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-market [symbol] [script-id]",
		Short: "Add a market",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scriptID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddMarketRequest(
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

func txAddMarketForAsset() *cobra.Command {
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
}

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

func txFetchPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch-price [source-port] [source-channel] [symbols] [script-id]",
		Short: "Fetch price from Oracle",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scriptID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			timeoutHeight, err := GetPacketTimeoutHeight(cmd)
			if err != nil {
				return err
			}

			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}

			absoluteTimeouts, err := cmd.Flags().GetBool(flagAbsoluteTimeouts)
			if err != nil {
				return err
			}

			feeLimit, err := GetFeeLimit(cmd)
			if err != nil {
				return err
			}

			prepareGas, err := cmd.Flags().GetUint64(flagPrepareGas)
			if err != nil {
				return err
			}

			executeGas, err := cmd.Flags().GetUint64(flagExecuteGas)
			if err != nil {
				return err
			}

			if !absoluteTimeouts {
				state, height, _, err := ibcchannelclientutils.QueryLatestConsensusState(ctx, args[0], args[1])
				if err != nil {
					return err
				}

				if !timeoutHeight.IsZero() {
					timeoutHeight.RevisionHeight += height.RevisionHeight
					timeoutHeight.RevisionNumber += height.RevisionNumber
				}
				if timeoutTimestamp != 0 {
					timeoutTimestamp += state.GetTimestamp()
				}
			}

			msg := types.NewMsgFetchPriceRequest(
				ctx.FromAddress,
				args[0],
				args[1],
				timeoutHeight,
				timeoutTimestamp,
				strings.Split(args[2], ","),
				scriptID,
				feeLimit,
				prepareGas,
				executeGas,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagPacketTimeoutHeight, "0-1000", "packet timeout block height")
	cmd.Flags().Duration(flagPacketTimeoutTimestamp, 10*time.Minute, "packet timeout timestamp")
	cmd.Flags().Bool(flagAbsoluteTimeouts, false, "timeout flags are used as absolute timeouts")
	cmd.Flags().String(flagFeeLimit, "", "fee limit")
	cmd.Flags().Uint64(flagPrepareGas, 0, "prepare gas")
	cmd.Flags().Uint64(flagExecuteGas, 0, "execute gas")

	return cmd
}
