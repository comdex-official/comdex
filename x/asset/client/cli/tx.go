package cli

import (
	"fmt"
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
		Use:   "add-assets [name] [Denom] [Decimals] [isOnChain]",
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

			isOnChain, err := ParseStringFromString(args[3], ",")
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
				newIsOnChain := ParseBoolFromString(isOnChain[i])
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
				newCollateralWeight, err := sdk.NewDecFromStr(collateralWeight[i])
				if err != nil {
					return err
				}
				newLiquidationThreshold, err := sdk.NewDecFromStr(liquidationThreshold[i])
				if err != nil {
					return err
				}
				newIsBridgedAsset := ParseBoolFromString(isBridgedAsset[i])

				assets = append(assets, types.ExtendedAsset{
					AssetId:              assetId[i],
					CollateralWeight:     newCollateralWeight,
					LiquidationThreshold: newLiquidationThreshold,
					IsBridgedAsset:       newIsBridgedAsset,
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
			newCollateralWeight, err := sdk.NewDecFromStr(collateralWeight)
			if err != nil {
				return err
			}

			liquidationThreshold, err := cmd.Flags().GetString(flagLiquidationThreshold)
			if err != nil {
				return err
			}
			newLiquidationThreshold, err := sdk.NewDecFromStr(liquidationThreshold)
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

			newIsBridgedAsset := ParseBoolFromString(isBridgedAsset)

			from := clientCtx.GetFromAddress()

			asset := types.ExtendedAsset{
				Id:                   id,
				CollateralWeight:     newCollateralWeight,
				LiquidationThreshold: newLiquidationThreshold,
				IsBridgedAsset:       newIsBridgedAsset,
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
		Use:   "add-lend-asset-pairs [flags]",
		Short: "Add lend asset pairs",
		Long: `Must provide path to a add white listed pairs JSON file (--add-white-whitelisted-pairs-file) describing the white listed pairs to be created
Sample json content
{
	"pair_id" :"",
	"module-account" :"",
	"base_borrow_rate_asset_1" :"",
	"base_borrow_rate_asset_2" :"",
	"base_lend_rate_asset_1" :"",
	"base_lend_rate_asset_2" :"",
	"title" :"",
	"description" :"",
	"deposit" :""
}`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateWhiteListedPairsMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateWhiteListedPairsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

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
			baseBorrowRateAsset1, err := cmd.Flags().GetString(flagBaseBorrowRateAsset1)
			if err != nil {
				return err
			}
			newBaseBorrowRateAsset1, err := sdk.NewDecFromStr(baseBorrowRateAsset1)
			if err != nil {
				return err
			}

			baseBorrowRateAsset2, err := cmd.Flags().GetString(flagBaseBorrowRateAsset2)
			if err != nil {
				return err
			}
			newBaseBorrowRateAsset2, err := sdk.NewDecFromStr(baseBorrowRateAsset2)
			if err != nil {
				return err
			}

			baseLendRateAsset1, err := cmd.Flags().GetString(flagBaseLendRateAsset1)
			if err != nil {
				return err
			}
			newBaseLendRateAsset1, err := sdk.NewDecFromStr(baseLendRateAsset1)
			if err != nil {
				return err
			}

			baseLendRateAsset2, err := cmd.Flags().GetString(flagBaseLendRateAsset2)
			if err != nil {
				return err
			}
			newBaseLendRateAsset2, err := sdk.NewDecFromStr(baseLendRateAsset2)
			if err != nil {
				return err
			}

			pair := types.ExtendedPairLend{
				Id:                    id,
				ModuleAcc:             moduleAcc,
				BaseBorrowRateAsset_1: newBaseBorrowRateAsset1,
				BaseBorrowRateAsset_2: newBaseBorrowRateAsset2,
				BaseLendRateAsset_1:   newBaseLendRateAsset1,
				BaseLendRateAsset_2:   newBaseLendRateAsset2,
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
	cmd.Flags().String(flagBaseBorrowRateAsset1, "", "base borrow rate asset1")
	cmd.Flags().String(flagBaseBorrowRateAsset2, "", "base borrow rate asset2")
	cmd.Flags().String(flagBaseLendRateAsset1, "", "base lend rate asset1")
	cmd.Flags().String(flagBaseLendRateAsset2, "", "base lend rate asset2")

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
		Use:   "add-asset-mapping [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Add asset mapping",
		Long: `Must provide path to a add asset mapping JSON file (--add-asset-mapping-file) describing the asset mapping to be created
Sample json content
{
	"app_id" :"",
	"asset_id" :"",
	"genesis_supply" :"",
	"is_gov_token" :"",
	"recipient" :"",
	"title" :"",
	"description" :"",
	"deposit" :""
}`,
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateAssetMappingMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)

		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateAssetMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func NewCmdSubmitAddExtendedPairsVaultProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pairs-vault [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Add pairs vault ",
		Long: `Must provide path to a extended pair vault JSON file (--extended-pair-vault-file) describing the extended pair to be created
Sample json content
{
	"app_mapping_id" : "",
	"pair_id" : "",
	"liquidation_ratio" : "",
	"stability_fee" : "",
	"closing_fee" : "",
	"liquidation_penalty" : "",
	"draw_down_fee" : "",
	"is_vault_active" : "",
	"debt_ceiling" : "",
	"debt_floor" : "",
	"is_psm_pair" : "",
	"min_cr" : "",
	"pair_name" : "",
	"asset_out_oracle_price" : "",
	"asset_out_price" : "",
	"min_usd_value_left" : "",
	"title" :"",
	"description" :"",
	"deposit":""
	
	
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateExtendedPairVaultMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateExtendedPairVault())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func NewCreateExtendedPairVaultMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	extPairVault, err := parseExtendPairVaultFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse extPairVault: %w", err)
	}

	appMappingId, err := strconv.ParseUint(extPairVault.AppMappingID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	pairId, err := ParseUint64SliceFromString(extPairVault.PairID, ",")
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

	debtCeiling, err := ParseStringFromString(extPairVault.DebtCeiling, ",")
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

	if extPairVault.Title == "" {
		return txf, nil, types.ErrorProposalTitleMissing
	}

	if extPairVault.Description == "" {
		return txf, nil, types.ErrorProposalDescriptionMissing
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
		debtCeiling, ok := sdk.NewIntFromString(debtCeiling[i])
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

	content := types.NewAddExtendedPairsVaultProposa(extPairVault.Title, extPairVault.Description, pairs)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

func NewCreateAssetMappingMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetMapping, err := parseAssetMappingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse assetMapping: %w", err)
	}

	appId, err := strconv.ParseUint(assetMapping.AppID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetId, err := ParseUint64SliceFromString(assetMapping.AssetID, ",")
	if err != nil {
		return txf, nil, err
	}
	genesisSupply, err := ParseStringFromString(assetMapping.GenesisSupply, ",")
	if err != nil {
		return txf, nil, err
	}
	isGovToken, err := ParseStringFromString(assetMapping.IsGovToken, ",")
	if err != nil {
		return txf, nil, err
	}
	recipient, err := ParseStringFromString(assetMapping.Recipient, ",")
	if err != nil {
		return txf, nil, err
	}

	if assetMapping.Title == "" {
		return txf, nil, types.ErrorProposalTitleMissing
	}

	if assetMapping.Description == "" {
		return txf, nil, types.ErrorProposalDescriptionMissing
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
			return txf, nil, types.ErrorInvalidGenesisSupply
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

	deposit, err := sdk.ParseCoinsNormalized(assetMapping.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddAssetMapingProposa(assetMapping.Title, assetMapping.Description, aMap)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

func NewCreateWhiteListedPairsMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	whiteListedPairs, err := parseWhiteListedPairsFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse whiteListedPairs: %w", err)
	}

	pairId, err := ParseUint64SliceFromString(whiteListedPairs.PairID, ",")
	if err != nil {
		return txf, nil, err
	}

	moduleAccount, err := ParseStringFromString(whiteListedPairs.ModuleAccount, ",")
	if err != nil {
		return txf, nil, err
	}

	baseBorrowRateAsset1, err := ParseStringFromString(whiteListedPairs.BaseBorrowRateAsset1, ",")
	if err != nil {
		return txf, nil, err
	}
	baseBorrowRateAsset2, err := ParseStringFromString(whiteListedPairs.BaseBorrowRateAsset2, ",")
	if err != nil {
		return txf, nil, err
	}
	baseLendRateAsset1, err := ParseStringFromString(whiteListedPairs.BaseLendRateAsset1, ",")
	if err != nil {
		return txf, nil, err
	}
	baseLendRateAsset2, err := ParseStringFromString(whiteListedPairs.BaseLendRateAsset2, ",")
	if err != nil {
		return txf, nil, err
	}

	var pairs []types.ExtendedPairLend
	for i := range pairId {

		newBaseBorrowRateAsset1, err := sdk.NewDecFromStr(baseBorrowRateAsset1[i])
		if err != nil {
			return txf, nil, err
		}
		newBaseBorrowRateAsset2, err := sdk.NewDecFromStr(baseBorrowRateAsset2[i])
		if err != nil {
			return txf, nil, err
		}
		newBaseLendRateAsset1, err := sdk.NewDecFromStr(baseLendRateAsset1[i])
		if err != nil {
			return txf, nil, err
		}
		newBaseLendRateAsset2, err := sdk.NewDecFromStr(baseLendRateAsset2[i])
		if err != nil {
			return txf, nil, err
		}
		pairs = append(pairs, types.ExtendedPairLend{
			PairId:                pairId[i],
			ModuleAcc:             moduleAccount[i],
			BaseBorrowRateAsset_1: newBaseBorrowRateAsset1,
			BaseBorrowRateAsset_2: newBaseBorrowRateAsset2,
			BaseLendRateAsset_1:   newBaseLendRateAsset1,
			BaseLendRateAsset_2:   newBaseLendRateAsset2,
		})
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(whiteListedPairs.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddWhitelistedPairsProposal(whiteListedPairs.Title, whiteListedPairs.Description, pairs)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}
