package types

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

const (
	ProposalAddWhitelistedAssets = "AddWhitelistedAssets"
	ProposalUpdateAsset          = "UpdateWhitelistedAsset"
	ProposalAddPairs             = "AddWhitelistedPairs"
	ProposalUpdatePair           = "UpdateWhitelistedPair"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddWhitelistedAssets)
	govtypes.RegisterProposalTypeCodec(&AddWhitelistedAssetsProposal{}, "comdex/AddWhitelistedAssetsProposal")

	govtypes.RegisterProposalType(ProposalUpdateAsset)
	govtypes.RegisterProposalTypeCodec(&UpdateWhitelistedAssetProposal{}, "comdex/UpdateWhitelistedAssetProposal")

	govtypes.RegisterProposalType(ProposalAddPairs)
	govtypes.RegisterProposalTypeCodec(&AddWhitelistedPairsProposal{}, "comdex/AddWhitelistedPairsProposal")

	govtypes.RegisterProposalType(ProposalUpdatePair)
	govtypes.RegisterProposalTypeCodec(&UpdateWhitelistedPairProposal{}, "comdex/UpdateWhitelistedPairProposal")
}

var (
	_ govtypes.Content = &AddWhitelistedAssetsProposal{}
	_ govtypes.Content = &UpdateWhitelistedAssetProposal{}
	_ govtypes.Content = &AddWhitelistedPairsProposal{}
	_ govtypes.Content = &UpdateWhitelistedAssetProposal{}
)

func NewAddWhitelistedAssetsProposal(title, description string, assets []Asset) govtypes.Content {
	return &AddWhitelistedAssetsProposal{
		Title:       title,
		Description: description,
		Assets:      assets,
	}
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

func NewUpdateWhitelistedAssetProposal(title, description string, asset Asset) govtypes.Content {
	return &UpdateWhitelistedAssetProposal{
		Title:       title,
		Description: description,
		Asset:       asset,
	}
}

func (p *UpdateWhitelistedAssetProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateWhitelistedAssetProposal) ProposalType() string { return ProposalUpdateAsset }

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

func NewAddWhitelistedPairsProposal(title, description string, pairs []Pair) govtypes.Content {
	return &AddWhitelistedPairsProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *AddWhitelistedPairsProposal) ProposalRoute() string { return RouterKey }

func (p *AddWhitelistedPairsProposal) ProposalType() string { return ProposalAddPairs }

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

func NewUpdateWhitelistedPairProposal(title, description string, pair Pair) govtypes.Content {
	return &UpdateWhitelistedPairProposal{
		Title:       title,
		Description: description,
		Pair:        pair,
	}
}

func (p *UpdateWhitelistedPairProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateWhitelistedPairProposal) ProposalType() string { return ProposalUpdatePair }

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
