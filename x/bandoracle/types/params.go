package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyCreator        = []byte("Creator")
	KeyOracleScriptId = []byte("OracleScriptId")
	KeySourceChannel  = []byte("SourceChannel")
	KeyAskCount       = []byte("AskCount")
	KeyMinCount       = []byte("MinCount")
	KeyFeeLimit       = []byte("FeeLimit")
	KeyPrepareGas     = []byte("PrepareGas")
	KeyExecuteGas     = []byte("ExecuteGas")
)

var (
	DefaultCreator        = ModuleName
	DefaultOracleScriptId = "112"
	DefaultSourceChannel  = "channel-0"
	DefaultAskCount       = "5"
	DefaultMinCount       = "3"
	DefaultFeeLimit       = sdk.Coins{sdk.NewCoin("uband", sdk.NewInt(250000))}
	DefaultPrepareGas     = "6000"
	DefaultExecuteGas     = "6000"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		Creator:        DefaultCreator,
		OracleScriptId: DefaultOracleScriptId,
		SourceChannel:  DefaultSourceChannel,
		AskCount:       DefaultAskCount,
		MinCount:       DefaultMinCount,
		FeeLimit:       DefaultFeeLimit,
		PrepareGas:     DefaultPrepareGas,
		ExecuteGas:     DefaultExecuteGas,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCreator, &p.Creator, ValidateCreator),
		paramtypes.NewParamSetPair(KeyOracleScriptId, &p.OracleScriptId, ValidateOracleScriptID),
		paramtypes.NewParamSetPair(KeySourceChannel, &p.SourceChannel, ValidateSourceChannel),
		paramtypes.NewParamSetPair(KeyAskCount, &p.AskCount, ValidateAskCount),
		paramtypes.NewParamSetPair(KeyMinCount, &p.MinCount, ValidateMinCount),
		paramtypes.NewParamSetPair(KeyFeeLimit, &p.FeeLimit, ValidateFeeLimit),
		paramtypes.NewParamSetPair(KeyPrepareGas, &p.PrepareGas, ValidatePrepareGas),
		paramtypes.NewParamSetPair(KeyExecuteGas, &p.ExecuteGas, ValidateExecuteGas),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

func ValidateCreator(i interface{}) error {
	return nil
}

func ValidateOracleScriptID(i interface{}) error {
	return nil
}

func ValidateSourceChannel(i interface{}) error {
	return nil
}

func ValidateAskCount(i interface{}) error {
	return nil
}

func ValidateMinCount(i interface{}) error {
	return nil
}

func ValidateFeeLimit(i interface{}) error {
	return nil
}

func ValidatePrepareGas(i interface{}) error {
	return nil
}

func ValidateExecuteGas(i interface{}) error {
	return nil
}
