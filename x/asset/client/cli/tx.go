package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/comdex-official/comdex/x/asset/types"
)

func NewCmdSubmitAddMultipleAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-multiple-assets [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit multiple assets",
		Long: `Must provide path to a add assets in JSON file (--add-assets) describing the asset in app to be created
Sample json content
{
	"name" :"ATOM,CMDX,CMST,OSMO,cATOM,cCMDX,cCMST,cOSMO,HARBOR",
	"denom" :"uatom,ucmdx,ucmst,uosmo,ucatom,uccmdx,uccmst,ucosmo,uharbor",
	"decimals" :"1000000,1000000,1000000,1000000,1000000,1000000,1000000,1000000,1000000",
	"is_on_chain" :"0,0,0,0,0,0,0,0,1",
	"asset_oracle_price" :"1,1,0,1,0,0,0,0,0",
	"title" :"Add assets for applications to be deployed on comdex testnet",
	"description" :"This proposal it to add following assets ATOM,CMDX,CMST,OSMO,cATOM,cCMDX,cCMST,cOSMO,HARBOR to be then used on harbor, commodo and cswap apps",
	"deposit" :"1000000000ucmdx"
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateMultipleAssets(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateAssetsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func NewCreateMultipleAssets(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetsMapping, err := parseAssetsMappingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse assetsMapping: %w", err)
	}

	names, err := ParseStringFromString(assetsMapping.Name, ",")
	if err != nil {
		return txf, nil, err
	}
	denoms, err := ParseStringFromString(assetsMapping.Denom, ",")
	if err != nil {
		return txf, nil, err
	}

	decimals, err := ParseStringFromString(assetsMapping.Decimals, ",")
	if err != nil {
		return txf, nil, err
	}

	isOnChain, err := ParseStringFromString(assetsMapping.IsOnChain, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOraclePrice, err := ParseStringFromString(assetsMapping.AssetOraclePrice, ",")
	if err != nil {
		return txf, nil, err
	}
	isCdpMintable, err := ParseStringFromString(assetsMapping.IsCdpMintable, ",")
	if err != nil {
		return txf, nil, err
	}

	from := clientCtx.GetFromAddress()

	var assets []types.Asset
	for i := range names {
		newIsOnChain := ParseBoolFromString(isOnChain[i])
		newAssetOraclePrice := ParseBoolFromString(assetOraclePrice[i])
		newIsCdpMintable := ParseBoolFromString(isCdpMintable[i])
		newDecimals, ok := sdk.NewIntFromString(decimals[i])
		if !ok {
			return txf, nil, types.ErrorInvalidDecimals
		}
		assets = append(assets, types.Asset{
			Name:                  names[i],
			Denom:                 denoms[i],
			Decimals:              newDecimals,
			IsOnChain:             newIsOnChain,
			IsOraclePriceRequired: newAssetOraclePrice,
			IsCdpMintable:         newIsCdpMintable,
		})
	}

	deposit, err := sdk.ParseCoinsNormalized(assetsMapping.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddMultipleAssetsProposal(assetsMapping.Title, assetsMapping.Description, assets)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func NewCmdSubmitAddAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-assets [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit assets",
		Long: `Must provide path to a add assets in JSON file (--add-assets) describing the asset in app to be created
Sample json content
{
	"name" :"ATOM",
	"denom" :"uatom",
	"decimals" :"1000000",
	"is_on_chain" :"0",
	"asset_oracle_price" :"1",
	"title" :"Add assets for applications to be deployed on comdex testnet",
	"description" :"This proposal it to add following assets ATOM to be then used on harbor, commodo and cswap apps",
	"deposit" :"1000000000ucmdx"
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateAssets(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateAssetsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func NewCreateAssets(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetsMapping, err := parseAssetsMappingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse assetsMapping: %w", err)
	}

	names := assetsMapping.Name

	denoms := assetsMapping.Denom

	decimals := assetsMapping.Decimals

	isOnChain := ParseBoolFromString(assetsMapping.IsOnChain)
	if err != nil {
		return txf, nil, err
	}

	assetOraclePrice := ParseBoolFromString(assetsMapping.AssetOraclePrice)
	if err != nil {
		return txf, nil, err
	}

	isCdpMintable := ParseBoolFromString(assetsMapping.IsCdpMintable)
	if err != nil {
		return txf, nil, err
	}

	from := clientCtx.GetFromAddress()

	newDecimals, ok := sdk.NewIntFromString(decimals)
	if !ok {
		return txf, nil, types.ErrorInvalidDecimals
	}
	assets := types.Asset{
		Name:                  names,
		Denom:                 denoms,
		Decimals:              newDecimals,
		IsOnChain:             isOnChain,
		IsOraclePriceRequired: assetOraclePrice,
		IsCdpMintable:         isCdpMintable,
	}

	deposit, err := sdk.ParseCoinsNormalized(assetsMapping.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddAssetsProposal(assetsMapping.Title, assetsMapping.Description, assets)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func NewCmdSubmitUpdateAssetProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-asset [id] [name] [denom] [decimal] [is-oracle-price-required]",
		Args:  cobra.ExactArgs(5),
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

			name := args[1]

			denom := args[2]

			decimals, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrorInvalidDecimals
			}
			assetOraclePrice := args[4]

			newAssetOraclePrice := ParseBoolFromString(assetOraclePrice)

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
				Id:                    id,
				Name:                  name,
				Denom:                 denom,
				Decimals:              decimals,
				IsOraclePriceRequired: newAssetOraclePrice,
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

			assetIn, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetOut, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			pairs := types.Pair{
				AssetIn:  assetIn,
				AssetOut: assetOut,
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

func NewCmdSubmitUpdatePairProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pair [id] [asset-in] [asset-out]",
		Short: "Update a pair",
		Args:  cobra.ExactArgs(3),
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

			pair := types.Pair{
				Id:       id,
				AssetIn:  assetIn,
				AssetOut: assetOut,
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

			content := types.NewUpdatePairProposal(title, description, pair)

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

func NewCmdSubmitAddMultiplePairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-multiple-pairs [asset-in] [asset-out]",
		Short: "Add multiple pair",
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

			content := types.NewAddMultiplePairsProposal(title, description, pairs)

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

func NewCmdSubmitAddAppProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-app [name] [short_name] [min_gov_deposit] [gov_time_in_seconds]",
		Args:  cobra.ExactArgs(4),
		Short: "Add app",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]

			shortName := args[1]

			minGovDeposit := args[2]

			govTimeIn, err := strconv.ParseUint(args[3], 10, 64)
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

			var aMap types.AppData
			var bMap []types.MintGenesisToken
			newMinGovDeposit, ok := sdk.NewIntFromString(minGovDeposit)

			if err != nil {
				return err
			}
			if !ok {
				return types.ErrorInvalidMinGovSupply
			}
			aMap = types.AppData{
				Name:             name,
				ShortName:        shortName,
				MinGovDeposit:    newMinGovDeposit,
				GovTimeInSeconds: govTimeIn,
				GenesisToken:     bMap,
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewAddAppProposal(title, description, aMap)

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

func NewCmdSubmitUpdateGovTimeInAppProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gov-time-app [app_id] [gov_time_in_seconds] [min_gov_deposit]",
		Args:  cobra.ExactArgs(3),
		Short: "Update gov time in app",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			govTimeIn, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			minGovDeposit, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return types.ErrorInvalidMinGovSupply
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

			aTime := types.AppAndGovTime{
				AppId:            appID,
				GovTimeInSeconds: govTimeIn,
				MinGovDeposit:    minGovDeposit,
			}

			content := types.NewUpdateGovTimeInAppProposal(title, description, aTime)

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

func NewCmdSubmitAddAssetInAppProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-in-app [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Add asset in app",
		Long: `Must provide path to a add asset in app JSON file (--add-asset-in-app-file) describing the asset in app to be created
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

			txf, msg, err := NewCreateAssetInAppMsg(clientCtx, txf, cmd.Flags())
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

func NewCreateAssetInAppMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetMapping, err := parseAssetMappingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse assetMapping: %w", err)
	}

	appID, err := strconv.ParseUint(assetMapping.AppID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetIDs, err := strconv.ParseUint(assetMapping.AssetID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	genesisSupply := assetMapping.GenesisSupply

	isGovToken := ParseBoolFromString(assetMapping.IsGovToken)
	if err != nil {
		return txf, nil, err
	}
	recipient := assetMapping.Recipient
	if err != nil {
		return txf, nil, err
	}

	if assetMapping.Title == "" {
		return txf, nil, types.ErrorProposalTitleMissing
	}

	if assetMapping.Description == "" {
		return txf, nil, types.ErrorProposalDescriptionMissing
	}

	var bMap []types.MintGenesisToken
	newGenesisSupply, ok := sdk.NewIntFromString(genesisSupply)
	if !ok {
		return txf, nil, types.ErrorInvalidGenesisSupply
	}
	address, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		panic(err)
	}
	var cmap types.MintGenesisToken

	cmap.AssetId = assetIDs
	cmap.GenesisSupply = newGenesisSupply
	cmap.IsGovToken = isGovToken
	cmap.Recipient = address.String()

	bMap = append(bMap, cmap)

	aMap := types.AppData{
		Id:           appID,
		GenesisToken: bMap,
	}

	deposit, err := sdk.ParseCoinsNormalized(assetMapping.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddAssetInAppProposal(assetMapping.Title, assetMapping.Description, aMap)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}

func NewCmdSubmitAddMultipleAssetsPairsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-multiple-assets-pairs [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit multiple assets and pairs",
		Long: `Must provide path to a add assets and pairs in JSON file (--add-assets-pairs-file) describing the assets and pairs to be created
Sample json content
{
	"name" :"ATOM,CMDX,OSMO,cATOM,cCMDX,cCMST,cOSMO,HARBOR",
	"denom" :"uatom,ucmdx,uosmo,ucatom,uccmdx,uccmst,ucosmo,uharbor",
	"decimals" :"1000000,1000000,1000000,1000000,1000000,1000000,1000000,1000000",
	"is_on_chain" :"0,0,0,0,0,0,0,1",
	"asset_oracle_price" :"1,1,1,0,0,0,0,0",
    "is_cdp_mintable" :"1,1,1,0,0,0,0,0",
    "asset_out" :"3,3,3,3,3,3,3,3",
	"title" :"Add assets and pairs for applications to be deployed on comdex testnet",
	"description" :"This proposal it to add following assets ATOM,CMDX,CMST,OSMO,cATOM,cCMDX,cCMST,cOSMO,HARBOR with pairs to be then used on harbor, commodo and cswap apps",
	"deposit" :"1000000000ucmdx"
}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateMultipleAssetsPairs(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateAssetsPairsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func NewCreateMultipleAssetsPairs(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	assetsPairsMapping, err := parseAssetsPairsMappingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse assetsPairs: %w", err)
	}

	names, err := ParseStringFromString(assetsPairsMapping.Name, ",")
	if err != nil {
		return txf, nil, err
	}
	denoms, err := ParseStringFromString(assetsPairsMapping.Denom, ",")
	if err != nil {
		return txf, nil, err
	}

	decimals, err := ParseStringFromString(assetsPairsMapping.Decimals, ",")
	if err != nil {
		return txf, nil, err
	}

	isOnChain, err := ParseStringFromString(assetsPairsMapping.IsOnChain, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOraclePrice, err := ParseStringFromString(assetsPairsMapping.AssetOraclePrice, ",")
	if err != nil {
		return txf, nil, err
	}
	isCdpMintable, err := ParseStringFromString(assetsPairsMapping.IsCdpMintable, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOut, err := ParseStringFromString(assetsPairsMapping.AssetOut, ",")
	if err != nil {
		return txf, nil, err
	}

	from := clientCtx.GetFromAddress()

	var assets []types.AssetPair
	for i := range names {
		newIsOnChain := ParseBoolFromString(isOnChain[i])
		newAssetOraclePrice := ParseBoolFromString(assetOraclePrice[i])
		newIsCdpMintable := ParseBoolFromString(isCdpMintable[i])
		newDecimals, ok := sdk.NewIntFromString(decimals[i])
		if !ok {
			return txf, nil, types.ErrorInvalidDecimals
		}
		newAssetOut, _ := sdk.NewIntFromString(assetOut[i])

		assets = append(assets, types.AssetPair{
			Name:                  names[i],
			Denom:                 denoms[i],
			Decimals:              newDecimals,
			IsOnChain:             newIsOnChain,
			IsOraclePriceRequired: newAssetOraclePrice,
			IsCdpMintable:         newIsCdpMintable,
			AssetOut:              newAssetOut.Uint64(),
		})
	}

	deposit, err := sdk.ParseCoinsNormalized(assetsPairsMapping.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAddMultipleAssetsPairsProposal(assetsPairsMapping.Title, assetsPairsMapping.Description, assets)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

func NewAddAssetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset [name] [denom] [decimals] [is-on-chain] [is-oracle-price-required] [is-cdp-mintable]",
		Args:  cobra.ExactArgs(6),
		Short: "Register/Whitelist a new asset in asset module ",
		Long: strings.TrimSpace(
			fmt.Sprintf(` Add / Register / Whitelist new asset in asset module
Example:
$ %s tx %s add-asset CMDX ucmdx 1000000 true true true --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			denom := args[1]

			decimals, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("parse decimals: %w", err)
			}

			isOnChain, err := strconv.ParseBool(args[3])
			if err != nil {
				return fmt.Errorf("parse is-on-chain: %w", err)
			}

			isOraclePriceRequired, err := strconv.ParseBool(args[4])
			if err != nil {
				return fmt.Errorf("parse is-oracle-price-required: %w", err)
			}

			isCdpMintable, err := strconv.ParseBool(args[5])
			if err != nil {
				return fmt.Errorf("parse is-cdp-mintable: %w", err)
			}

			msg := types.NewMsgAddAsset(clientCtx.GetFromAddress(), name, denom, decimals, isOnChain, isOraclePriceRequired, isCdpMintable)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
