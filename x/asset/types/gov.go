package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalAddAssets          = "AddAssets"
	ProposalUpdateAsset        = "UpdateAsset"
	ProposalAddPairs           = "AddPairs"
	ProposalUpdatePair         = "UpdatePair"
	ProposalAddApp             = "AddApp"
	ProposalAddAssetInApp      = "AddAssetInApp"
	ProposalUpdateGovTimeInApp = "UpdateGovTimeInApp"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddAssets)
	govtypes.RegisterProposalTypeCodec(&AddAssetsProposal{}, "comdex/AddAssetsProposal")

	govtypes.RegisterProposalType(ProposalUpdateAsset)
	govtypes.RegisterProposalTypeCodec(&UpdateAssetProposal{}, "comdex/UpdateAssetProposal")

	govtypes.RegisterProposalType(ProposalAddPairs)
	govtypes.RegisterProposalTypeCodec(&AddPairsProposal{}, "comdex/AddPairsProposal")

	govtypes.RegisterProposalType(ProposalUpdatePair)
	govtypes.RegisterProposalTypeCodec(&UpdatePairProposal{}, "comdex/UpdatePairProposal")

	govtypes.RegisterProposalType(ProposalUpdateGovTimeInApp)
	govtypes.RegisterProposalTypeCodec(&UpdateGovTimeInAppProposal{}, "comdex/UpdateGovTimeInAppProposal")

	govtypes.RegisterProposalType(ProposalAddApp)
	govtypes.RegisterProposalTypeCodec(&AddAppProposal{}, "comdex/AddAppProposal")

	govtypes.RegisterProposalType(ProposalAddAssetInApp)
	govtypes.RegisterProposalTypeCodec(&AddAssetInAppProposal{}, "comdex/AddAssetInAppProposal")
}

var (
	_ govtypes.Content = &AddAssetsProposal{}
	_ govtypes.Content = &UpdateAssetProposal{}
	_ govtypes.Content = &AddPairsProposal{}
	_ govtypes.Content = &UpdatePairProposal{}
	_ govtypes.Content = &UpdateGovTimeInAppProposal{}
	_ govtypes.Content = &AddAppProposal{}
	_ govtypes.Content = &AddAssetInAppProposal{}
)

func NewAddAssetsProposal(title, description string, assets Asset) govtypes.Content {
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

	if err := p.Assets.Validate(); err != nil {
		return err
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

	if err := p.Asset.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddPairsProposal(title, description string, pairs Pair) govtypes.Content {
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

	if err := p.Pairs.Validate(); err != nil {
		return err
	}

	return nil
}

func NewUpdatePairProposal(title, description string, pair Pair) govtypes.Content {
	return &UpdatePairProposal{
		Title:       title,
		Description: description,
		Pairs:       pair,
	}
}

func (p *UpdatePairProposal) GetTitle() string {
	return p.Title
}

func (p *UpdatePairProposal) GetDescription() string {
	return p.Description
}

func (p *UpdatePairProposal) ProposalRoute() string { return RouterKey }

func (p *UpdatePairProposal) ProposalType() string { return ProposalUpdatePair }

func (p *UpdatePairProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pair := p.Pairs
	if err := pair.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddAppProposal(title, description string, amap AppData) govtypes.Content {
	return &AddAppProposal{
		Title:       title,
		Description: description,
		App:         amap,
	}
}

func (p *AddAppProposal) GetTitle() string {
	return p.Title
}

func (p *AddAppProposal) GetDescription() string {
	return p.Description
}

func (p *AddAppProposal) ProposalRoute() string { return RouterKey }

func (p *AddAppProposal) ProposalType() string { return ProposalAddApp }

func (p *AddAppProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewUpdateGovTimeInAppProposal(title, description string, aTime AppAndGovTime) govtypes.Content {
	return &UpdateGovTimeInAppProposal{
		Title:       title,
		Description: description,
		GovTime:     aTime,
	}
}

func (p *UpdateGovTimeInAppProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateGovTimeInAppProposal) GetDescription() string {
	return p.Description
}

func (p *UpdateGovTimeInAppProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateGovTimeInAppProposal) ProposalType() string {
	return ProposalUpdateGovTimeInApp
}

func (p *UpdateGovTimeInAppProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewAddAssetInAppProposal(title, description string, amap AppData) govtypes.Content {
	return &AddAssetInAppProposal{
		Title:       title,
		Description: description,
		App:         amap,
	}
}

func (p *AddAssetInAppProposal) GetTitle() string {
	return p.Title
}

func (p *AddAssetInAppProposal) GetDescription() string {
	return p.Description
}

func (p *AddAssetInAppProposal) ProposalRoute() string { return RouterKey }

func (p *AddAssetInAppProposal) ProposalType() string { return ProposalAddAssetInApp }

func (p *AddAssetInAppProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}
