package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	flag "github.com/spf13/pflag"
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
		Use:   "add-assets [name] [Denom] [Decimals] [isOnchain]",
		Args:  cobra.ExactArgs(4),
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
				assets = append(assets, types.Asset{
					Name:      names[i],
					Denom:     denoms[i],
					Decimals:  decimals[i],
					IsOnchain: newIsOnChain,
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

			assetId, err := ParseUint64SliceFromString(args[0], ",")
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
			for i := range assetId {
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
					AssetId:              assetId[i],
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

			pairId, err := ParseUint64SliceFromString(args[0], ",")
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
			for i := range pairId {

				newBaseBorrowRateAsset1, err := sdk.NewDecFromStr(baseborrowrateasset1[i])
				if err != nil {
					return err
				}
				newBaseBorrowRateAsset2, err := sdk.NewDecFromStr(baseborrowrateasset2[i])
				if err != nil {
					return err
				}
				newBaseLendRateAsset1, err := sdk.NewDecFromStr(baselendrateasset1[i])
				if err != nil {
					return err
				}
				newBaseLendRateAsset2, err := sdk.NewDecFromStr(baselendrateasset2[i])
				if err != nil {
					return err
				}
				pairs = append(pairs, types.ExtendedPairLend{
					PairId:                pairId[i],
					ModuleAcc:             moduleAccnt[i],
					BaseBorrowRateAsset_1: newBaseBorrowRateAsset1,
					BaseBorrowRateAsset_2: newBaseBorrowRateAsset2,
					BaseLendRateAsset_1:   newBaseLendRateAsset1,
					BaseLendRateAsset_2:   newBaseLendRateAsset2,
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

func NewCmdSubmitAddAppMappingProposal() *cobra.Command {
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

			shortName := args[1]

			minGovDeposit := args[2]

			govTimeIn, err := time.ParseDuration(args[3] + "s")
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
			var bMap []types.MintGenesisToken
			newMinGovDeposit, ok := sdk.NewIntFromString(minGovDeposit)

			if err != nil {
				return err
			}
			if !ok {
				return types.ErrorInvalidMinGovSupply
			}
			aMap = append(aMap, types.AppMapping{
				Name:             name,
				ShortName:        shortName,
				MinGovDeposit:    newMinGovDeposit,
				GovTimeInSeconds: govTimeIn.Seconds(),
				GenesisToken:     bMap,
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

func NewCmdSubmitAddAssetMappingProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-mapping [app_id] [asset_id] [genesis_supply] [isgovToken] [recipient]",
		Args:  cobra.ExactArgs(5),
		Short: "Add asset mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}
			genesisSupply, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}
			isGovToken, err := ParseStringFromString(args[3], ",")
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
			var bMap []types.MintGenesisToken
			for i := range assetId {
				newIsGovToken := ParseBoolFromString(isGovToken[i])
				newGenesisSupply, ok := sdk.NewIntFromString(genesisSupply[i])
				address, err := sdk.AccAddressFromBech32(recipient[i])
				if err != nil {
					panic(err)
				}
				if !ok {
					return types.ErrorInvalidGenesisSupply
				}
				var cmap types.MintGenesisToken

				cmap.AssetId = assetId[i]
				cmap.GenesisSupply = &newGenesisSupply
				cmap.IsgovToken = newIsGovToken
				cmap.Recipient = address.String()

				bMap = append(bMap, cmap)
			}
			aMap = append(aMap, types.AppMapping{
				Id:           appId,
				GenesisToken: bMap,
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
		Use:   "add-pairs-vault [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Add pairs vault ",
		Long:  `Must provide path to a extended pair vault JSON file (--extended-pair-vault-file) describing the extended pair to be created`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildCreateExtendedPairVaultMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateExtendedPaiVault())
	//cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagExtendedPairVaultFile)

	return cmd
}

func NewBuildCreateExtendedPairVaultMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	extPairVault, err := parseExtendPairVaultFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse extPairVault: %w", err)
	}

	appMappingId, err := strconv.ParseUint(extPairVault.AppMappingId, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	pairId, err := ParseUint64SliceFromString(extPairVault.PairId, ",")
	if err != nil {
		return txf, nil, err
	}

	liquidationRatio, err := ParseStringFromString(extPairVault.LiquidationRatio, ",")
	if err != nil {
		return txf, nil, err
	}

	stabilityFee, err := ParseStringFromString(extPairVault.StabilityFee, ",")
	if err != nil {
		return txf, nil, err
	}

	closingFee, err := ParseStringFromString(extPairVault.ClosingFee, ",")
	if err != nil {
		return txf, nil, err
	}

	liquidationPenalty, err := ParseStringFromString(extPairVault.LiquidationPenalty, ",")
	if err != nil {
		return txf, nil, err
	}

	drawDownFee, err := ParseStringFromString(extPairVault.DrawDownFee, ",")
	if err != nil {
		return txf, nil, err
	}

	isVaultActive, err := ParseStringFromString(extPairVault.IsVaultActive, ",")
	if err != nil {
		return txf, nil, err
	}

	debtCieling, err := ParseStringFromString(extPairVault.DebtCeiling, ",")
	if err != nil {
		return txf, nil, err
	}

	debtFloor, err := ParseStringFromString(extPairVault.DebtFloor, ",")
	if err != nil {
		return txf, nil, err
	}

	isPsmPair, err := ParseStringFromString(extPairVault.IsPsmPair, ",")
	if err != nil {
		return txf, nil, err
	}

	minCr, err := ParseStringFromString(extPairVault.MinCr, ",")
	if err != nil {
		return txf, nil, err
	}

	pairName, err := ParseStringFromString(extPairVault.PairName, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOutOraclePrice, err := ParseStringFromString(extPairVault.AssetOutOraclePrice, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOutPrice, err := ParseUint64SliceFromString(extPairVault.AssetOutPrice, ",")
	if err != nil {
		return txf, nil, err
	}

	minUsdValueLeft, err := ParseUint64SliceFromString(extPairVault.MinUsdValueLeft, ",")
	if err != nil {
		return txf, nil, err
	}

	title, err := ParseStringFromString(extPairVault.Title, ",")
	if err != nil {
		return txf, nil, err
	}

	description, err := ParseStringFromString(extPairVault.Description, ",")
	if err != nil {
		return txf, nil, err
	}

	from := clientCtx.GetFromAddress()

	var pairs []types.ExtendedPairVault
	for i := range pairId {
		newLiquidationRatio, err := sdk.NewDecFromStr(liquidationRatio[i])
		if err != nil {
			return txf, nil, err
		}
		newStabilityFee, err := sdk.NewDecFromStr(stabilityFee[i])
		if err != nil {
			return txf, nil, err
		}
		newClosingFee, err := sdk.NewDecFromStr(closingFee[i])
		if err != nil {
			return txf, nil, err
		}
		newLiquidationPenalty, err := sdk.NewDecFromStr(liquidationPenalty[i])
		if err != nil {
			return txf, nil, err
		}
		newDrawDownFee, err := sdk.NewDecFromStr(drawDownFee[i])
		if err != nil {
			return txf, nil, err
		}
		newMinCr, err := sdk.NewDecFromStr(minCr[i])
		if err != nil {
			return txf, nil, err
		}
		newIsVaultActive := ParseBoolFromString(isVaultActive[i])
		if err != nil {
			return txf, nil, err
		}
		debtCeiling, ok := sdk.NewIntFromString(debtCieling[i])
		if !ok {
			return txf, nil, types.ErrorInvalidDebtCeiling
		}
		newDebtFloor, ok := sdk.NewIntFromString(debtFloor[i])
		if !ok {
			return txf, nil, types.ErrorInvalidDebtFloor
		}
		newIsPsmPair := ParseBoolFromString(isPsmPair[i])
		newAssetOutOraclePrice := ParseBoolFromString(assetOutOraclePrice[i])
		pairs = append(pairs, types.ExtendedPairVault{
			AppMappingId:        appMappingId,
			PairId:              pairId[i],
			LiquidationRatio:    newLiquidationRatio,
			StabilityFee:        newStabilityFee,
			ClosingFee:          newClosingFee,
			LiquidationPenalty:  newLiquidationPenalty,
			DrawDownFee:         newDrawDownFee,
			IsVaultActive:       newIsVaultActive,
			DebtCeiling:         debtCeiling,
			DebtFloor:           newDebtFloor,
			IsPsmPair:           newIsPsmPair,
			MinCr:               newMinCr,
			PairName:            pairName[i],
			AssetOutOraclePrice: newAssetOutOraclePrice,
			AssetOutPrice:       assetOutPrice[i],
			MinUsdValueLeft:     minUsdValueLeft[i],
		})
	}

	deposit, err := sdk.ParseCoinsNormalized(extPairVault.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddExtendedPairsVaultProposa(title[0], description[0], pairs)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}
