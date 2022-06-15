package cli

import (
	"fmt"
	"strconv"

	"github.com/comdex-official/comdex/x/auction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func txPlaceSurplusBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid-surplus [auction-id] [bid] [app-id] [auction-mapping-id]",
		Short: "Place a bid on an auction",
		Args:  cobra.ExactArgs(4),
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

			appId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			auctionMappingId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			msg := types.NewMsgPlaceSurplusBid(clientCtx.GetFromAddress(), id, amt, appId, auctionMappingId)
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

func txPlaceDebtBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid-debt [auction-id] [bid] [user expected Token] [app-id] [auction-mapping-id] ",
		Short: "Place a Debt bid on an auction",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			bid, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			userExpectedToken, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			appId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			auctionMappingId, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			msg := types.NewMsgPlaceDebtBid(clientCtx.GetFromAddress(), id, bid, userExpectedToken, appId, auctionMappingId)
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

func txPlaceDutchBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid-dutch [auction-id] [amount] [app-id] [auction-mapping-id]",
		Short: "Place a Dutch bid on an auction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			amt, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			// max := sdk.MustNewDecFromStr(args[1])

			appId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			auctionMappingId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			msg := types.NewMsgPlaceDutchBid(clientCtx.GetFromAddress(), auctionId, amt, appId, auctionMappingId)
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
