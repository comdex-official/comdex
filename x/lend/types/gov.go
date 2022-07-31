package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	ProposalAddLendPairs       = "ProposalAddLendPairs"
	ProposalUpdateLendPair     = "ProposalUpdateLendPair"
	ProposalAddPool            = "ProposalAddPool"
	ProposalAddAssetToPair     = "ProposalAddAssetToPair"
	ProposalAddAssetRatesStats = "ProposalAddAssetRatesStats"
	ProposalAddAuctionParams = "ProposalAddAuctionParams"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddLendPairs)
	govtypes.RegisterProposalTypeCodec(&LendPairsProposal{}, "comdex/AddLendPairsProposal")
	govtypes.RegisterProposalType(ProposalUpdateLendPair)
	govtypes.RegisterProposalTypeCodec(&UpdatePairProposal{}, "comdex/UpdateLenddPairProposal")
	govtypes.RegisterProposalType(ProposalAddPool)
	govtypes.RegisterProposalTypeCodec(&AddPoolsProposal{}, "comdex/AddPoolsProposal")
	govtypes.RegisterProposalType(ProposalAddAssetToPair)
	govtypes.RegisterProposalTypeCodec(&AddAssetToPairProposal{}, "comdex/AddAssetToPairProposal")
	govtypes.RegisterProposalType(ProposalAddAssetRatesStats)
	govtypes.RegisterProposalTypeCodec(&AddAssetRatesStats{}, "comdex/AddAssetRatesStats")
	govtypes.RegisterProposalType(ProposalAddAuctionParams)
	govtypes.RegisterProposalTypeCodec(&AddAuctionParamsProposal{}, "comdex/AddAuctionParamsProposal")
}

var (
	_ govtypes.Content = &LendPairsProposal{}
	_ govtypes.Content = &UpdatePairProposal{}
	_ govtypes.Content = &AddPoolsProposal{}
	_ govtypes.Content = &AddAssetToPairProposal{}
	_ govtypes.Content = &AddAssetRatesStats{}
	_ govtypes.Content = &AddAuctionParamsProposal{}
)

func NewAddLendPairsProposal(title, description string, pairs []Extended_Pair) govtypes.Content {
	return &LendPairsProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *LendPairsProposal) ProposalRoute() string { return RouterKey }

func (p *LendPairsProposal) ProposalType() string { return ProposalAddLendPairs }

func (p *LendPairsProposal) ValidateBasic() error {
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

func NewUpdateLendPairProposal(title, description string, pair Extended_Pair) govtypes.Content {
	return &UpdatePairProposal{
		Title:       title,
		Description: description,
		Pair:        pair,
	}
}

func (p *UpdatePairProposal) ProposalRoute() string { return RouterKey }

func (p *UpdatePairProposal) ProposalType() string { return ProposalUpdateLendPair }

func (p *UpdatePairProposal) ValidateBasic() error {
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

func NewAddPoolProposal(title, description string, pool Pool) govtypes.Content {
	return &AddPoolsProposal{
		Title:       title,
		Description: description,
		Pool:        pool,
	}
}

func (p *AddPoolsProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddPoolsProposal) ProposalType() string {
	return ProposalAddPool
}

func (p *AddPoolsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pool := p.Pool
	if err := pool.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddAssetToPairProposal(title, description string, AssetToPairMapping AssetToPairMapping) govtypes.Content {
	return &AddAssetToPairProposal{
		Title:              title,
		Description:        description,
		AssetToPairMapping: AssetToPairMapping,
	}
}

func (p *AddAssetToPairProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddAssetToPairProposal) ProposalType() string {
	return ProposalAddAssetToPair
}

func (p *AddAssetToPairProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pool := p.AssetToPairMapping
	if err := pool.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddAssetRatesStats(title, description string, AssetRatesStats []AssetRatesStats) govtypes.Content {
	return &AddAssetRatesStats{
		Title:           title,
		Description:     description,
		AssetRatesStats: AssetRatesStats,
	}
}

func (p *AddAssetRatesStats) ProposalRoute() string {
	return RouterKey
}

func (p *AddAssetRatesStats) ProposalType() string {
	return ProposalAddAssetRatesStats
}

func (p *AddAssetRatesStats) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	for _, assetRatesStats := range p.AssetRatesStats {
		if err := assetRatesStats.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewAddAuctionParams(title, description string, AddAuctionParams AuctionParams) govtypes.Content {
	return &AddAuctionParamsProposal{
		Title:           title,
		Description:     description,
		AuctionParams: AddAuctionParams,
	}
}

func (p *AddAuctionParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddAuctionParamsProposal) ProposalType() string {
	return ProposalAddAssetRatesStats
}

func (p *AddAuctionParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}
