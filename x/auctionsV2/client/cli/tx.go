package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
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
		Use:                        "auctions",
		Short:                      "AuctionsV2 module sub-commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txPlaceMarketDutchBid(),
		txDepositLimitDutchBid(),
		txCancelLimitDutchBid(),
		txWithdrawLimitDutchBid(),
	)

	return cmd
}

func txPlaceMarketDutchBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market-bid-request [auction-id] [bid-amount]",
		Short: "Place a market bid on a dutch auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			amt, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceMarketBid(clientCtx.GetFromAddress().String(), id, amt)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txDepositLimitDutchBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-limit-bid-request [collateral-token-id] [debt-token-id] [discount] [bid-amount]",
		Short: "Place a limit bid on a dutch auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralTokenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("collateralTokenID '%s' not a valid uint", args[0])
			}
			debtTokenID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("debtTokenID '%s' not a valid uint", args[1])
			}

			premiumDiscount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("premiumDiscount '%s' not a valid int", args[2])
			}

			amt, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositLimitBid(clientCtx.GetFromAddress().String(), collateralTokenID, debtTokenID, premiumDiscount, amt)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txCancelLimitDutchBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-limit-bid-request [collateral-token-id] [debt-token-id] [discount]",
		Short: "Cancel a limit bid on a dutch auction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralTokenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("collateralTokenID '%s' not a valid uint", args[0])
			}
			debtTokenID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("debtTokenID '%s' not a valid uint", args[1])
			}

			premiumDiscount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("premiumDiscount '%s' not a valid int", args[2])
			}

			msg := types.NewMsgCancelLimitBid(clientCtx.GetFromAddress().String(), collateralTokenID, debtTokenID, premiumDiscount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txWithdrawLimitDutchBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-limit-bid-request [collateral-token-id] [debt-token-id] [discount] [bid-amount]",
		Short: "Withdraw a limit bid on a dutch auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralTokenID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("collateralTokenID '%s' not a valid uint", args[0])
			}
			debtTokenID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("debtTokenID '%s' not a valid uint", args[1])
			}

			premiumDiscount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("premiumDiscount '%s' not a valid int", args[2])
			}

			amt, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawLimitBid(clientCtx.GetFromAddress().String(), collateralTokenID, debtTokenID, premiumDiscount, amt)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
