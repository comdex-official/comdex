package params

import (
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Name defines the application name of the Comdex network.
	Name = "comdex"

	// BondDenom defines the native staking token denomination.
	BondDenom = "ucmdx"

	// DisplayDenom defines the name, symbol, and display value of the cmdx token.
	DisplayDenom = "CMDX"

	// DefaultGasLimit - set to the same value as cosmos-sdk flags.DefaultGasLimit
	// this value is currently only used in tests.
	DefaultGasLimit = 200000
)

var (
	// ProtocolMinGasPrice is a consensus controlled gas price. Each validator must set his
	// `minimum-gas-prices` in app.toml config to value above ProtocolMinGasPrice.
	// Transactions with gas-price smaller than ProtocolMinGasPrice will fail during DeliverTx.
	ProtocolMinGasPrice = sdk.NewDecCoinFromDec(BondDenom, sdk.MustNewDecFromStr("0.00"))

	// DefaultGovPeriod is 3 days. This should be long enough for validators to react,
	// and short enough for the team to list new assets competitively.
	DefaultGovPeriod = time.Hour * 24 * 3
)

// SetAddressPrefixes call this in init to set sdk config's bech32 address prefix
func SetAddressPrefixes() {
	SetAddressConfig()
	if AccountAddressPrefix != Name {
		log.Fatal("AccountAddresPrefix must equal Name")
	}
}
