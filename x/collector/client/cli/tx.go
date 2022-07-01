package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"strconv"

	"github.com/comdex-official/comdex/x/collector/types"
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
	return cmd
}

// NewCmdLookupTableParams cmd for lookup table param proposal updates.
func NewCmdLookupTableParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-params [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Set collector lookup params for collector module",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateLookupTableParams(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateLookupTableParamsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewCreateLookupTableParams(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	lookupTableParams, err := parseLookupTableParamsFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse lookup table params: %w", err)
	}

	appID, err := ParseUint64SliceFromString(lookupTableParams.AppID, ",")
	if err != nil {
		return txf, nil, err
	}

	collectorAssetID, err := ParseUint64SliceFromString(lookupTableParams.CollectorAssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	secondaryAssetID, err := ParseUint64SliceFromString(lookupTableParams.SecondaryAssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	surplusThreshold, err := ParseUint64SliceFromString(lookupTableParams.SurplusThreshold, ",")
	if err != nil {
		return txf, nil, err
	}

	debtThreshold, err := ParseUint64SliceFromString(lookupTableParams.DebtThreshold, ",")
	if err != nil {
		return txf, nil, err
	}

	lockerSavingRate, err := ParseStringFromString(lookupTableParams.LockerSavingRate, ",")
	if err != nil {
		return txf, nil, err
	}

	lotSize, err := ParseUint64SliceFromString(lookupTableParams.LotSize, ",")
	if err != nil {
		return txf, nil, err
	}

	bidFactor, err := ParseStringFromString(lookupTableParams.BidFactor, ",")
	if err != nil {
		return txf, nil, err
	}

	debtLotSize, err := ParseUint64SliceFromString(lookupTableParams.DebtLotSize, ",")
	if err != nil {
		return txf, nil, err
	}

	var LookupTableRecords []types.CollectorLookupTable
	for i := range appID {
		newBidFactor, err := sdk.NewDecFromStr(bidFactor[i])
		if err != nil {
			return txf, nil, err
		}

		newLockerSavingRate, err := sdk.NewDecFromStr(lockerSavingRate[i])
		if err != nil {
			return txf, nil, err
		}

		LookupTableRecords = append(LookupTableRecords, types.CollectorLookupTable{
			AppId:            appID[i],
			CollectorAssetId: collectorAssetID[i],
			SecondaryAssetId: secondaryAssetID[i],
			SurplusThreshold: surplusThreshold[i],
			DebtThreshold:    debtThreshold[i],
			LockerSavingRate: newLockerSavingRate,
			LotSize:          lotSize[i],
			BidFactor:        newBidFactor,
			DebtLotSize:      debtLotSize[i],
		})
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(lookupTableParams.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewLookupTableParamsProposal(lookupTableParams.Title, lookupTableParams.Description, LookupTableRecords)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}

// NewCmdAuctionControlProposal cmd to update controls for auction params.
func NewCmdAuctionControlProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction-control [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Set auction control",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewCreateAuctionControlParams(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetAuctionControlParamsMapping())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	return cmd
}

func NewCreateAuctionControlParams(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	auctionControlParams, err := parseAuctionControlParamsFlags(fs)

	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse auction control params: %w", err)
	}
	appID, err := strconv.ParseUint(auctionControlParams.AppID, 10, 64)
	if err != nil {
		return txf, nil, err
	}

	assetID, err := ParseUint64SliceFromString(auctionControlParams.AssetID, ",")
	if err != nil {
		return txf, nil, err
	}

	surplusAuction, err := ParseStringFromString(auctionControlParams.SurplusAuction, ",")
	if err != nil {
		return txf, nil, err
	}
	debtAuction, err := ParseStringFromString(auctionControlParams.DebtAuction, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOutOraclePrice, err := ParseStringFromString(auctionControlParams.AssetOutOraclePrice, ",")
	if err != nil {
		return txf, nil, err
	}

	assetOutPrice, err := ParseUint64SliceFromString(auctionControlParams.AssetOutPrice, ",")
	if err != nil {
		return txf, nil, err
	}

	var collectorAuctionLookupRecords types.CollectorAuctionLookupTable
	for i := range assetID {
		newSurplusAuction, err := strconv.ParseBool(surplusAuction[i])
		if err != nil {
			return txf, nil, err
		}
		newDebtAuction, err := strconv.ParseBool(debtAuction[i])
		if err != nil {
			return txf, nil, err
		}
		newAssetOutOraclePrice, err := strconv.ParseBool(assetOutOraclePrice[i])
		if err != nil {
			return txf, nil, err
		}
		collectorAuctionLookupRecords.AppId = appID
		var AssetIDToAuctionLookup types.AssetIdToAuctionLookupTable
		AssetIDToAuctionLookup.AssetId = assetID[i]
		AssetIDToAuctionLookup.IsSurplusAuction = newSurplusAuction
		AssetIDToAuctionLookup.IsDebtAuction = newDebtAuction
		AssetIDToAuctionLookup.AssetOutOraclePrice = newAssetOutOraclePrice
		AssetIDToAuctionLookup.AssetOutPrice = assetOutPrice[i]
		collectorAuctionLookupRecords.AssetIdToAuctionLookup = append(collectorAuctionLookupRecords.AssetIdToAuctionLookup, AssetIDToAuctionLookup)
	}

	from := clientCtx.GetFromAddress()

	deposit, err := sdk.ParseCoinsNormalized(auctionControlParams.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewAuctionLookupTableProposal(auctionControlParams.Title, auctionControlParams.Description, collectorAuctionLookupRecords)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}
	return txf, msg, nil
}
