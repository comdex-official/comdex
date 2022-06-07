package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
)

type XCreateExtPairVaultInputs createExtPairVaultInputs
type XCreateAddAssetMappingInputs createAddAssetMappingInputs
type XCreateAddWhiteListedPairInputs createAddWhiteListedPairsInputs

type XCreateExtPairVaultInputsExceptions struct {
	XCreateExtPairVaultInputs
	Other *string // Other won't raise an error
}

type XCreateAddAssetMappingInputsExceptions struct {
	XCreateAddAssetMappingInputs
	Other *string // Other won't raise an error
}

type XCreateAddWhiteListedPairInputsExceptions struct {
	XCreateAddWhiteListedPairInputs
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

func (release *createAddWhiteListedPairsInputs) UnmarshalJSON(data []byte) error {
	var createAddWhiteListedPairInputsE XCreateAddWhiteListedPairInputsExceptions
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&createAddWhiteListedPairInputsE); err != nil {
		return err
	}

	*release = createAddWhiteListedPairsInputs(createAddWhiteListedPairInputsE.XCreateAddWhiteListedPairInputs)
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

func parseWhiteListedPairsFlags(fs *pflag.FlagSet) (*createAddWhiteListedPairsInputs, error) {
	whiteListedAssetPairs := &createAddWhiteListedPairsInputs{}
	whiteListedAssetPairsFile, _ := fs.GetString(FlagAddWhiteListedPairsFile)

	if whiteListedAssetPairsFile == "" {
		return nil, fmt.Errorf("must pass in add white listed pairs json using the --%s flag", FlagAddWhiteListedPairsFile)
	}

	contents, err := ioutil.ReadFile(whiteListedAssetPairsFile)
	if err != nil {
		return nil, err
	}

	// make exception if unknown field exists
	err = whiteListedAssetPairs.UnmarshalJSON(contents)
	if err != nil {
		return nil, err
	}

	return whiteListedAssetPairs, nil
}
