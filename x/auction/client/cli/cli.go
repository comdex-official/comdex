package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "auction",
		Short:                      "Auction module sub-commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryParams(),
		querySurplusAuction(),
		querySurplusAuctions(),
		querySurplusBiddings(),
		queryDebtAuction(),
		queryDebtAuctions(),
		queryDebtBidding(),
		queryDutchAuction(),
		queryDutchAuctions(),
		queryDutchBiddings(),
		queryProtocolStats(),
		queryAuctionParams(),
		queryDutchLendAuction(),
		queryDutchLendAuctions(),
		queryDutchLendBiddings(),
	)

	return cmd
}

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "auction",
		Short:                      "Auction module sub-commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txPlaceSurplusBid(),
		txPlaceDebtBid(),
		txPlaceDutchBid(),
		txPlaceDutchLendBid(),
	)

	return cmd
}
