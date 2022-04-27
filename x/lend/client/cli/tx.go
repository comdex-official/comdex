package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/lend/types"
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

	cmd.AddCommand(
		txLend(),
		txWithdraw(), //withdraw collateral partially or fully
		txBorrowAsset(),
		txRepayAsset(), //including functionality of both repaying and closing position
		txFundModuleAccounts(),
	)

	return cmd
}

func txLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend [amount]",
		Short: "lend a whitelisted asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLend(ctx.GetFromAddress(), asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [amount]",
		Short: "withdraw lent asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(ctx.GetFromAddress(), asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txBorrowAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow [amount]",
		Short: "borrow a whitelisted asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBorrow(ctx.GetFromAddress(), asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRepayAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repay [amount]",
		Short: "repay borrowed asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRepay(ctx.GetFromAddress(), asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func NewCmdSubmitAddWhitelistedAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-whitelisted-assets [name] [Denom] [Decimals] [Collateral_Weight] [Liquidation_Threshold] [Base_Borrow_Rate] [Base_Lend_Rate]",
		Args:  cobra.ExactArgs(7),
		Short: "Add whitelisted assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			names, err := ParseStringFromString(args[0], ",")
			if err != nil {
				return err
			}
			denoms, err := ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}

			decimals, err := ParseInt64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			collateralWeight, err := ParseStringFromString(args[3], ",")

			liquidationThreshold, err := ParseStringFromString(args[4], ",")

			baseBorrowRate, err := ParseStringFromString(args[5], ",")

			baseLendRate, err := ParseStringFromString(args[6], ",")

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := ctx.GetFromAddress()

			var assets []types.Asset
			for i, _ := range names {
				newcollateralWeigt, _ := sdk.NewDecFromStr(collateralWeight[i])
				newliquidationThreshold, _ := sdk.NewDecFromStr(liquidationThreshold[i])
				newbaseBorrowRate, _ := sdk.NewDecFromStr(baseBorrowRate[i])
				newbaseLendRate, _ := sdk.NewDecFromStr(baseLendRate[i])

				assets = append(assets, types.Asset{
					Name:                 names[i],
					Denom:                denoms[i],
					Decimals:             decimals[i],
					CollateralWeight:     newcollateralWeigt,
					LiquidationThreshold: newliquidationThreshold,
					BaseBorrowRate:       newbaseBorrowRate,
					BaseLendRate:         newbaseLendRate,
				})
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewAddWhitelistedAssetsProposal(title, description, assets)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func txFundModuleAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-module [module-name] [amount]",
		Short: "Deposit amount to the respective module account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moduleName := args[0]

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgFundModuleAccounts(moduleName, ctx.GetFromAddress(), amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
