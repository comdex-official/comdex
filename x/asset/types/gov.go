package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalAddAssets                 = "AddAssets"
	ProposalUpdateAsset               = "UpdateAsset"
	ProposalAddPairs                  = "AddPairs"
	ProposalAddAppMapping             = "AddAppMapping"
	ProposalAddAssetMapping           = "AddAssetMapping"
	ProposalAddExtendedPairsVault     = "AddExtendedPairsVault"
	ProposalUpdateGovTimeInAppMapping = "UpdateGovTimeInAppMapping"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddAssets)
	govtypes.RegisterProposalTypeCodec(&AddAssetsProposal{}, "comdex/AddAssetsProposal")

	govtypes.RegisterProposalType(ProposalUpdateAsset)
	govtypes.RegisterProposalTypeCodec(&UpdateAssetProposal{}, "comdex/UpdateAssetProposal")

	govtypes.RegisterProposalType(ProposalAddPairs)
	govtypes.RegisterProposalTypeCodec(&AddPairsProposal{}, "comdex/AddPairsProposal")

	govtypes.RegisterProposalType(ProposalUpdateGovTimeInAppMapping)
	govtypes.RegisterProposalTypeCodec(&UpdateGovTimeInAppMappingProposal{}, "comdex/UpdateGovTimeInAppMappingProposal")

	govtypes.RegisterProposalType(ProposalAddAppMapping)
	govtypes.RegisterProposalTypeCodec(&AddAppMappingProposal{}, "comdex/AddAppMappingProposal")

	govtypes.RegisterProposalType(ProposalAddAssetMapping)
	govtypes.RegisterProposalTypeCodec(&AddAssetMappingProposal{}, "comdex/AddAssetMappingProposal")

	govtypes.RegisterProposalType(ProposalAddExtendedPairsVault)
	govtypes.RegisterProposalTypeCodec(&AddExtendedPairsVaultProposal{}, "comdex/AddExtendedPairsVaultProposal")
}

var (
	_ govtypes.Content = &AddAssetsProposal{}
	_ govtypes.Content = &UpdateAssetProposal{}
	_ govtypes.Content = &AddPairsProposal{}
	_ govtypes.Content = &UpdateGovTimeInAppMappingProposal{}
	_ govtypes.Content = &AddAppMappingProposal{}
	_ govtypes.Content = &AddAssetMappingProposal{}
	_ govtypes.Content = &AddExtendedPairsVaultProposal{}
)

func NewAddAssetsProposal(title, description string, assets []Asset) govtypes.Content {
	return &AddAssetsProposal{
		Title:       title,
		Description: description,
		Assets:      assets,
	}
}
func (p *AddAssetsProposal) GetTitle() string {
	return p.Title
}

func (p *AddAssetsProposal) GetDescription() string {
	return p.Description
}
func (p *AddAssetsProposal) ProposalRoute() string { return RouterKey }

func (p *AddAssetsProposal) ProposalType() string { return ProposalAddAssets }

func (p *AddAssetsProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.Assets) == 0 {
		return ErrorEmptyProposalAssets
	}

	for _, asset := range p.Assets {
		if err := asset.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewUpdateAssetProposal(title, description string, asset Asset) govtypes.Content {
	return &UpdateAssetProposal{
		Title:       title,
		Description: description,
		Asset:       asset,
	}
}

func (p *UpdateAssetProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateAssetProposal) GetDescription() string {
	return p.Description
}

func (p *UpdateAssetProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateAssetProposal) ProposalType() string { return ProposalUpdateAsset }

func (p *UpdateAssetProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	asset := p.Asset
	if err := asset.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddPairsProposal(title, description string, pairs []Pair) govtypes.Content {
	return &AddPairsProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *AddPairsProposal) GetTitle() string {
	return p.Title
}

func (p *AddPairsProposal) GetDescription() string {
	return p.Description
}
func (p *AddPairsProposal) ProposalRoute() string { return RouterKey }

func (p *AddPairsProposal) ProposalType() string { return ProposalAddPairs }

func (p *AddPairsProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.Pairs) == 0 {
		return ErrorEmptyProposalAssets
	}

	for _, pair := range p.Pairs {
		if err := pair.Validate(); err != nil {
			return err
		}
	}

	return nil
}


func NewAddAppMapingProposa(title, description string, amap []AppMapping) govtypes.Content {
	return &AddAppMappingProposal{
		Title:       title,
		Description: description,
		App:         amap,
	}
}

func (p *AddAppMappingProposal) GetTitle() string {
	return p.Title
}

func (p *AddAppMappingProposal) GetDescription() string {
	return p.Description
}

func (p *AddAppMappingProposal) ProposalRoute() string { return RouterKey }

func (p *AddAppMappingProposal) ProposalType() string { return ProposalAddAppMapping }

func (p *AddAppMappingProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.App) == 0 {
		return ErrorEmptyProposalAssets
	}

	return nil
}

func NewUpdateGovTimeInAppMappingProposal(title, description string, aTime AppAndGovTime) govtypes.Content {
	return &UpdateGovTimeInAppMappingProposal{
		Title:       title,
		Description: description,
		GovTime:     aTime,
	}
}

func (p *UpdateGovTimeInAppMappingProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateGovTimeInAppMappingProposal) GetDescription() string {
	return p.Description
}

func (p *UpdateGovTimeInAppMappingProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateGovTimeInAppMappingProposal) ProposalType() string {
	return ProposalUpdateGovTimeInAppMapping
}

func (p *UpdateGovTimeInAppMappingProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewAddAssetMapingProposa(title, description string, amap []AppMapping) govtypes.Content {
	return &AddAssetMappingProposal{
		Title:       title,
		Description: description,
		App:         amap,
	}
}

func (p *AddAssetMappingProposal) GetTitle() string {
	return p.Title
}

func (p *AddAssetMappingProposal) GetDescription() string {
	return p.Description
}

func (p *AddAssetMappingProposal) ProposalRoute() string { return RouterKey }

func (p *AddAssetMappingProposal) ProposalType() string { return ProposalAddAssetMapping }

func (p *AddAssetMappingProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.App) == 0 {
		return ErrorEmptyProposalAssets
	}

	return nil
}

func NewAddExtendedPairsVaultProposa(title, description string, pairs []ExtendedPairVault) govtypes.Content {
	return &AddExtendedPairsVaultProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *AddExtendedPairsVaultProposal) GetTitle() string {
	return p.Title
}

func (p *AddExtendedPairsVaultProposal) GetDescription() string {
	return p.Description
}

func (p *AddExtendedPairsVaultProposal) ProposalRoute() string { return RouterKey }

func (p *AddExtendedPairsVaultProposal) ProposalType() string { return ProposalAddExtendedPairsVault }

func (p *AddExtendedPairsVaultProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.Pairs) == 0 {
		return ErrorEmptyProposalAssets
	}

	return nil
}
