package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
)

type (
	XCreateAddAssetMappingInputs  createAddAssetMappingInputs
	XCreateAddAssetsMappingInputs createAddAssetsMappingInputs
)

type XCreateAddAssetMappingInputsExceptions struct {
	XCreateAddAssetMappingInputs
	Other *string // Other won't raise an error
}

type XCreateAddAssetsMappingInputsExceptions struct {
	XCreateAddAssetsMappingInputs
	Other *string // Other won't raise an error
}

// UnmarshalJSON should error if there are fields unexpected.

func (release *createAddAssetMappingInputs) UnmarshalJSON(data []byte) error {
	var createAddAssetMappingInputsE XCreateAddAssetMappingInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createAddAssetMappingInputsE); err != nil {
		return err
	}

	*release = createAddAssetMappingInputs(createAddAssetMappingInputsE.XCreateAddAssetMappingInputs)
	return nil
}

func (release *createAddAssetsMappingInputs) UnmarshalJSON(data []byte) error {
	var createAddAssetsMappingInputsE XCreateAddAssetsMappingInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createAddAssetsMappingInputsE); err != nil {
		return err
	}

	*release = createAddAssetsMappingInputs(createAddAssetsMappingInputsE.XCreateAddAssetsMappingInputs)
	return nil
}

func parseAssetMappingFlags(fs *pflag.FlagSet) (*createAddAssetMappingInputs, error) {
	assetMapping := &createAddAssetMappingInputs{}
	addAssetMappingFile, _ := fs.GetString(FlagAddAssetMappingFile)

	if addAssetMappingFile == "" {
		return nil, fmt.Errorf("must pass in add asset mapping json using the --%s flag", FlagAddAssetMappingFile)
	}

	contents, err := ioutil.ReadFile(addAssetMappingFile) //nolint:gosec
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = assetMapping.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return assetMapping, nil
}

func parseAssetsMappingFlags(fs *pflag.FlagSet) (*createAddAssetsMappingInputs, error) {
	assetsMapping := &createAddAssetsMappingInputs{}
	addAssetsMappingFile, _ := fs.GetString(FlagAddAssetsMappingFile)

	if addAssetsMappingFile == "" {
		return nil, fmt.Errorf("must pass in add asset mapping json using the --%s flag", FlagAddAssetMappingFile)
	}

	contents, err := ioutil.ReadFile(addAssetsMappingFile) //nolint:gosec
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = assetsMapping.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return assetsMapping, nil
}
