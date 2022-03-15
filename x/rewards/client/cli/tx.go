package cli

import (
	"strconv"
	"time"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
)

func AddNewMintingRewardsProposalCLIHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-new-mint-rewards [collateral-denom] [casset-denom] [total-rewards] [casset-maxcap] [duration-days]",
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

			cAssetDenom := args[1]
			err = sdk.ValidateDenom(cAssetDenom)
			if err != nil {
				return err
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

			content := types.AddNewMintRewardsProposalContent(
				title,
				description,
				collateralDenom,
				cAssetDenom,
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

func txDepositMintingRewardAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [minting-reward-id] [start-timestamp]",
		Short: "Deposit amount for the rewards added through governance proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mintingRewardId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			startTime, err := time.Parse("2006-01-02 15:04:05", args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositMintingRewardAmount(mintingRewardId, clientCtx.GetFromAddress(), startTime)
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

func txUpdateMintingRewardStartTimestamp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [minting-reward-id] [start-timestamp]",
		Short: "update existing minging reward",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mintingRewardId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			newStartTimestamp, err := time.Parse("2006-01-02 15:04:05", args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateMintRewardStartTime(mintingRewardId, clientCtx.GetFromAddress(), newStartTimestamp)
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

func DisableMintingRewardsProposalCLIHandler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-mint-rewards [minting-reward-id]",
		Short: "add new mint rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mintingRewardId, err := strconv.ParseUint(args[0], 10, 64)
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

			content := types.DisableMintRewardsProposalContent(
				title,
				description,
				mintingRewardId,
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
