package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type (
	XCreateWhitelistLiquidationMappingInputs createWhitelistLiquidationInputs
)

type XCreateWhitelistLiquidationInputsExceptions struct {
	XCreateWhitelistLiquidationMappingInputs
	Other *string // Other won't raise an error
}

func parseLiquidationWhitelistingFlags(fs *pflag.FlagSet) (*createWhitelistLiquidationInputs, error) {
	whitelistLiquidation := &createWhitelistLiquidationInputs{}
	whitelistLiquidationFile, _ := fs.GetString(FlagWhitelistLiquidation)

	if whitelistLiquidationFile == "" {
		return nil, fmt.Errorf("must pass in add asset mapping json using the --%s flag", whitelistLiquidationFile)
	}

	contents, err := os.ReadFile(whitelistLiquidationFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = whitelistLiquidation.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return whitelistLiquidation, nil
}

func (release *createWhitelistLiquidationInputs) UnmarshalJSON(data []byte) error {
	var createWhitelistLiquidationInputsE XCreateWhitelistLiquidationInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createWhitelistLiquidationInputsE); err != nil {
		return err
	}

	*release = createWhitelistLiquidationInputs(createWhitelistLiquidationInputsE.XCreateWhitelistLiquidationMappingInputs)
	return nil
}
