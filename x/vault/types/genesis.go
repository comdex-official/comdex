package types

func NewGenesisState(vaults []Vault, stableMintVault []StableMintVault) *GenesisState {
	return &GenesisState{
		Vaults:          vaults,
		StableMintVault: stableMintVault,
		// AppExtendedPairVaultMapping: appExtendedPairVaultMapping,
		// UserVaultAssetMapping:       userVaultAssetMapping,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]Vault{},
		[]StableMintVault{},
		// []AppExtendedPairVaultMappingData{},
		// []UserVaultAssetMapping{},
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
