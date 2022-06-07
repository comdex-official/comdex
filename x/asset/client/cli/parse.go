package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
)

type XCreateExtPairVaultInputs createExtPairVaultInputs

type XCreateExtPairVaultInputsExceptions struct {
	XCreateExtPairVaultInputs
	Other *string // Other won't raise an error
}

// UnmarshalJSON should error if there are fields unexpected.
func (release *createExtPairVaultInputs) UnmarshalJSON(data []byte) error {
	var createExtendedPairVaultE XCreateExtPairVaultInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createExtendedPairVaultE); err != nil {
		return err
	}

	*release = createExtPairVaultInputs(createExtendedPairVaultE.XCreateExtPairVaultInputs)
	return nil
}

func parseExtendPairVaultFlags(fs *pflag.FlagSet) (*createExtPairVaultInputs, error) {
	extPairVault := &createExtPairVaultInputs{}
	extPairVaultFile, _ := fs.GetString(FlagExtendedPairVaultFile)

	if extPairVaultFile == "" {
		return nil, fmt.Errorf("must pass in a Extended Pair Vault json using the --%s flag", FlagExtendedPairVaultFile)
	}

	contents, err := ioutil.ReadFile(extPairVaultFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = extPairVault.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return extPairVault, nil
}
