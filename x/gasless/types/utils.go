package types

import (
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive private plan gas tank address.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return address.Module(moduleName, []byte(name))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}

// ItemExists returns true if item exists in array else false .
func ItemExists(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

func NewGasProviderResponse(gasProvider GasProvider, balances sdk.Coins) GasProviderResponse {
	return GasProviderResponse{
		Id:                     gasProvider.Id,
		Creator:                gasProvider.Creator,
		GasTankAddress:         gasProvider.GasTank,
		GasTankBalances:        balances,
		IsActive:               gasProvider.IsActive,
		MaxTxsCountPerConsumer: gasProvider.MaxTxsCountPerConsumer,
		MaxFeeUsagePerConsumer: gasProvider.MaxFeeUsagePerConsumer,
		MaxFeeUsagePerTx:       gasProvider.MaxFeeUsagePerTx,
		TxsAllowed:             gasProvider.TxsAllowed,
		ContractsAllowed:       gasProvider.ContractsAllowed,
		AuthorizedActors:       gasProvider.AuthorizedActors,
		FeeDenom:               gasProvider.FeeDenom,
	}
}
