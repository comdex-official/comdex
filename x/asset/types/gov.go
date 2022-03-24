package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalAddAssets = "AddAssets"
	ProposalUpdateAsset = "UpdateAsset"
	ProposalAddPairs = "AddPairs"
	ProposalUpdatePair = "UpdatePair"
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
}

var (
	_ govtypes.Content = &AddAssetsProposal{}
	_ govtypes.Content = &UpdateAssetProposal{}
	_ govtypes.Content = &AddPairsProposal{}
	_ govtypes.Content = &UpdateAssetProposal{}
)

func NewAddAssetsProposal(title, description string, assets []Asset) govtypes.Content {
	return &AddAssetsProposal{
		Title:       title,
		Description: description,
		Assets:      assets,
	}
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
		Asset:      asset,
	}
}

func (p *UpdateAssetProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateAssetProposal) ProposalType() string { return ProposalUpdateAsset }

func (p *UpdateAssetProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	 asset :=  p.Asset
		if err := asset.Validate(); err != nil {
			return err
	}

	return nil
}



func NewAddPairsProposal(title, description string, pairs []Pair) govtypes.Content {
	return &AddPairsProposal{
		Title:       title,
		Description: description,
		Pairs:      pairs,
	}
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

func NewUpdatePairProposal(title, description string, pair Pair) govtypes.Content {
	return &UpdatePairProposal{
		Title:       title,
		Description: description,
		Pair:      pair,
	}
}

func (p *UpdatePairProposal) ProposalRoute() string { return RouterKey }

func (p *UpdatePairProposal) ProposalType() string { return ProposalUpdatePair }

func (p *UpdatePairProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pair :=  p.Pair
	if err := pair.Validate(); err != nil {
		return err
	}

	return nil
}
