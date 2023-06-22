package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type (
	XAuctionParamsMappingInputs createAuctionParamsInputs
)

type XAuctionParamsMappingInputsExceptions struct {
	XAuctionParamsMappingInputs
	Other *string // Other won't raise an error
}

func parseAuctionParamsFlags(fs *pflag.FlagSet) (*createAuctionParamsInputs, error) {
	auctionParams := &createAuctionParamsInputs{}
	addAuctionParamsFile, _ := fs.GetString(FlagAddAuctionParams)

	if addAuctionParamsFile == "" {
		return nil, fmt.Errorf("must pass in add asset mapping json using the --%s flag", addAuctionParamsFile)
	}

	contents, err := os.ReadFile(addAuctionParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = auctionParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return auctionParams, nil
}

func (release *createAuctionParamsInputs) UnmarshalJSON(data []byte) error {
	var addAuctionParamsInputsE XAuctionParamsMappingInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&addAuctionParamsInputsE); err != nil {
		return err
	}

	*release = createAuctionParamsInputs(addAuctionParamsInputsE.XAuctionParamsMappingInputs)
	return nil
}
