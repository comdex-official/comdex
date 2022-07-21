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
		txWithdraw(), //withdraw collateral partially or fully
		txDeposit(),
		txCloseLend(),
		txBorrowAsset(),
		txDrawAsset(),
		txRepayAsset(), //including functionality of both repaying and closing position
		txDepositBorrowAsset(),
		txCloseBorrowAsset(),
		txFundModuleAccounts(),
	)

	return cmd
}

func txLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lend [Asset_Id] [Amount] [Pool_Id] [App_Id]",
		Short: "lend a whitelisted asset",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[0], 10, 64)
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

			appId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgLend(ctx.GetFromAddress().String(), assetId, asset, pool, appId)

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
		Use:   "deposit [lendId] [Amount]",
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
		Use:   "close-lend [lendId]",
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

	assetIn, err := ParseUint64SliceFromString(newLendPairs.AssetIn, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOut, err := ParseUint64SliceFromString(newLendPairs.AssetOut, ",")
	if err != nil {
		return txf, nil, err
	}

	isInterPool, err := ParseUint64SliceFromString(newLendPairs.IsInterPool, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOutPoolID, err := ParseUint64SliceFromString(newLendPairs.AssetOutPoolID, ",")
	if err != nil {
		return txf, nil, err
	}

	var pairs []types.Extended_Pair
	for i := range assetIn {
		interPool := ParseBoolFromString(isInterPool[i])
		pairs = append(pairs, types.Extended_Pair{
			AssetIn:        assetIn[i],
			AssetOut:       assetOut[i],
			IsInterPool:    interPool,
			AssetOutPoolId: assetOutPoolID[i],
		})
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
			AssetId:   assetID[i],
			IsBridged: bridged,
		})
	}
	pool = types.Pool{
		ModuleName:           moduleName,
		MainAssetId:          mainAssetID,
		FirstBridgedAssetId:  firstBridgedAssetID,
		SecondBridgedAssetId: secondBridgedAssetID,
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
				AssetId: assetID,
				PoolId:  poolID,
				PairId:  pairIDs,
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

	assetID, err := ParseUint64SliceFromString(assetRatesStatsInput.AssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	uOptimal, err := ParseStringFromString(assetRatesStatsInput.UOptimal, ",")
	if err != nil {
		return txf, nil, err
	}
	base, err := ParseStringFromString(assetRatesStatsInput.Base, ",")
	if err != nil {
		return txf, nil, err
	}
	slope1, err := ParseStringFromString(assetRatesStatsInput.Slope1, ",")
	if err != nil {
		return txf, nil, err
	}
	slope2, err := ParseStringFromString(assetRatesStatsInput.Slope2, ",")
	if err != nil {
		return txf, nil, err
	}
	enableStableBorrow, err := ParseUint64SliceFromString(assetRatesStatsInput.EnableStableBorrow, ",")
	if err != nil {
		return txf, nil, err
	}
	stableBase, err := ParseStringFromString(assetRatesStatsInput.StableBase, ",")
	if err != nil {
		return txf, nil, err
	}
	stableSlope1, err := ParseStringFromString(assetRatesStatsInput.StableSlope1, ",")
	if err != nil {
		return txf, nil, err
	}
	stableSlope2, err := ParseStringFromString(assetRatesStatsInput.StableSlope2, ",")
	if err != nil {
		return txf, nil, err
	}
	ltv, err := ParseStringFromString(assetRatesStatsInput.LTV, ",")
	if err != nil {
		return txf, nil, err
	}
	liquidationThreshold, err := ParseStringFromString(assetRatesStatsInput.LiquidationThreshold, ",")
	if err != nil {
		return txf, nil, err
	}
	liquidationPenalty, err := ParseStringFromString(assetRatesStatsInput.LiquidationPenalty, ",")
	if err != nil {
		return txf, nil, err
	}
	reserveFactor, err := ParseStringFromString(assetRatesStatsInput.ReserveFactor, ",")
	if err != nil {
		return txf, nil, err
	}
	cAssetID, err := ParseUint64SliceFromString(assetRatesStatsInput.CAssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	var assetRatesStats []types.AssetRatesStats
	for i := range assetID {
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
			AssetId:              assetID[i],
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
			CAssetId:             cAssetID[i],
		},
		)
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
