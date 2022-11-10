package types

func NewGenesisState(vaults []Vault, stableMintVault []StableMintVault, appExtendedPairVaultMapping []AppExtendedPairVaultMappingData, userVaultAssetMapping []OwnerAppExtendedPairVaultMappingData, lengthOfvaults uint64) *GenesisState {
	return &GenesisState{
		Vaults:                      vaults,
		StableMintVault:             stableMintVault,
		AppExtendedPairVaultMapping: appExtendedPairVaultMapping,
		UserVaultAssetMapping:       userVaultAssetMapping,
		LengthOfVaults:              lengthOfvaults,
	}
}

func DefaultGenesisState() *GenesisState {
	var length uint64
	return NewGenesisState(
		[]Vault{},
		[]StableMintVault{},
		[]AppExtendedPairVaultMappingData{},
		[]OwnerAppExtendedPairVaultMappingData{},
		length,
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
