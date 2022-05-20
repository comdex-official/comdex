package cli

import (
	"fmt"
	"strconv"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/collector/types"
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

	// this line is used by starport scaffolding # 1

	return cmd
}

func NewCmdLookupTableParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collector-lookup-params [app-id] [collector_asset_id] [secondary_asset_id] [surplus_threshold] [debt_threshold] [locker_saving_rate] [lot_size] [bid_factor]",
		Args:  cobra.ExactArgs(8),
		Short: "Set collector lookup params for collector module",
		RunE: func(cmd *cobra.Command, args []string) error {

			appId, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			collector_asset_id, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}
			secondary_asset_id, err := ParseUint64SliceFromString(args[2], ",")
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
			lot_size, err := ParseUint64SliceFromString(args[6], ",")
			if err != nil {
				return err
			}
			bid_factor, err := ParseStringFromString(args[7], ",")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var LookupTableRecords []types.CollectorLookupTable
			for i := range appId {
				newbid_factor, _ := sdk.NewDecFromStr(bid_factor[i])
				newlockerSavingRate, _ := sdk.NewDecFromStr(lockerSavingRate[i])
				LookupTableRecords = append(LookupTableRecords, types.CollectorLookupTable{
					AppId: appId[i],
					CollectorAssetId: collector_asset_id[i],
					SecondaryAssetId: secondary_asset_id[i],
					SurplusThreshold: surplusThreshold[i],
					DebtThreshold: debtThreshold[i],
					LockerSavingRate: &newlockerSavingRate,
					LotSize: lot_size[i],
					BidFactor: &newbid_factor,
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

func NewCmdAuctionControlProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction-control [app-id] [surplus_auction] [debt_auction]",
		Args:  cobra.ExactArgs(3),
		Short: "Set auction control",
		RunE: func(cmd *cobra.Command, args []string) error {

			appId, err := ParseUint64SliceFromString(args[0], ",")
			if err != nil {
				return err
			}

			surplus_auction, err := ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}
			debt_auction, err := ParseStringFromString(args[2], ",")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var AuctionToAppTableRecords []types.AppIdToAuctionLookupTable
			for i := range appId {
				newSurplusAuction, err := strconv.ParseBool(surplus_auction[i])
				if err != nil {
					return err
				}
				newDebtAuction, err := strconv.ParseBool(debt_auction[i])
				if err != nil {
					return err
				}
				AuctionToAppTableRecords = append(AuctionToAppTableRecords, types.AppIdToAuctionLookupTable{
					AppId: appId[i],
					SurplusAuction: newSurplusAuction,
					DebtAuction: newDebtAuction,
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

			content := types.NewAuctionLookupTableProposal(title, description, AuctionToAppTableRecords)

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