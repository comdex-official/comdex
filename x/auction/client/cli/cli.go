package cli

import (
	"strings"

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
		queryFilterDutchAuctions(),
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
	)

	return cmd
}

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}
