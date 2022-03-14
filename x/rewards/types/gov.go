package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeNewMintRewards = "NewMintRewardsProposal"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeNewMintRewards)
	govtypes.RegisterProposalTypeCodec(&NewMintRewardsProposal{}, "comdex/NewMintRewardsProposal")
}

var _ govtypes.Content = &NewMintRewardsProposal{}

func AddNewMintRewardsProposalContent(
	title,
	description string,
	collateralDenom string,
	cassetsDenoms []string,
	total_rewards sdk.Coin,
	casset_maxcap uint64,
	duration_days uint64,
) govtypes.Content {
	return &NewMintRewardsProposal{
		Title:       title,
		Description: description,
		MintRewards: &MintingRewards{
			Id:                0,
			AllowedCollateral: collateralDenom,
			AllowedCassets:    cassetsDenoms,
			TotalRewards:      total_rewards,
			CassetMaxCap:      casset_maxcap,
			DurationDays:      duration_days,
			IsActive:          false,
		},
	}
}

func (p *NewMintRewardsProposal) GetTitle() string { return p.Title }

func (p *NewMintRewardsProposal) GetDescription() string { return p.Description }

func (p *NewMintRewardsProposal) ProposalRoute() string { return RouterKey }

func (p *NewMintRewardsProposal) ProposalType() string { return ProposalTypeNewMintRewards }

func (p *NewMintRewardsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	return nil
}
