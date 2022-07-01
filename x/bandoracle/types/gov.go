package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalFetchPrice = "FetchPrice"
)

func init() {
	govtypes.RegisterProposalType(ProposalFetchPrice)
	govtypes.RegisterProposalTypeCodec(&FetchPriceProposal{}, "comdex/FetchPriceProposal")
}

var (
	_ govtypes.Content = &FetchPriceProposal{}
)

func NewFetchPriceProposal(title, description string, fetchPrice MsgFetchPriceData) govtypes.Content {
	return &FetchPriceProposal{
		Title:       title,
		Description: description,
		FetchPrice:  fetchPrice,
	}
}

func (p *FetchPriceProposal) ProposalRoute() string { return RouterKey }

func (p *FetchPriceProposal) ProposalType() string { return ProposalFetchPrice }

func (p *FetchPriceProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	return nil
}
