package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
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
