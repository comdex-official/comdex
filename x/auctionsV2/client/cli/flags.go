package cli

import flag "github.com/spf13/pflag"

const (
	FlagAddAuctionParams = "add-auction-params"
)

func FlagSetAuctionParams() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAuctionParams, "", "add auction params json file path")
	return fs
}

type createAuctionParamsInputs struct {
	AuctionDurationSeconds string `json:"auction_duration_seconds"`
	Step                   string `json:"step"`
	WithdrawalFee          string `json:"withdrawal_fee"`
	ClosingFee             string `json:"closing_fee"`
	MinUsdValueLeft        string `json:"min_usd_value_left"`
	BidFactor              string `json:"bid_factor"`
	LiquidationPenalty     string `json:"liquidation_penalty"`
	AuctionBonus           string `json:"auction_bonus"`
	Title                  string
	Description            string
	Deposit                string
}
