package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"strings"
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
	if err := ValidateCreator(p.Creator); err != nil {
		return err
	}
	if err := ValidateOracleScriptID(p.OracleScriptId); err != nil {
		return err
	}
	if err := ValidateSourceChannel(p.SourceChannel); err != nil {
		return err
	}
	if err := ValidateAskCount(p.AskCount); err != nil {
		return err
	}
	if err := ValidateMinCount(p.MinCount); err != nil {
		return err
	}
	if err := ValidateFeeLimit(p.FeeLimit); err != nil {
		return err
	}
	if err := ValidatePrepareGas(p.PrepareGas); err != nil {
		return err
	}
	if err := ValidateExecuteGas(p.ExecuteGas); err != nil {
		return err
	}
	return nil
}

func ValidateCreator(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("creator cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateOracleScriptID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("oraclesccriptid cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateSourceChannel(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("sourcechannel cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateAskCount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("askcount cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateMinCount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mincount cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateFeeLimit(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("feelimit creation fee: %+v", i)
	}

	return nil
}

func ValidatePrepareGas(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("preparegas cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidateExecuteGas(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("executegas cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
