package types

import (
	fmt "fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MustMarshalTxGPIDS returns the TxGPIDS bytes.
// It throws panic if it fails.
func MustMarshalTxGPIDS(cdc codec.BinaryCodec, txGPIDS TxGPIDS) []byte {
	return cdc.MustMarshal(&txGPIDS)
}

// MustUnmarshalTxGPIDS return the unmarshalled TxGPIDS from bytes.
// It throws panic if it fails.
func MustUnmarshalTxGPIDS(cdc codec.BinaryCodec, value []byte) TxGPIDS {
	txGPIDS, err := UnmarshalTxGPIDS(cdc, value)
	if err != nil {
		panic(err)
	}

	return txGPIDS
}

// UnmarshalTxGPIDS returns the TxGPIDS from bytes.
func UnmarshalTxGPIDS(cdc codec.BinaryCodec, value []byte) (txGPIDS TxGPIDS, err error) {
	err = cdc.Unmarshal(value, &txGPIDS)
	return txGPIDS, err
}

// MustMarshalGasProvider returns the GasProvider bytes.
// It throws panic if it fails.
func MustMarshalGasProvider(cdc codec.BinaryCodec, gasProvider GasProvider) []byte {
	return cdc.MustMarshal(&gasProvider)
}

// MustUnmarshalGasProvider return the unmarshalled GasProvider from bytes.
// It throws panic if it fails.
func MustUnmarshalGasProvider(cdc codec.BinaryCodec, value []byte) GasProvider {
	gasProvider, err := UnmarshalGasProvider(cdc, value)
	if err != nil {
		panic(err)
	}

	return gasProvider
}

// UnmarshalGasProvider returns the GasProvider from bytes.
func UnmarshalGasProvider(cdc codec.BinaryCodec, value []byte) (gasProvider GasProvider, err error) {
	err = cdc.Unmarshal(value, &gasProvider)
	return gasProvider, err
}

// MustMarshalGasConsumer returns the GasConsumer bytes.
// It throws panic if it fails.
func MustMarshalGasConsumer(cdc codec.BinaryCodec, gasConsumer GasConsumer) []byte {
	return cdc.MustMarshal(&gasConsumer)
}

// MustUnmarshalGasConsumer return the unmarshalled GasConsumer from bytes.
// It throws panic if it fails.
func MustUnmarshalGasConsumer(cdc codec.BinaryCodec, value []byte) GasConsumer {
	gasConsumer, err := UnmarshalGasConsumer(cdc, value)
	if err != nil {
		panic(err)
	}

	return gasConsumer
}

// UnmarshalGasConsumer returns the GasConsumer from bytes.
func UnmarshalGasConsumer(cdc codec.BinaryCodec, value []byte) (gasConsumer GasConsumer, err error) {
	err = cdc.Unmarshal(value, &gasConsumer)
	return gasConsumer, err
}

func GasTankAddress(gasProviderID uint64) sdk.AccAddress {
	return DeriveAddress(
		AddressType32Bytes,
		ModuleName,
		strings.Join([]string{GasTankAddressPrefix, strconv.FormatUint(gasProviderID, 10)}, ModuleAddressNameSplitter))
}

func NewGasProvider(
	id uint64,
	creator sdk.AccAddress,
	maxTxsCountPerConsumer uint64,
	maxFeeUsagePerConsumer sdkmath.Int,
	maxFeeUsagePerTx sdkmath.Int,
	txsAllowed []string,
	contractsAllowed []string,
	feeDenom string,
) GasProvider {
	return GasProvider{
		Id:                     id,
		Creator:                creator.String(),
		GasTank:                GasTankAddress(id).String(),
		IsActive:               true,
		MaxTxsCountPerConsumer: maxTxsCountPerConsumer,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		TxsAllowed:             RemoveDuplicates(txsAllowed),
		ContractsAllowed:       RemoveDuplicates(contractsAllowed),
		AuthorizedActors:       []string{},
		FeeDenom:               feeDenom,
	}
}

func (gasProvider GasProvider) GetGasTankReserveAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(gasProvider.GasTank)
	if err != nil {
		panic(err)
	}
	return addr
}

func (gasProvider GasProvider) Validate() error {
	if gasProvider.Id == 0 {
		return fmt.Errorf("pair id must not be 0")
	}
	if err := sdk.ValidateDenom(gasProvider.FeeDenom); err != nil {
		return fmt.Errorf("invalid fee denom: %w", err)
	}
	if gasProvider.MaxTxsCountPerConsumer == 0 {
		return fmt.Errorf("max tx count per consumer must not be 0")
	}
	if !gasProvider.MaxFeeUsagePerTx.IsPositive() {
		return fmt.Errorf("max_fee_usage_per_tx should be positive")
	}
	if !gasProvider.MaxFeeUsagePerConsumer.IsPositive() {
		return fmt.Errorf("max_fee_usage_per_consumer should be positive")
	}
	if len(gasProvider.TxsAllowed) == 0 && len(gasProvider.ContractsAllowed) == 0 {
		return fmt.Errorf("atleast one tx or contract is required to initialize")
	}

	return nil
}
