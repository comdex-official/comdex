package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type (
	XAddNewLendPairsInputs            addNewLendPairsInputs
	XAddLendPoolInputs                addLendPoolInputs
	XAddAssetRatesParamsInputs        addAssetRatesParamsInputs
	XSetAuctionParamsInputs           addNewAuctionParamsInputs
	XAddLendPoolPairsInputs           addLendPoolPairsInputs
	XAddAssetRatesLendPoolPairsInputs addAssetRatesPoolPairsInputs
	XAddDepreciatePoolsInputs         addDepreciatePoolsInputs
)

type XAddNewLendPairsInputsExceptions struct {
	XAddNewLendPairsInputs
	Other *string // Other won't raise an error
}

type XAddPoolInputsExceptions struct {
	XAddLendPoolInputs
	Other *string // Other won't raise an error
}

type XAddPoolPairsInputsExceptions struct {
	XAddLendPoolPairsInputs
	Other *string // Other won't raise an error
}

type XAddAssetRatesPoolPairsInputsExceptions struct {
	XAddAssetRatesLendPoolPairsInputs
	Other *string // Other won't raise an error
}

type XAddAssetRatesParamsInputsExceptions struct {
	XAddAssetRatesParamsInputs
	Other *string // Other won't raise an error
}

type XAddDepreciatePoolsInputsExceptions struct {
	XAddDepreciatePoolsInputs
	Other *string // Other won't raise an error
}

type XSetAuctionParamsInputsExceptions struct {
	XSetAuctionParamsInputs
	Other *string // Other won't raise an error
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *addNewLendPairsInputs) UnmarshalJSON(data []byte) error {
	var addNewLendPairsParamsE XAddNewLendPairsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addNewLendPairsParamsE); err != nil {
		return err
	}

	*release = addNewLendPairsInputs(addNewLendPairsParamsE.XAddNewLendPairsInputs)
	return nil
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *addLendPoolInputs) UnmarshalJSON(data []byte) error {
	var addPoolParamsE XAddPoolInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addPoolParamsE); err != nil {
		return err
	}

	*release = addLendPoolInputs(addPoolParamsE.XAddLendPoolInputs)
	return nil
}

func (release *addLendPoolPairsInputs) UnmarshalJSON(data []byte) error {
	var addPoolParamsE XAddPoolPairsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addPoolParamsE); err != nil {
		return err
	}

	*release = addLendPoolPairsInputs(addPoolParamsE.XAddLendPoolPairsInputs)
	return nil
}

func (release *addAssetRatesPoolPairsInputs) UnmarshalJSON(data []byte) error {
	var addPoolParamsE XAddAssetRatesPoolPairsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addPoolParamsE); err != nil {
		return err
	}

	*release = addAssetRatesPoolPairsInputs(addPoolParamsE.XAddAssetRatesLendPoolPairsInputs)
	return nil
}

func (release *addAssetRatesParamsInputs) UnmarshalJSON(data []byte) error {
	var addAssetRatesParamsE XAddAssetRatesParamsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addAssetRatesParamsE); err != nil {
		return err
	}

	*release = addAssetRatesParamsInputs(addAssetRatesParamsE.XAddAssetRatesParamsInputs)
	return nil
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *addNewAuctionParamsInputs) UnmarshalJSON(data []byte) error {
	var setAuctionParamsE XSetAuctionParamsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&setAuctionParamsE); err != nil {
		return err
	}

	*release = addNewAuctionParamsInputs(setAuctionParamsE.XSetAuctionParamsInputs)
	return nil
}

func parseAddNewLendPairsFlags(fs *pflag.FlagSet) (*addNewLendPairsInputs, error) {
	addLendPairsParams := &addNewLendPairsInputs{}
	addLendPairsParamsFile, _ := fs.GetString(FlagNewLendPairFile)

	if addLendPairsParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add new lend pairs json using the --%s flag", FlagNewLendPairFile)
	}

	contents, err := os.ReadFile(addLendPairsParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addLendPairsParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addLendPairsParams, nil
}

func parseAddPoolFlags(fs *pflag.FlagSet) (*addLendPoolInputs, error) {
	addPoolParams := &addLendPoolInputs{}
	addPoolParamsParamsFile, _ := fs.GetString(FlagAddLendPoolFile)

	if addPoolParamsParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add new pool json using the --%s flag", FlagAddLendPoolFile)
	}

	contents, err := os.ReadFile(addPoolParamsParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addPoolParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addPoolParams, nil
}

func parseAddPoolPairsFlags(fs *pflag.FlagSet) (*addLendPoolPairsInputs, error) {
	addPoolPairsParams := &addLendPoolPairsInputs{}
	addPoolPairsParamsFile, _ := fs.GetString(FlagAddLendPoolPairsFile)

	if addPoolPairsParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add new pool pairs json using the --%s flag", FlagAddLendPoolPairsFile)
	}

	contents, err := os.ReadFile(addPoolPairsParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addPoolPairsParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addPoolPairsParams, nil
}

func parseAddAssetratesPoolPairsFlags(fs *pflag.FlagSet) (*addAssetRatesPoolPairsInputs, error) {
	addPoolPairsParams := &addAssetRatesPoolPairsInputs{}
	addPoolPairsParamsFile, _ := fs.GetString(FlagAddAssetRatesPoolPairsFile)

	if addPoolPairsParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add new asset rates, pool & pairs json using the --%s flag", FlagAddAssetRatesPoolPairsFile)
	}

	contents, err := os.ReadFile(addPoolPairsParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addPoolPairsParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addPoolPairsParams, nil
}

func parseAssetRateStatsFlags(fs *pflag.FlagSet) (*addAssetRatesParamsInputs, error) {
	addAssetRatesParams := &addAssetRatesParamsInputs{}
	addAssetRatesParamsFile, _ := fs.GetString(FlagAddAssetRatesParamsFile)

	if addAssetRatesParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add asset rates stats json using the --%s flag", FlagAddAssetRatesParamsFile)
	}

	contents, err := os.ReadFile(addAssetRatesParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addAssetRatesParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addAssetRatesParams, nil
}

func parseAuctionPramsFlags(fs *pflag.FlagSet) (*addNewAuctionParamsInputs, error) {
	addNewAuctionParams := &addNewAuctionParamsInputs{}
	addNewAuctionParamsFile, _ := fs.GetString(FlagSetAuctionParamsFile)

	if addNewAuctionParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add auction params json using the --%s flag", FlagSetAuctionParamsFile)
	}

	contents, err := os.ReadFile(addNewAuctionParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addNewAuctionParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addNewAuctionParams, nil
}

func parseDepreciatePoolsFlags(fs *pflag.FlagSet) (*addDepreciatePoolsInputs, error) {
	depreciatePools := &addDepreciatePoolsInputs{}
	depreciatePoolsFile, _ := fs.GetString(FlagDepreciatePoolsFile)

	if depreciatePoolsFile == "" {
		return nil, fmt.Errorf("must pass in a appID and poolID for depreciate pool json using the --%s flag", FlagAddAssetRatesPoolPairsFile)
	}

	contents, err := os.ReadFile(depreciatePoolsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = depreciatePools.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return depreciatePools, nil
}

func (release *addDepreciatePoolsInputs) UnmarshalJSON(data []byte) error {
	var addDepreciatePoolsE XAddDepreciatePoolsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addDepreciatePoolsE); err != nil {
		return err
	}

	*release = addDepreciatePoolsInputs(addDepreciatePoolsE.XAddDepreciatePoolsInputs)
	return nil
}
