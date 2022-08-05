package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
)

type XAddNewLendPairsInputs addNewLendPairsInputs
type XAddLendPoolInputs addLendPoolInputs
type XAddAssetRatesStatsInputs addAssetRatesStatsInputs
type XSetAuctionParamsInputs addNewAuctionParamsInputs

type XAddNewLendPairsInputsExceptions struct {
	XAddNewLendPairsInputs
	Other *string // Other won't raise an error
}

type XAddPoolInputsExceptions struct {
	XAddLendPoolInputs
	Other *string // Other won't raise an error
}
type XAddAssetRatesStatsInputsExceptions struct {
	XAddAssetRatesStatsInputs
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

func (release *addAssetRatesStatsInputs) UnmarshalJSON(data []byte) error {
	var addAssetRatesStatsE XAddAssetRatesStatsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addAssetRatesStatsE); err != nil {
		return err
	}

	*release = addAssetRatesStatsInputs(addAssetRatesStatsE.XAddAssetRatesStatsInputs)
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

	contents, err := ioutil.ReadFile(addLendPairsParamsFile) //nolint:gosec
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

	contents, err := ioutil.ReadFile(addPoolParamsParamsFile) //nolint:gosec
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

func parseAssetRateStatsFlags(fs *pflag.FlagSet) (*addAssetRatesStatsInputs, error) {
	addAssetRatesStats := &addAssetRatesStatsInputs{}
	addAssetRatesStatsFile, _ := fs.GetString(FlagAddAssetRatesStatsFile)

	if addAssetRatesStatsFile == "" {
		return nil, fmt.Errorf("must pass in a add asset rates stats json using the --%s flag", FlagAddAssetRatesStatsFile)
	}

	contents, err := ioutil.ReadFile(addAssetRatesStatsFile) //nolint:gosec
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = addAssetRatesStats.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return addAssetRatesStats, nil
}

func parseAuctionPramsFlags(fs *pflag.FlagSet) (*addNewAuctionParamsInputs, error) {
	addNewAuctionParams := &addNewAuctionParamsInputs{}
	addNewAuctionParamsFile, _ := fs.GetString(FlagSetAuctionParamsFile)

	if addNewAuctionParamsFile == "" {
		return nil, fmt.Errorf("must pass in a add auction params json using the --%s flag", FlagSetAuctionParamsFile)
	}

	contents, err := ioutil.ReadFile(addNewAuctionParamsFile) //nolint:gosec
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
