package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/expected"
)

// BulkSendCoinsOperation holds a list of SendCoins operations for bulk execution.
type BulkSendCoinsOperation struct {
	Inputs  []banktypes.Input
	Outputs []banktypes.Output
}

// NewBulkSendCoinsOperation returns an empty BulkSendCoinsOperation.
func NewBulkSendCoinsOperation() *BulkSendCoinsOperation {
	return &BulkSendCoinsOperation{
		Inputs:  []banktypes.Input{},
		Outputs: []banktypes.Output{},
	}
}

// QueueSendCoins queues a BankKeeper.SendCoins operation for later execution.
func (op *BulkSendCoinsOperation) QueueSendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
	if amt.IsValid() && !amt.IsZero() {
		op.Inputs = append(op.Inputs, banktypes.NewInput(fromAddr, amt))
		op.Outputs = append(op.Outputs, banktypes.NewOutput(toAddr, amt))
	}
}

// Run runs BankKeeper.InputOutputCoins once for queued operations.
func (op *BulkSendCoinsOperation) Run(ctx sdk.Context, bankKeeper expected.BankKeeper) error {
	if len(op.Inputs) > 0 && len(op.Outputs) > 0 {
		return bankKeeper.InputOutputCoins(ctx, op.Inputs, op.Outputs)
	}
	return nil
}

// IsTooSmallOrderAmount returns whether the order amount is too small for
// matching, based on the order price.
func IsTooSmallOrderAmount(amt sdk.Int, price sdk.Dec) bool {
	return amt.LT(amm.MinCoinAmount) || price.MulInt(amt).LT(amm.MinCoinAmount.ToDec())
}

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive private plan farming pool address, and staking reserve address.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return sdk.AccAddress(address.Module(moduleName, []byte(name)))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}
