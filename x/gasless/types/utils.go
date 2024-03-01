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

func RemoveDuplicates(input []string) []string {
	uniqueMap := make(map[string]bool)
	for _, str := range input {
		uniqueMap[str] = true
	}
	uniqueSlice := make([]string, 0, len(uniqueMap))
	for str := range uniqueMap {
		uniqueSlice = append(uniqueSlice, str)
	}
	return uniqueSlice
}

func RemoveDuplicatesUint64(input []uint64) []uint64 {
	uniqueMap := make(map[uint64]bool)
	var uniqueList []uint64
	for _, v := range input {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			uniqueList = append(uniqueList, v)
		}
	}
	return uniqueList
}

func RemoveValueFromListUint64(list []uint64, x uint64) []uint64 {
	var newList []uint64
	for _, v := range list {
		if v != x {
			newList = append(newList, v)
		}
	}
	return newList
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

func NewConsumptionDetail(
	txsAllowed uint64,
	feeConsumptionAllowed sdk.Coin,
) *ConsumptionDetail {
	return &ConsumptionDetail{
		IsBlocked:                  false,
		TotalTxsAllowed:            txsAllowed,
		TotalTxsMade:               0,
		TotalFeeConsumptionAllowed: feeConsumptionAllowed,
		TotalFeesConsumed:          sdk.NewCoin(feeConsumptionAllowed.Denom, sdk.ZeroInt()),
		Usage: &Usage{
			Txs:       make(map[string]*UsageDetails),
			Contracts: make(map[string]*UsageDetails),
		},
	}
}

func NewTxGPIDS(tpoc string) TxGPIDS {
	return TxGPIDS{
		TxPathOrContractAddress: tpoc,
		GasProviderIds:          []uint64{},
	}
}
