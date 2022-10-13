package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type (
	XCreateAddExternalLendRewardsMappingInputs createAddLendExternalRewardsInputs
)

type XCreateAddExternalLendRewardsMappingInputsExceptions struct {
	XCreateAddExternalLendRewardsMappingInputs
	Other *string // Other won't raise an error
}

func (release *createAddLendExternalRewardsInputs) UnmarshalJSON(data []byte) error {
	var createAddLendExternalRewardsMappingInputsE XCreateAddExternalLendRewardsMappingInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createAddLendExternalRewardsMappingInputsE); err != nil {
		return err
	}

	*release = createAddLendExternalRewardsInputs(createAddLendExternalRewardsMappingInputsE.XCreateAddExternalLendRewardsMappingInputs)
	return nil
}

func parseAddExternalLendRewardsMappingFlags(fs *pflag.FlagSet) (*createAddLendExternalRewardsInputs, error) {
	assetMapping := &createAddLendExternalRewardsInputs{}
	addAssetMappingFile, _ := fs.GetString(FlagAddLendExternalRewardsFile)

	if addAssetMappingFile == "" {
		return nil, fmt.Errorf("must pass in add External rewards lend mapping json using the --%s flag", FlagAddLendExternalRewardsFile)
	}

	contents, err := os.ReadFile(addAssetMappingFile)
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
