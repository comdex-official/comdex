package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/collector/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
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
	return cmd
}

// NewCmdLookupTableParams cmd for lookup table param proposal updates
func NewCmdLookupTableParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-params [app-id] [collector_asset_id] [secondary_asset_id] [surplus_threshold] [debt_threshold] [locker_saving_rate] [lot_size] [bid_factor] [debt_lot_size]",
		Args:  cobra.ExactArgs(9),
		Short: "Set collector lookup params for collector module",
		RunE: func(cmd *cobra.Command, args []string) error {

			appId, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			collectorAssetId, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			secondaryAssetId, err := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			surplusThreshold, err := ParseUint64SliceFromString(args[3], ",")
			if err != nil {
				return err
			}

			debtThreshold, err := ParseUint64SliceFromString(args[4], ",")
			if err != nil {
				return err
			}

			lockerSavingRate, err := ParseStringFromString(args[5], ",")
			if err != nil {
				return err
			}

			lotSize, err := ParseUint64SliceFromString(args[6], ",")
			if err != nil {
				return err
			}

			bidFactor, err := ParseStringFromString(args[7], ",")
			if err != nil {
				return err
			}
			debt_lot_size, err := ParseUint64SliceFromString(args[8], ",")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var LookupTableRecords []types.CollectorLookupTable
			for i := range appId {
				newBidFactor, err := sdk.NewDecFromStr(bidFactor[i])
				if err != nil {
					return err
				}

				newLockerSavingRate, err := sdk.NewDecFromStr(lockerSavingRate[i])
				if err != nil {
					return err
				}

				LookupTableRecords = append(LookupTableRecords, types.CollectorLookupTable{
					AppId:            appId[i],
					CollectorAssetId: collectorAssetId[i],
					SecondaryAssetId: secondaryAssetId[i],
					SurplusThreshold: surplusThreshold[i],
					DebtThreshold:    debtThreshold[i],
					LockerSavingRate: newLockerSavingRate,
					LotSize:          lotSize[i],
					BidFactor:        newBidFactor,
					DebtLotSize: debt_lot_size[i],
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

			content := types.NewLookupTableParamsProposal(title, description, LookupTableRecords)

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

// NewCmdAuctionControlProposal cmd to update controls for auction params
func NewCmdAuctionControlProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction-control [app-id] [asset_id] [surplus_auction] [debt_auction] [asset_out_oracle_price] [asset_out_price]",
		Args:  cobra.ExactArgs(6),
		Short: "Set auction control",
		RunE: func(cmd *cobra.Command, args []string) error {

			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			surplusAuction, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}
			debtAuction, err := ParseStringFromString(args[3], ",")
			if err != nil {
				return err
			}

			asset_out_oracle_price, err := ParseStringFromString(args[4], ",")
			if err != nil {
				return err
			}

			asset_out_price, err := ParseUint64SliceFromString(args[5], ",")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var collectorAuctionLookupRecords types.CollectorAuctionLookupTable
			for i := range assetId {
				newSurplusAuction, err := strconv.ParseBool(surplusAuction[i])
				if err != nil {
					return err
				}
				newDebtAuction, err := strconv.ParseBool(debtAuction[i])
				if err != nil {
					return err
				}
				newasset_out_oracle_price, err := strconv.ParseBool(asset_out_oracle_price[i])
				if err != nil {
					return err
				}
				collectorAuctionLookupRecords.AppId = appId
				var AssetIdToAuctionLookup types.AssetIdToAuctionLookupTable
				AssetIdToAuctionLookup.AssetId = assetId[i]
				AssetIdToAuctionLookup.IsSurplusAuction = newSurplusAuction
				AssetIdToAuctionLookup.IsDebtAuction = newDebtAuction
				AssetIdToAuctionLookup.AssetOutOraclePrice = newasset_out_oracle_price
				AssetIdToAuctionLookup.AssetOutPrice = asset_out_price[i]
				collectorAuctionLookupRecords.AssetIdToAuctionLookup = append(collectorAuctionLookupRecords.AssetIdToAuctionLookup, AssetIdToAuctionLookup)
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

			content := types.NewAuctionLookupTableProposal(title, description, collectorAuctionLookupRecords)

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
