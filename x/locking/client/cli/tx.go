package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/locking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
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

	cmd.AddCommand(
		NewLockTokensCmd(),
		NewBeginUnlockByIDCmd(),
	)
	return cmd
}

// NewLockTokensCmd lock tokens into bonding pool from user's account.
func NewLockTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-tokens [tokens]",
		Short: "lock tokens into bonding pool from user account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			durationStr, err := cmd.Flags().GetString(FlagDuration)
			if err != nil {
				return err
			}

			duration, err := time.ParseDuration(durationStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgLockTokens(
				clientCtx.GetFromAddress(),
				duration,
				coin,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetLockTokens())
	flags.AddTxFlagsToCmd(cmd)
	err := cmd.MarkFlagRequired(FlagDuration)
	if err != nil {
		panic(err)
	}
	return cmd
}

// NewBeginUnlockByIDCmd unlock individual period lock by ID.
func NewBeginUnlockByIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "begin-unlock [lock-id] [tokens]",
		Short: "begin unlock individual period lock by ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse lock-id id: %w", err)
			}

			unlockingCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBeginUnlockingTokens(
				clientCtx.GetFromAddress(),
				uint64(id),
				unlockingCoin,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
