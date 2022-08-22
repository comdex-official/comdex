package types

func NewGenesisState(vaults []Vault, stableMintVault []StableMintVault, appExtendedPairVaultMapping []AppExtendedPairVaultMappingData, userVaultAssetMapping []OwnerAppExtendedPairVaultMappingData) *GenesisState {
	return &GenesisState{
		Vaults:                      vaults,
		StableMintVault:             stableMintVault,
		AppExtendedPairVaultMapping: appExtendedPairVaultMapping,
		UserVaultAssetMapping:       userVaultAssetMapping,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]Vault{},
		[]StableMintVault{},
		[]AppExtendedPairVaultMappingData{},
		[]OwnerAppExtendedPairVaultMappingData{},
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
