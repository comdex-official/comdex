package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalAddAssets              = "AddAssets"
	ProposalUpdateAsset            = "UpdateAsset"
	ProposalAddPairs               = "AddPairs"
	ProposalAddWhitelistedAssets   = "AddWhitelistedAssets"
	ProposalUpdateWhitelistedAsset = "UpdateWhitelistedAsset"
	ProposalAddWhitelistedPairs    = "AddWhitelistedPairs"
	ProposalUpdateWhitelistedPair  = "UpdateWhitelistedPair"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddAssets)
	govtypes.RegisterProposalTypeCodec(&AddAssetsProposal{}, "comdex/AddAssetsProposal")

	govtypes.RegisterProposalType(ProposalUpdateAsset)
	govtypes.RegisterProposalTypeCodec(&UpdateAssetProposal{}, "comdex/UpdateAssetProposal")

	govtypes.RegisterProposalType(ProposalAddPairs)
	govtypes.RegisterProposalTypeCodec(&AddPairsProposal{}, "comdex/AddPairsProposal")

	govtypes.RegisterProposalType(ProposalAddWhitelistedAssets)
	govtypes.RegisterProposalTypeCodec(&AddWhitelistedAssetsProposal{}, "comdex/AddWhitelistedAssetsProposal")

	govtypes.RegisterProposalType(ProposalUpdateWhitelistedAsset)
	govtypes.RegisterProposalTypeCodec(&UpdateWhitelistedAssetProposal{}, "comdex/UpdateWhitelistedAssetProposal")

	govtypes.RegisterProposalType(ProposalAddWhitelistedPairs)
	govtypes.RegisterProposalTypeCodec(&AddWhitelistedPairsProposal{}, "comdex/AddWhitelistedPairsProposal")

	govtypes.RegisterProposalType(ProposalUpdateWhitelistedPair)
	govtypes.RegisterProposalTypeCodec(&UpdateWhitelistedPairProposal{}, "comdex/UpdateWhitelistedPairProposal")

}

var (
	_ govtypes.Content = &AddAssetsProposal{}
	_ govtypes.Content = &UpdateAssetProposal{}
	_ govtypes.Content = &AddPairsProposal{}
	_ govtypes.Content = &AddWhitelistedAssetsProposal{}
	_ govtypes.Content = &UpdateWhitelistedAssetProposal{}
	_ govtypes.Content = &AddWhitelistedPairsProposal{}
	_ govtypes.Content = &UpdateWhitelistedPairProposal{}
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

func NewAddWhitelistedAssetsProposal(title, description string, assets []ExtendedAsset) govtypes.Content {
	return &AddWhitelistedAssetsProposal{
		Title:       title,
		Description: description,
		Assets:      assets,
	}
}

func (p *AddWhitelistedAssetsProposal) GetTitle() string {
	return p.Title
}

func (p *AddWhitelistedAssetsProposal) GetDescription() string {
	return p.Description
}

func (p *AddWhitelistedAssetsProposal) ProposalRoute() string { return RouterKey }

func (p *AddWhitelistedAssetsProposal) ProposalType() string { return ProposalAddWhitelistedAssets }


func (p *AddWhitelistedAssetsProposal) ValidateBasic() error {

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

func NewUpdateWhitelistedAssetProposal(title, description string, asset ExtendedAsset) govtypes.Content {
	return &UpdateWhitelistedAssetProposal{
		Title:       title,
		Description: description,
		Asset:       asset,
	}
}

func (p *UpdateWhitelistedAssetProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateWhitelistedAssetProposal) GetDescription() string {
	return p.Description
}

func (p *UpdateWhitelistedAssetProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateWhitelistedAssetProposal) ProposalType() string { return ProposalUpdateWhitelistedAsset }

func (p *UpdateWhitelistedAssetProposal) ValidateBasic() error {

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

func NewAddWhitelistedPairsProposal(title, description string, pairs []ExtendedPairLend) govtypes.Content {
	return &AddWhitelistedPairsProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *AddWhitelistedPairsProposal) GetTitle() string {
	return p.Title
}

func (p *AddWhitelistedPairsProposal) GetDescription() string {
	return p.Description
}

func (p *AddWhitelistedPairsProposal) ProposalRoute() string { return RouterKey }

func (p *AddWhitelistedPairsProposal) ProposalType() string { return ProposalAddWhitelistedPairs }


func (p *AddWhitelistedPairsProposal) ValidateBasic() error {

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

func NewUpdateWhitelistedPairProposal(title, description string, pair ExtendedPairLend) govtypes.Content {
	return &UpdateWhitelistedPairProposal{
		Title:       title,
		Description: description,
		Pair:        pair,
	}
}

func (p *UpdateWhitelistedPairProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateWhitelistedPairProposal) GetDescription() string {
	return p.Description
}

func (p *UpdateWhitelistedPairProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateWhitelistedPairProposal) ProposalType() string { return ProposalUpdateWhitelistedPair }

func (p *UpdateWhitelistedPairProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pair := p.Pair
	if err := pair.Validate(); err != nil {
		return err
	}

	return nil
}
