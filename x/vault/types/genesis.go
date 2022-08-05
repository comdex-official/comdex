package types

func NewGenesisState(vaults []Vault, stableMintVault []StableMintVault, appExtendedPairVaultMapping []AppExtendedPairVaultMapping, userVaultAssetMapping []UserVaultAssetMapping) *GenesisState {
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
		[]AppExtendedPairVaultMapping{},
		[]UserVaultAssetMapping{},
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
