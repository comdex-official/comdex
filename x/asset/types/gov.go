package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalAddAsset = "AddAssets"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddAsset)
	govtypes.RegisterProposalTypeCodec(&AddAssetsProposal{}, "comdex/AddAssetsProposal")
}

var (
	_ govtypes.Content = &AddAssetsProposal{}
)

func NewUpdateLiquidationRatioProposal(title, description string, assets []Asset) govtypes.Content {
	return &AddAssetsProposal{
		Title:       title,
		Description: description,
		Assets:      assets,
	}
}

func (p *AddAssetsProposal) GetTitle() string { return p.Title }

func (p *AddAssetsProposal) GetDescription() string { return p.Description }

func (p *AddAssetsProposal) ProposalRoute() string { return RouterKey }

func (p *AddAssetsProposal) ProposalType() string { return ProposalAddAsset }

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
