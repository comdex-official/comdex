package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"strconv"

	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module.
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
		txWithdraw(),
		txDeposit(),
		txCloseLend(),
		txBorrowAsset(),
		txDrawAsset(),
		txRepayAsset(),
		txDepositBorrowAsset(),
		txCloseBorrowAsset(),
		txBorrowAssetAlternate(),
		txFundModuleAccounts(),
		txCalculateBorrowInterest(),
		txCalculateLendRewards(),
	)

	return cmd
}

func txLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend [asset-id] [amount] [pool-id] [app-id]",
		Short: "lend a whitelisted asset",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[0], 10, 64)
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

			appID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgLend(ctx.GetFromAddress().String(), assetID, asset, pool, appID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [lend-id] [amount]",
		Short: "withdraw lent asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(ctx.GetFromAddress().String(), lendID, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [lend-id] [amount]",
		Short: "deposit into a lent position",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(ctx.GetFromAddress().String(), lendID, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txCloseLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-lend [lend-id]",
		Short: "close a lent position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			lendID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCloseLend(ctx.GetFromAddress().String(), lendID)

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
			lendID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pairID, err := strconv.ParseUint(args[1], 10, 64)
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

			msg := types.NewMsgBorrow(ctx.GetFromAddress().String(), lendID, pairID, isStableBorrow, amountIn, amountOut)

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

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRepay(ctx.GetFromAddress().String(), borrowID, asset)

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

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDraw(ctx.GetFromAddress().String(), borrowID, asset)

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

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositBorrow(ctx.GetFromAddress().String(), borrowID, asset)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txCloseBorrowAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-borrow [borrow-id] ",
		Short: " close borrow position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCloseBorrow(ctx.GetFromAddress().String(), borrowID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txBorrowAssetAlternate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "borrow-alternate [asset-id] [pool-id] [amount-in] [pair-id] [is-stable-borrow] [amount-out] [app-id]",
		Short: "directly borrow from a whitelisted asset",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amountIn, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			pairID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			stableBorrow, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			isStableBorrow := ParseBoolFromString(stableBorrow)

			amountOut, err := sdk.ParseCoinNormalized(args[5])
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgBorrowAlternate(ctx.GetFromAddress().String(), assetID, poolID, amountIn, pairID, isStableBorrow, amountOut, appID)

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

			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgFundModuleAccounts(moduleName, assetID, ctx.GetFromAddress().String(), amount)
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

func CmdAddNewLendPairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lend-pairs [flags]",
		Short: "Add lend asset pairs",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateNewLendPairs(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetNewLendPairsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewCreateNewLendPairs(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	newLendPairs, err := parseAddNewLendPairsFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse add lend pairs : %w", err)
	}

	assetIn, err := strconv.ParseUint(newLendPairs.AssetIn, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetOut, err := strconv.ParseUint(newLendPairs.AssetOut, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	isInterPool, err := strconv.ParseUint(newLendPairs.IsInterPool, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetOutPoolID, err := strconv.ParseUint(newLendPairs.AssetOutPoolID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	minUSDValueLeft, err := strconv.ParseUint(newLendPairs.MinUSDValueLeft, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	interPool := ParseBoolFromString(isInterPool)
	pairs := types.Extended_Pair{
		AssetIn:         assetIn,
		AssetOut:        assetOut,
		IsInterPool:     interPool,
		AssetOutPoolID:  assetOutPoolID,
		MinUsdValueLeft: minUSDValueLeft,
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(newLendPairs.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddLendPairsProposal(newLendPairs.Title, newLendPairs.Description, pairs)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

func CmdUpdateLendPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-lend-pair [id] [asset-in] [asset-out] [min-usd-value-left]",
		Short: "Update a lend pair",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetIn, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			assetOut, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			minUsdValueLeft, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			pair := types.Extended_Pair{
				Id:              id,
				AssetIn:         assetIn,
				AssetOut:        assetOut,
				MinUsdValueLeft: minUsdValueLeft,
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

			content := types.NewUpdateLendPairsProposal(title, description, pair)

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
		Use:   "add-lend-pool [flag] ",
		Short: "Add lend pool ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateLendPool(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetAddLendPoolMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewCreateLendPool(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	newLendPool, err := parseAddPoolFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse add lend pool : %w", err)
	}

	moduleName := newLendPool.ModuleName
	cPoolName := newLendPool.CPoolName

	mainAssetID, err := strconv.ParseUint(newLendPool.MainAssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	firstBridgedAssetID, err := strconv.ParseUint(newLendPool.FirstBridgedAssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	secondBridgedAssetID, err := strconv.ParseUint(newLendPool.SecondBridgedAssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	reserveFunds, err := strconv.ParseUint(newLendPool.ReserveFunds, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetID, err := ParseUint64SliceFromString(newLendPool.AssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	isBridgedAsset, err := ParseUint64SliceFromString(newLendPool.IsBridgedAsset, ",")
	if err != nil {
		return txf, nil, err
	}
	var pool types.Pool
	var assetData []types.AssetDataPoolMapping

	for i := range assetID {
		bridged := ParseBoolFromString(isBridgedAsset[i])
		assetData = append(assetData, types.AssetDataPoolMapping{
			AssetID:   assetID[i],
			IsBridged: bridged,
		})
	}
	pool = types.Pool{
		ModuleName:           moduleName,
		MainAssetId:          mainAssetID,
		FirstBridgedAssetID:  firstBridgedAssetID,
		SecondBridgedAssetID: secondBridgedAssetID,
		CPoolName:            cPoolName,
		ReserveFunds:         reserveFunds,
		AssetData:            assetData,
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(newLendPool.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddPoolProposal(newLendPool.Title, newLendPool.Description, pool)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func CmdAddAssetToPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-to-pair-mapping [asset_id] [pool_id] [pair_id] ",
		Short: "Add asset to pair mapping ",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			rawPairID, _ := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}
			var pairIDs []uint64
			for i := range rawPairID {
				pairIDs = append(pairIDs, rawPairID[i])
			}
			assetToPairMapping := types.AssetToPairMapping{
				AssetID: assetID,
				PoolID:  poolID,
				PairID:  pairIDs,
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

func CmdAddNewAssetRatesStatsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-rates-stats [flags]",
		Short: "Add lend asset pairs",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateAssetRatesStats(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetAddAssetRatesStatsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewCreateAssetRatesStats(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetRatesStatsInput, err := parseAssetRateStatsFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse asset rates stats : %w", err)
	}

	assetID, err := strconv.ParseUint(assetRatesStatsInput.AssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	uOptimal := assetRatesStatsInput.UOptimal

	base := assetRatesStatsInput.Base

	slope1 := assetRatesStatsInput.Slope1

	slope2 := assetRatesStatsInput.Slope2

	enableStableBorrow, err := strconv.ParseUint(assetRatesStatsInput.EnableStableBorrow, 10, 64)
	if err != nil {
		return txf, nil, err
	}
	stableBase := assetRatesStatsInput.StableBase

	stableSlope1 := assetRatesStatsInput.StableSlope1

	stableSlope2 := assetRatesStatsInput.StableSlope2

	ltv := assetRatesStatsInput.LTV

	liquidationThreshold := assetRatesStatsInput.LiquidationThreshold

	liquidationPenalty := assetRatesStatsInput.LiquidationPenalty

	liquidationBonus := assetRatesStatsInput.LiquidationBonus

	reserveFactor := assetRatesStatsInput.ReserveFactor

	cAssetID, err := strconv.ParseUint(assetRatesStatsInput.CAssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	newUOptimal, _ := sdk.NewDecFromStr(uOptimal)
	newBase, _ := sdk.NewDecFromStr(base)
	newSlope1, _ := sdk.NewDecFromStr(slope1)
	newSlope2, _ := sdk.NewDecFromStr(slope2)
	newEnableStableBorrow := ParseBoolFromString(enableStableBorrow)
	newStableBase, _ := sdk.NewDecFromStr(stableBase)
	newStableSlope1, _ := sdk.NewDecFromStr(stableSlope1)
	newStableSlope2, _ := sdk.NewDecFromStr(stableSlope2)
	newLTV, _ := sdk.NewDecFromStr(ltv)
	newLiquidationThreshold, _ := sdk.NewDecFromStr(liquidationThreshold)
	newLiquidationPenalty, _ := sdk.NewDecFromStr(liquidationPenalty)
	newLiquidationBonus, _ := sdk.NewDecFromStr(liquidationBonus)
	newReserveFactor, _ := sdk.NewDecFromStr(reserveFactor)

	assetRatesStats := types.AssetRatesStats{
		AssetID:              assetID,
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
		LiquidationBonus:     newLiquidationBonus,
		ReserveFactor:        newReserveFactor,
		CAssetID:             cAssetID,
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(assetRatesStatsInput.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddAssetRatesStats(assetRatesStatsInput.Title, assetRatesStatsInput.Description, assetRatesStats)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func CmdAddNewAuctionParamsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-auction-params [flags]",
		Short: "Add auction params",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewAddAuctionParams(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetAuctionParams())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewAddAuctionParams(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	auctionParamsInput, err := parseAuctionPramsFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse auction params : %w", err)
	}

	appID, err := strconv.ParseUint(auctionParamsInput.AppID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	auctionDurationSeconds, err := strconv.ParseUint(auctionParamsInput.AuctionDurationSeconds, 10, 64)
	if err != nil {
		return txf, nil, err
	}
	buffer, _ := sdk.NewDecFromStr(auctionParamsInput.Buffer)

	cusp, _ := sdk.NewDecFromStr(auctionParamsInput.Cusp)

	step, _ := sdk.NewIntFromString(auctionParamsInput.Step)

	priceFunctionType, err := strconv.ParseUint(auctionParamsInput.PriceFunctionType, 10, 64)
	if err != nil {
		return txf, nil, err
	}
	dutchID, err := strconv.ParseUint(auctionParamsInput.DutchID, 10, 64)
	if err != nil {
		return txf, nil, err
	}
	bidDurationSeconds, err := strconv.ParseUint(auctionParamsInput.BidDurationSeconds, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	auctionParams := types.AuctionParams{
		AppId:                  appID,
		AuctionDurationSeconds: auctionDurationSeconds,
		Buffer:                 buffer,
		Cusp:                   cusp,
		Step:                   step,
		PriceFunctionType:      priceFunctionType,
		DutchId:                dutchID,
		BidDurationSeconds:     bidDurationSeconds,
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(auctionParamsInput.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddAuctionParams(auctionParamsInput.Title, auctionParamsInput.Description, auctionParams)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func txCalculateBorrowInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate-borrow-interest [borrow-id] ",
		Short: " calculate borrow interest for a borrow position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCalculateBorrowInterest(ctx.GetFromAddress().String(), borrowID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txCalculateLendRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate-lend-rewards [lend-id] ",
		Short: " calculate lend rewards for a lend position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			borrowID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCalculateLendRewards(ctx.GetFromAddress().String(), borrowID)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
