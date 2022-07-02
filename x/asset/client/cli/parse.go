package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
)

type XCreateAddAssetMappingInputs createAddAssetMappingInputs

type XCreateAddAssetMappingInputsExceptions struct {
	XCreateAddAssetMappingInputs
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

func parseAssetMappingFlags(fs *pflag.FlagSet) (*createAddAssetMappingInputs, error) {
	assetMapping := &createAddAssetMappingInputs{}
	addAssetMappingFile, _ := fs.GetString(FlagAddAssetMappingFile)

	if addAssetMappingFile == "" {
		return nil, fmt.Errorf("must pass in add asset mapping json using the --%s flag", FlagAddAssetMappingFile)
	}

	contents, err := ioutil.ReadFile(addAssetMappingFile)
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
