package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeNewMintRewards     = "NewMintRewardsProposal"
	ProposalTypeDisableMintRewards = "DisbaleMintRewardsProposal"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeNewMintRewards)
	govtypes.RegisterProposalTypeCodec(&NewMintRewardsProposal{}, "comdex/NewMintRewardsProposal")

	govtypes.RegisterProposalType(ProposalTypeDisableMintRewards)
	govtypes.RegisterProposalTypeCodec(&DisbaleMintRewardsProposal{}, "comdex/DisbaleMintRewardsProposal")
}

var _ govtypes.Content = &NewMintRewardsProposal{}
var _ govtypes.Content = &DisbaleMintRewardsProposal{}

func AddNewMintRewardsProposalContent(
	title,
	description string,
	collateralDenom string,
	cassetsDenom string,
	total_rewards sdk.Coin,
	casset_maxcap uint64,
	duration_days uint64,
	minimumLockupTimeInSeconds uint64,
) govtypes.Content {
	return &NewMintRewardsProposal{
		Title:       title,
		Description: description,
		MintRewards: &MintingRewards{
			Id:                   0,
			AllowedCollateral:    collateralDenom,
			AllowedCasset:        cassetsDenom,
			TotalRewards:         total_rewards,
			CassetMaxCap:         casset_maxcap,
			DurationDays:         duration_days,
			IsActive:             false,
			AvailableRewards:     sdk.NewCoin(total_rewards.Denom, sdk.NewInt(0)),
			Depositor:            nil,
			MinLockupTimeSeconds: minimumLockupTimeInSeconds,
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

func DisableMintRewardsProposalContent(
	title,
	description string,
	mintingRewardId uint64,
) govtypes.Content {
	return &DisbaleMintRewardsProposal{
		Title:        title,
		Description:  description,
		MintRewardId: mintingRewardId,
	}
}

func (p *DisbaleMintRewardsProposal) GetTitle() string { return p.Title }

func (p *DisbaleMintRewardsProposal) GetDescription() string { return p.Description }

func (p *DisbaleMintRewardsProposal) ProposalRoute() string { return RouterKey }

func (p *DisbaleMintRewardsProposal) ProposalType() string { return ProposalTypeNewMintRewards }

func (p *DisbaleMintRewardsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	return nil
}
