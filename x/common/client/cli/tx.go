package cli

import (
	"fmt"
	"time"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/comdex-official/comdex/x/common/types"
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

	cmd.AddCommand(CmdRegisterContract())

	// this line is used by starport scaffolding # 1

	return cmd 
}

func CmdRegisterContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-contract [game name] [game id] [contract address]",
		Short: "Register game contract",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			gameId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("game-id '%s' not a valid uint", args[1])
			}

			msg := types.NewMsgRegisterContract(
				clientCtx.GetFromAddress().String(),
				args[0],
				gameId,
				args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
