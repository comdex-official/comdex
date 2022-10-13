package types

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

const (
	ProposalAddLendRewards = "AddLendRewards"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddLendRewards)
	govtypes.RegisterProposalTypeCodec(&AddLendExternalRewardsProposal{}, "comdex/AddLendExternalRewardsProposal")
}

var (
	_ govtypes.Content = &AddLendExternalRewardsProposal{}
)

func NewAddExternalLendRewardsProposal(title, description string, lendExternalRewards LendExternalRewards) govtypes.Content {
	return &AddLendExternalRewardsProposal{
		Title:               title,
		Description:         description,
		LendExternalRewards: lendExternalRewards,
	}
}

func (p *AddLendExternalRewardsProposal) GetTitle() string {
	return p.Title
}

func (p *AddLendExternalRewardsProposal) GetDescription() string {
	return p.Description
}
func (p *AddLendExternalRewardsProposal) ProposalRoute() string { return RouterKey }

func (p *AddLendExternalRewardsProposal) ProposalType() string { return ProposalAddLendRewards }

func (p *AddLendExternalRewardsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}
