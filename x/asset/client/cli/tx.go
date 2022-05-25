package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/asset/types"
)

func NewCmdSubmitAddAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-assets [name] [Denom] [Decimals]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			var assets []types.Asset
			for i := range names {
				assets = append(assets, types.Asset{
					Name:     names[i],
					Denom:    denoms[i],
					Decimals: decimals[i],
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

			content := types.NewAddAssetsProposal(title, description, assets)

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

func NewCmdSubmitUpdateAssetProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-asset [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update an Asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			denom, err := cmd.Flags().GetString(flagDenom)
			if err != nil {
				return err
			}

			decimals, err := cmd.Flags().GetInt64(flagDecimals)
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

			asset := types.Asset{
				Id:       id,
				Name:     name,
				Denom:    denom,
				Decimals: decimals,
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUpdateAssetProposal(title, description, asset)

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
	cmd.Flags().String(flagName, "", "name")
	cmd.Flags().String(flagDenom, "", "denomination")
	cmd.Flags().Int64(flagDecimals, -1, "decimals")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func NewCmdSubmitAddPairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pairs [asset-in] [asset-out]",
		Short: "Add a pair",
		Args:  cobra.ExactArgs(2),
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

			var pairs []types.Pair
			for i := range assetIn {
				pairs = append(pairs, types.Pair{
					AssetIn:  assetIn[i],
					AssetOut: assetOut[i],
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

			content := types.NewAddPairsProposal(title, description, pairs)

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

func NewCmdSubmitAddWhitelistedAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lend-assets [asset_id] [Collateral_Weight] [Liquidation_Threshold] [Is_Bridged_Asset]",
		Args:  cobra.ExactArgs(4),
		Short: "Add lend assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			collateralWeight, err := ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}

			liquidationThreshold, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}

			isBridgedAsset, err := ParseStringFromString(args[3], ",")
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

			from := ctx.GetFromAddress()

			var assets []types.ExtendedAsset
			for i := range assetID {
				newcollateralWeigt, _ := sdk.NewDecFromStr(collateralWeight[i])
				newliquidationThreshold, _ := sdk.NewDecFromStr(liquidationThreshold[i])
				newisBridgedAsset := ParseBoolFromString(isBridgedAsset[i])

				assets = append(assets, types.ExtendedAsset{
					AssetId:              assetID[i],
					CollateralWeight:     newcollateralWeigt,
					LiquidationThreshold: newliquidationThreshold,
					IsBridgedAsset:       newisBridgedAsset,
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

func NewCmdUpdateWhitelistedAssetProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-lend-asset [asset_id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update lend assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			collateralWeight, err := cmd.Flags().GetString(flagCollateralWeight)
			if err != nil {
				return err
			}
			newcollateralWeight, _ := sdk.NewDecFromStr(collateralWeight)

			liquidationThreshold, err := cmd.Flags().GetString(flagLiquidationThreshold)
			if err != nil {
				return err
			}
			newliquidationThreshold, _ := sdk.NewDecFromStr(liquidationThreshold)

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			isBridgedAsset, err := cmd.Flags().GetString(flagIsBridgedAsset)
			if err != nil {
				return err
			}

			newisBridgedAsset := ParseBoolFromString(isBridgedAsset)

			from := clientCtx.GetFromAddress()

			asset := types.ExtendedAsset{
				Id:                   id,
				CollateralWeight:     newcollateralWeight,
				LiquidationThreshold: newliquidationThreshold,
				IsBridgedAsset:       newisBridgedAsset,
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUpdateWhitelistedAssetProposal(title, description, asset)

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
	cmd.Flags().String(flagName, "", "name")
	cmd.Flags().String(flagDenom, "", "denomination")
	cmd.Flags().String(flagCollateralWeight, "", "collateralWeight")
	cmd.Flags().String(flagLiquidationThreshold, "", "liquidationThreshold")
	cmd.Flags().String(flagIsBridgedAsset, "", "isBridgedAsset")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func NewCmdAddWhitelistedPairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lend-asset-pairs [pair_id] [Module-Account] [Base_Borrow_Rate_Asset1] [Base_Borrow_Rate_Asset2] [Base_Lend_Rate_Asset1] [Base_Lend_Rate_Asset2]",
		Short: "Add lend asset pairs",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pairID, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			moduleAccnt, err := ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}

			baseborrowrateasset1, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}
			baseborrowrateasset2, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}
			baselendrateasset1, err := ParseStringFromString(args[4], ",")
			if err != nil {
				return err
			}
			baselendrateasset2, err := ParseStringFromString(args[5], ",")
			if err != nil {
				return err
			}

			var pairs []types.ExtendedPairLend
			for i := range pairID {
				newbaseborrowrateasset1, _ := sdk.NewDecFromStr(baseborrowrateasset1[i])
				newbaseborrowrateasset2, _ := sdk.NewDecFromStr(baseborrowrateasset2[i])
				newbaselendrateasset1, _ := sdk.NewDecFromStr(baselendrateasset1[i])
				newbaselendrateasset2, _ := sdk.NewDecFromStr(baselendrateasset2[i])
				pairs = append(pairs, types.ExtendedPairLend{
					PairId:                pairID[i],
					ModuleAcc:             moduleAccnt[i],
					BaseBorrowRateAsset_1: newbaseborrowrateasset1,
					BaseBorrowRateAsset_2: newbaseborrowrateasset2,
					BaseLendRateAsset_1:   newbaselendrateasset1,
					BaseLendRateAsset_2:   newbaselendrateasset2,
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

			content := types.NewAddWhitelistedPairsProposal(title, description, pairs)

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

func NewCmdUpdateWhitelistedPairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-lend-asset-pair [len_pair_id]",
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
			moduleAcc, err := cmd.Flags().GetString(flagModuleAcc)
			if err != nil {
				return err
			}
			baseborrowrateasset1, err := cmd.Flags().GetString(flagbaseborrowrateasset1)
			if err != nil {
				return err
			}
			newbaseborrowrateasset1, err := sdk.NewDecFromStr(baseborrowrateasset1)
			if err != nil {
				return err
			}

			baseborrowrateasset2, err := cmd.Flags().GetString(flagbaseborrowrateasset2)
			if err != nil {
				return err
			}
			newbaseborrowrateasset2, err := sdk.NewDecFromStr(baseborrowrateasset2)
			if err != nil {
				return err
			}

			baselendrateasset1, err := cmd.Flags().GetString(flagbaselendrateasset1)
			if err != nil {
				return err
			}
			newbaselendrateasset1, err := sdk.NewDecFromStr(baselendrateasset1)
			if err != nil {
				return err
			}

			baselendrateasset2, err := cmd.Flags().GetString(flagbaselendrateasset2)
			if err != nil {
				return err
			}
			newbaselendrateasset2, err := sdk.NewDecFromStr(baselendrateasset2)
			if err != nil {
				return err
			}

			pair := types.ExtendedPairLend{
				Id:                    id,
				ModuleAcc:             moduleAcc,
				BaseBorrowRateAsset_1: newbaseborrowrateasset1,
				BaseBorrowRateAsset_2: newbaseborrowrateasset2,
				BaseLendRateAsset_1:   newbaselendrateasset1,
				BaseLendRateAsset_2:   newbaselendrateasset2,
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

			content := types.NewUpdateWhitelistedPairProposal(title, description, pair)

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
	cmd.Flags().String(flagModuleAcc, "", "moduleAcc")
	cmd.Flags().String(flagbaseborrowrateasset1, "", "baseborrowrateasset1")
	cmd.Flags().String(flagbaseborrowrateasset2, "", "baseborrowrateasset2")
	cmd.Flags().String(flagbaselendrateasset1, "", "baselendrateasset1")
	cmd.Flags().String(flagbaselendrateasset2, "", "baselendrateasset2")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}
