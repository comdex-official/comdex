package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/esm/types"
)

func NewCmdSubmitToggleEsmProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle-esm [appId] [vaultStop] [lockerStop] [mintingStop]",
		Args:  cobra.ExactArgs(4),
		Short: "toggle emergency shutdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			vaultStop := ParseBoolFromString(args[1])
			lockerStop := ParseBoolFromString(args[2])
			mintingStop := ParseBoolFromString(args[3])

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			var esmActive types.EsmActive
			esmActive.AppId = appId
			esmActive.VaultStop = vaultStop
			esmActive.LockerStop = lockerStop
			esmActive.MintingStop = mintingStop
			

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewToggleEsmProposal(title, description, esmActive)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}
