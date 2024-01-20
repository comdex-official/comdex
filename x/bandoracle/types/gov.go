package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalFetchPrice = "FetchPrice"
)

func init() {
	govtypes.RegisterProposalType(ProposalFetchPrice)
}

var _ govtypes.Content = &FetchPriceProposal{}

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
	if p.FetchPrice.TwaBatchSize == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid batch size")
	}
	if err != nil {
		return err
	}
	return nil
}
