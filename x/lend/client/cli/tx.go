package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
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
		txDeposit(),
		txCloseLend(),
		txBorrowAsset(),
		txDrawAsset(),
		txRepayAsset(), //including functionality of both repaying and closing position
		txDepositBorrowAsset(),
		txFundModuleAccounts(),
	)

	return cmd
}

func txLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend [Asset_Id] [Amount] [Pool_Id]",
		Short: "lend a whitelisted asset",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pair, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			pool, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgLend(ctx.GetFromAddress(), pair, asset, pool)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [lendId] [Amount]",
		Short: "withdraw lent asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(ctx.GetFromAddress(), lendId, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [lendId] [Amount]",
		Short: "deposit into a lent position",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(ctx.GetFromAddress(), lendId, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txCloseLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-lend [lendId]",
		Short: "close a lent position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCloseLend(ctx.GetFromAddress(), lendId)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txBorrowAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow [lend-id] [pair-id] [is-stable-borrow] [amount-in] [amount-out]",
		Short: "borrow a whitelisted asset",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			lendId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pairId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			StableBorrow, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			isStableBorrow := ParseBoolFromString(StableBorrow)

			amountIn, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			amountOut, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgBorrow(ctx.GetFromAddress(), lendId, pairId, isStableBorrow, amountIn, amountOut)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRepayAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repay [borrow-id] [amount]",
		Short: "repay borrowed asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRepay(ctx.GetFromAddress(), borrowId, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txDrawAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw [borrow-id] [amount]",
		Short: "draw borrowed asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDraw(ctx.GetFromAddress(), borrowId, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txDepositBorrowAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-borrow [borrow-id] [amount]",
		Short: "deposit borrowed asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositBorrow(ctx.GetFromAddress(), borrowId, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txFundModuleAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-module [module-name] [asset_id] [amount]",
		Short: "Deposit amount to the respective module account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moduleName := args[0]

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgFundModuleAccounts(moduleName, assetId, ctx.GetFromAddress(), amount)
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

func CmdAddWNewLendPairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lend-pairs [asset_in] [asset_out] [is_inter_pool] [asset_out_pool_id] [liquidation_ratio]",
		Short: "Add lend asset pairs",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetIn, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			assetOut, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			isInterPool, err := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			assetOutPoolId, err := ParseUint64SliceFromString(args[3], ",")
			if err != nil {
				return err
			}

			var pairs []types.Extended_Pair
			for i := range assetIn {
				interPool := ParseBoolFromString(isInterPool[i])
				pairs = append(pairs, types.Extended_Pair{
					AssetIn:        assetIn[i],
					AssetOut:       assetOut[i],
					IsInterPool:    interPool,
					AssetOutPoolId: assetOutPoolId[i],
				})
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

			content := types.NewAddLendPairsProposal(title, description, pairs)

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

func CmdUpdateLendPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-lend-pair [len_pair_id]",
		Short: "Update a lend asset pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pair := types.Extended_Pair{
				Id: id,
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

			content := types.NewUpdateLendPairProposal(title, description, pair)

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

func CmdAddPoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lend-pool [module_name] [first_bridged_asset_id] [second_bridged_asset_id] [asset_id] [is_bridged_asset] ",
		Short: "Add lend pool ",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moduleName := args[0]

			firstBridgedAssetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			secondBridgedAssetId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := ParseUint64SliceFromString(args[3], ",")
			if err != nil {
				return err
			}

			isBridgedAsset, err := ParseUint64SliceFromString(args[4], ",")
			if err != nil {
				return err
			}
			var pool types.Pool
			var assetData []types.AssetDataPoolMapping

			for i := range assetId {
				bridged := ParseBoolFromString(isBridgedAsset[i])
				assetData = append(assetData, types.AssetDataPoolMapping{
					AssetId:   assetId[i],
					IsBridged: bridged,
				})
			}
			pool = types.Pool{
				ModuleName:           moduleName,
				FirstBridgedAssetId:  firstBridgedAssetId,
				SecondBridgedAssetId: secondBridgedAssetId,
				AssetData:            assetData,
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

			content := types.NewAddPoolProposal(title, description, pool)

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

func CmdAddAssetToPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-to-pair-mapping [asset_id] [pair_id] ",
		Short: "Add asset to pair mapping ",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			rawPairId, _ := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}
			var pairId []uint64
			for i := range rawPairId {

				pairId = append(pairId, rawPairId[i])
			}
			assetToPairMapping := types.AssetToPairMapping{
				AssetId: assetId,
				PairId:  pairId,
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

			content := types.NewAddAssetToPairProposal(title, description, assetToPairMapping)

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

func CmdAddWNewAssetRatesStatsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-rates-stats [asset_id] [u_optimal] [base] [slope1] [slope2] [stable_base] [stable_slope1] [stable_slope2] [ltv] [liquidation_threshold] [liquidation_penalty] [reserve_factor]",
		Short: "Add lend asset pairs",
		Args:  cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetId, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			uOptimal, err := ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}
			base, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}
			slope1, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}
			slope2, err := ParseStringFromString(args[4], ",")
			if err != nil {
				return err
			}
			enableStableBorrow, err := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}
			stableBase, err := ParseStringFromString(args[5], ",")
			if err != nil {
				return err
			}
			stableSlope1, err := ParseStringFromString(args[6], ",")
			if err != nil {
				return err
			}
			stableSlope2, err := ParseStringFromString(args[7], ",")
			if err != nil {
				return err
			}
			ltv, err := ParseStringFromString(args[8], ",")
			if err != nil {
				return err
			}
			liquidationThreshold, err := ParseStringFromString(args[9], ",")
			if err != nil {
				return err
			}
			liquidationPenalty, err := ParseStringFromString(args[10], ",")
			if err != nil {
				return err
			}
			reserveFactor, err := ParseStringFromString(args[11], ",")
			if err != nil {
				return err
			}

			var assetRatesStats []types.AssetRatesStats
			for i := range assetId {
				newUOptimal, _ := sdk.NewDecFromStr(uOptimal[i])
				newBase, _ := sdk.NewDecFromStr(base[i])
				newSlope1, _ := sdk.NewDecFromStr(slope1[i])
				newSlope2, _ := sdk.NewDecFromStr(slope2[i])
				newEnableStableBorrow := ParseBoolFromString(enableStableBorrow[i])
				newStableBase, _ := sdk.NewDecFromStr(stableBase[i])
				newStableSlope1, _ := sdk.NewDecFromStr(stableSlope1[i])
				newStableSlope2, _ := sdk.NewDecFromStr(stableSlope2[i])
				newLTV, _ := sdk.NewDecFromStr(ltv[i])
				newLiquidationThreshold, _ := sdk.NewDecFromStr(liquidationThreshold[i])
				newLiquidationPenalty, _ := sdk.NewDecFromStr(liquidationPenalty[i])
				newReserveFactor, _ := sdk.NewDecFromStr(reserveFactor[i])

				assetRatesStats = append(assetRatesStats, types.AssetRatesStats{
					AssetId:              assetId[i],
					UOptimal:             newUOptimal,
					Base:                 newBase,
					Slope1:               newSlope1,
					Slope2:               newSlope2,
					EnableStableBorrow:   newEnableStableBorrow,
					StableBase:           newStableBase,
					StableSlope1:         newStableSlope1,
					StableSlope2:         newStableSlope2,
					Ltv:                  newLTV,
					LiquidationThreshold: newLiquidationThreshold,
					LiquidationPenalty:   newLiquidationPenalty,
					ReserveFactor:        newReserveFactor,
				},
				)
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

			content := types.NewAddAssetRatesStats(title, description, assetRatesStats)

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
