package cli

import (
	"strconv"
	"strings"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
)

func AddNewMintingRewardsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-new-mint-rewards [collateral-denom] [casset-denoms] [total-rewards] [casset-maxcap] [duration-days]",
		Short: "add new mint rewards",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralDenom := args[0]
			err = sdk.ValidateDenom(collateralDenom)
			if err != nil {
				return err
			}

			var cAssetDenoms []string
			for _, denom := range strings.Split(args[1], ",") {
				denom = strings.TrimSpace(denom)
				err = sdk.ValidateDenom(denom)
				if err != nil {
					return err
				}
				cAssetDenoms = append(cAssetDenoms, denom)
			}

			totalRewards, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			cAssetMaxcap, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			durationDays, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewMintRewardsProposal(
				title,
				description,
				collateralDenom,
				cAssetDenoms,
				totalRewards,
				cAssetMaxcap,
				durationDays,
			)

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
