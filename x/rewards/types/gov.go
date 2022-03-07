package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeNewMintRewards = "NewMintRewards"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeNewMintRewards)
	govtypes.RegisterProposalTypeCodec(&NewMintRewards{}, "comdex/NewMintRewardsProposal")
}

var _ govtypes.Content = &NewMintRewards{}

func NewMintRewardsProposal(title, description string, collateralDenom string, cassetsDenoms []string, total_rewards sdk.Coin, casset_maxcap uint64, duration_days uint64) govtypes.Content {
	return &NewMintRewards{
		Title:             title,
		Description:       description,
		AllowedCollateral: collateralDenom,
		AllowedCassets:    cassetsDenoms,
		TotalRewards:      total_rewards,
		CassetMaxCap:      casset_maxcap,
		DurationDays:      duration_days,
	}
}

func (p *NewMintRewards) GetTitle() string { return p.Title }

func (p *NewMintRewards) GetDescription() string { return p.Description }

func (p *NewMintRewards) ProposalRoute() string { return RouterKey }

func (p *NewMintRewards) ProposalType() string { return ProposalTypeNewMintRewards }

func (p *NewMintRewards) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	return nil
}
