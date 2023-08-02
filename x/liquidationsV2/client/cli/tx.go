package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	flag "github.com/spf13/pflag"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
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

	cmd.AddCommand(
		txLiquidateInternalKeeper(),
		txLiquidateExternalKeeper(),
		txAppReserveFunds(),
	)
	return cmd
}

func txLiquidateInternalKeeper() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidate-internal-keeper [type] [id]",
		Short: "liquidate faulty positions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liqType, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidateInternalKeeperRequest(ctx.FromAddress, liqType, id)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txAppReserveFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app-reserve-funds [app-id] [asset-id] [amount]",
		Short: "app reserve funds",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgAppReserveFundsRequest(ctx.GetFromAddress().String(), appId, assetId, amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txLiquidateExternalKeeper() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidate-external-keeper [app-id] [owner] [collateral-token] [debt-token] [collateral-asset-id] [debt-asset-id] [is-debt-cmst]",
		Short: "liquidate faulty positions - external apps",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			AppId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			Owner := args[1]
			CollateralToken, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}
			DebtToken, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}
			CollateralAssetId, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			DebtAssetId, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}
			IsDebtCmst := ParseBoolFromString(args[6])

			msg := types.NewMsgLiquidateExternalKeeperRequest(ctx.FromAddress, AppId, Owner, CollateralToken, DebtToken, CollateralAssetId, DebtAssetId, IsDebtCmst)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func NewCmdSubmitWhitelistingLiquidationProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidation-whitelisting [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Submit liquidation whitelisting proposal",
		Long:  `Must provide path to a add liquidation whitelisting in JSON file (--add-liquidation-whitelisting) describing the it in an app to be created`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err !=nil {
				return err
			}
			txf = txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := WhitelistLiquidation(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetWhitelistLiquidation())
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func WhitelistLiquidation(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	liquidationWhitelisting, err := parseLiquidationWhitelistingFlags(fs)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse liquidationWhitelisting: %w", err)
	}
	from := clientCtx.GetFromAddress()

	appId, err := strconv.ParseUint(liquidationWhitelisting.AppId, 10, 64)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse appId: %w", err)
	}

	premium, err := sdk.NewDecFromStr(liquidationWhitelisting.Premium)
	if err != nil {
		return txf, nil, err
	}

	discount, err := sdk.NewDecFromStr(liquidationWhitelisting.Discount)
	if err != nil {
		return txf, nil, err
	}

	decrementFactor, err := strconv.ParseInt(liquidationWhitelisting.DecrementFactor, 10, 64)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse decrementFactor: %w", err)
	}

	decrementFactorEnglish, err := strconv.ParseInt(liquidationWhitelisting.DecrementFactorEnglish, 10, 64)
	if err != nil {
		return txf, nil, fmt.Errorf("failed to parse decrementFactor: %w", err)
	}

	dutchAuctionParam := types.DutchAuctionParam{
		Premium:         premium,
		Discount:        discount,
		DecrementFactor: sdk.NewInt(decrementFactor),
	}
	englishAuctionParam := types.EnglishAuctionParam{DecrementFactor: sdk.NewInt(decrementFactorEnglish)}

	liquidationWhitelistingStruct := types.LiquidationWhiteListing{
		AppId:               appId,
		Initiator:           ParseBoolFromString(liquidationWhitelisting.Initiator),
		IsDutchActivated:    ParseBoolFromString(liquidationWhitelisting.IsDutchActivated),
		DutchAuctionParam:   &dutchAuctionParam,
		IsEnglishActivated:  ParseBoolFromString(liquidationWhitelisting.IsEnglishActivated),
		EnglishAuctionParam: &englishAuctionParam,
	}

	deposit, err := sdk.ParseCoinsNormalized(liquidationWhitelisting.Deposit)
	if err != nil {
		return txf, nil, err
	}

	content := types.NewLiquidationWhiteListingProposal(liquidationWhitelisting.Title, liquidationWhitelisting.Description, liquidationWhitelistingStruct)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return txf, nil, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, msg, nil
}
