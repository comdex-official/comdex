package cli

import (
	"strconv"
	"time"

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
		Use:   "add-assets [name] [Denom] [Decimals] [isOnchain] [assetOraclePrice]",
		Args:  cobra.ExactArgs(5),
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

			isOnchain, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}

			assetOraclePrice, err := ParseStringFromString(args[4], ",")
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
				newIsOnChain := ParseBoolFromString(isOnchain[i])
				newAssetOraclePrice := ParseBoolFromString(assetOraclePrice[i])
				assets = append(assets, types.Asset{
					Name:             names[i],
					Denom:            denoms[i],
					Decimals:         decimals[i],
					IsOnChain:        newIsOnChain,
					AssetOraclePrice: newAssetOraclePrice,
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

			asset_id, err := ParseUint64SliceFromString(args[0], ",")
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
			for i := range asset_id {
				newcollateralWeigt, err := sdk.NewDecFromStr(collateralWeight[i])
				if err != nil {
					return err
				}
				newliquidationThreshold, err := sdk.NewDecFromStr(liquidationThreshold[i])
				if err != nil {
					return err
				}
				newisBridgedAsset := ParseBoolFromString(isBridgedAsset[i])

				assets = append(assets, types.ExtendedAsset{
					AssetId:              asset_id[i],
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
			newcollateralWeight, err := sdk.NewDecFromStr(collateralWeight)
			if err != nil {
				return err
			}

			liquidationThreshold, err := cmd.Flags().GetString(flagLiquidationThreshold)
			if err != nil {
				return err
			}
			newliquidationThreshold, err := sdk.NewDecFromStr(liquidationThreshold)
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

			pair_id, err := ParseUint64SliceFromString(args[0], ",")
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
			for i := range pair_id {

				newbaseborrowrateasset1, err := sdk.NewDecFromStr(baseborrowrateasset1[i])
				if err != nil {
					return err
				}
				newbaseborrowrateasset2, err := sdk.NewDecFromStr(baseborrowrateasset2[i])
				if err != nil {
					return err
				}
				newbaselendrateasset1, err := sdk.NewDecFromStr(baselendrateasset1[i])
				if err != nil {
					return err
				}
				newbaselendrateasset2, err := sdk.NewDecFromStr(baselendrateasset2[i])
				if err != nil {
					return err
				}
				pairs = append(pairs, types.ExtendedPairLend{
					PairId:                pair_id[i],
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

func NewCmdSubmitAddAppMapingProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-app-mapping [name] [short_name] [min_gov_deposit] [gov_time_in_seconds]",
		Args:  cobra.ExactArgs(4),
		Short: "Add app mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]

			short_name := args[1]

			min_gov_deposit := args[2]

			gov_time_in_seconds, err := time.ParseDuration(args[3] + "s")
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

			var aMap []types.AppMapping
			var bMap []*types.MintGenesisToken
			new_min_gov_deposit, ok := sdk.NewIntFromString(min_gov_deposit)

			if err != nil {
				return err
			}
			if !ok {
				return types.ErrorInvalidMinGovSupply
			}
			aMap = append(aMap, types.AppMapping{
				Name:             name,
				ShortName:        short_name,
				MinGovDeposit:    new_min_gov_deposit,
				GovTimeInSeconds: gov_time_in_seconds.Seconds(),
				MintGenesisToken: bMap,
			})

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewAddAppMapingProposa(title, description, aMap)

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

func NewCmdSubmitAddAssetMapingProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-mapping [app_id] [asset_id] [genesis_supply] [isgovToken] [recipient]",
		Args:  cobra.ExactArgs(5),
		Short: "Add asset mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_id, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}
			genesis_supply, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}
			isgovToken, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}
			recipient, err := ParseStringFromString(args[4], ",")
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

			var aMap []types.AppMapping
			var bMap []*types.MintGenesisToken
			for i := range asset_id {
				newisgovToken := ParseBoolFromString(isgovToken[i])
				newgenesis_supply, ok := sdk.NewIntFromString(genesis_supply[i])
				address, err := sdk.AccAddressFromBech32(recipient[i])
				if err != nil {
					panic(err)
				}
				if !ok {
					return types.ErrorInvalidGenesisSupply
				}
				var cmap types.MintGenesisToken

				cmap.AssetId = asset_id[i]
				cmap.GenesisSupply = &newgenesis_supply
				cmap.IsgovToken = newisgovToken
				cmap.Recipient = address.String()

				bMap = append(bMap, &cmap)
			}
			aMap = append(aMap, types.AppMapping{
				Id:               app_id,
				MintGenesisToken: bMap,
			})

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewAddAssetMapingProposa(title, description, aMap)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
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

func NewCmdSubmitAddExtendedPairsVaultProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pairs-vault [app_mapping_id] [pair_id] [liquidation_ratio] [stability_fee] [closing_fee] [liquidation_penalty] [draw_down_fee] [is_vault_active] [debt_cieling] [debt_floor] [is_psm_pair] [min_cr] [pair_name] [asset_out_oracle_price] [asset_out_price]",
		Args:  cobra.ExactArgs(15),
		Short: "Add pairs vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_mapping_id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pair_id, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			liquidation_ratio, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}

			stability_fee, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}

			closing_fee, err := ParseStringFromString(args[4], ",")
			if err != nil {
				return err
			}

			liquidation_penalty, err := ParseStringFromString(args[5], ",")
			if err != nil {
				return err
			}

			draw_down_fee, err := ParseStringFromString(args[6], ",")
			if err != nil {
				return err
			}

			is_vault_active, err := ParseStringFromString(args[7], ",")
			if err != nil {
				return err
			}

			debt_cieling, err := ParseStringFromString(args[8], ",")
			if err != nil {
				return err
			}

			debt_floor, err := ParseStringFromString(args[9], ",")
			if err != nil {
				return err
			}

			is_psm_pair, err := ParseStringFromString(args[10], ",")
			if err != nil {
				return err
			}

			min_cr, err := ParseStringFromString(args[11], ",")
			if err != nil {
				return err
			}

			pair_name, err := ParseStringFromString(args[12], ",")
			if err != nil {
				return err
			}

			asset_out_oracle_price, err := ParseStringFromString(args[13], ",")
			if err != nil {
				return err
			}

			asset_out_price, err := ParseUint64SliceFromString(args[14], ",")
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

			var pairs []types.ExtendedPairVault
			for i := range pair_id {
				newliquidation_ratio, err := sdk.NewDecFromStr(liquidation_ratio[i])
				if err != nil {
					return err
				}
				newstability_fee, err := sdk.NewDecFromStr(stability_fee[i])
				if err != nil {
					return err
				}
				newclosing_fee, err := sdk.NewDecFromStr(closing_fee[i])
				if err != nil {
					return err
				}
				newliquidation_penalty, err := sdk.NewDecFromStr(liquidation_penalty[i])
				if err != nil {
					return err
				}
				new_draw_down_fee, err := sdk.NewDecFromStr(draw_down_fee[i])
				if err != nil {
					return err
				}
				newmin_cr, err := sdk.NewDecFromStr(min_cr[i])
				if err != nil {
					return err
				}
				newis_vault_active := ParseBoolFromString(is_vault_active[i])
				if err != nil {
					return err
				}
				debt_ceiling, ok := sdk.NewIntFromString(debt_cieling[i])
				if !ok {
					return types.ErrorInvalidDebtCeiling
				}
				debt_floor, ok := sdk.NewIntFromString(debt_floor[i])
				if !ok {
					return types.ErrorInvalidDebtFloor
				}
				newis_psm_pair := ParseBoolFromString(is_psm_pair[i])
				newasset_out_oracle_price := ParseBoolFromString(asset_out_oracle_price[i])
				pairs = append(pairs, types.ExtendedPairVault{
					AppMappingId:        app_mapping_id,
					PairId:              pair_id[i],
					LiquidationRatio:    newliquidation_ratio,
					StabilityFee:        newstability_fee,
					ClosingFee:          newclosing_fee,
					LiquidationPenalty:  newliquidation_penalty,
					DrawDownFee:         new_draw_down_fee,
					IsVaultActive:       newis_vault_active,
					DebtCeiling:         debt_ceiling,
					DebtFloor:           debt_floor,
					IsPsmPair:           newis_psm_pair,
					MinCr:               newmin_cr,
					PairName:            pair_name[i],
					AssetOutOraclePrice: newasset_out_oracle_price,
					AssetOutPrice:       asset_out_price[i],
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

			content := types.NewAddExtendedPairsVaultProposa(title, description, pairs)

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
