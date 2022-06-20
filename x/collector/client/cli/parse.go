package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
)

type XCreateLookupTableParamsInputs createLookupTableParamsInputs
type XAuctionControlParamsInputs auctionControlParamsInputs

type XCreateLookupTableParamsInputsExceptions struct {
	XCreateLookupTableParamsInputs
	Other *string // Other won't raise an error
}

type XAuctionControlParamsInputsExceptions struct {
	XAuctionControlParamsInputs
	Other *string // Other won't raise an error
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *createLookupTableParamsInputs) UnmarshalJSON(data []byte) error {
	var createLookupTableParamsE XCreateLookupTableParamsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createLookupTableParamsE); err != nil {
		return err
	}

	*release = createLookupTableParamsInputs(createLookupTableParamsE.XCreateLookupTableParamsInputs)
	return nil
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *auctionControlParamsInputs) UnmarshalJSON(data []byte) error {
	var createAuctionControlParamsE XAuctionControlParamsInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createAuctionControlParamsE); err != nil {
		return err
	}

	*release = auctionControlParamsInputs(createAuctionControlParamsE.XAuctionControlParamsInputs)
	return nil
}

func parseLookupTableParamsFlags(fs *pflag.FlagSet) (*createLookupTableParamsInputs, error) {
	lookupTableParams := &createLookupTableParamsInputs{}
	lookupTableParamsFile, _ := fs.GetString(FlagAddLookupParamsTable)

	if lookupTableParamsFile == "" {
		return nil, fmt.Errorf("must pass in a lookup table params json using the --%s flag", FlagAddLookupParamsTable)
	}

	contents, err := ioutil.ReadFile(lookupTableParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = lookupTableParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return lookupTableParams, nil
}

func parseAuctionControlParamsFlags(fs *pflag.FlagSet) (*auctionControlParamsInputs, error) {
	auctionControlParams := &auctionControlParamsInputs{}
	auctionControlParamsFile, _ := fs.GetString(FlagAuctionControlParams)

	if auctionControlParamsFile == "" {
		return nil, fmt.Errorf("must pass in a auction control params json using the --%s flag", FlagAuctionControlParams)
	}

	contents, err := ioutil.ReadFile(auctionControlParamsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = auctionControlParams.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return auctionControlParams, nil
}
