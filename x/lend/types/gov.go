package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	ProposalAddLendPairs           = "ProposalAddLendPairs"
	ProposalAddMultipleLendPairs   = "ProposalAddMultipleLendPairs"
	ProposalAddPool                = "ProposalAddPool"
	ProposalAddAssetToPair         = "ProposalAddAssetToPair"
	ProposalAddAssetRatesParams    = "ProposalAddAssetRatesParams"
	ProposalAddAuctionParams       = "ProposalAddAuctionParams"
	ProposalAddMultipleAssetToPair = "ProposalAddMultipleAssetToPair"
	ProposalAddPoolPairs           = "ProposalAddPoolPairs"
	ProposalAddAssetRatesPoolPairs = "ProposalAddAssetRatesPoolPairs"
	ProposalDepreciatePool         = "ProposalDepreciatePool"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddLendPairs)
	govtypes.RegisterProposalTypeCodec(&LendPairsProposal{}, "comdex/AddLendPairsProposal")
	govtypes.RegisterProposalType(ProposalAddMultipleLendPairs)
	govtypes.RegisterProposalTypeCodec(&MultipleLendPairsProposal{}, "comdex/AddMultipleLendPairsProposal")
	govtypes.RegisterProposalType(ProposalAddPool)
	govtypes.RegisterProposalTypeCodec(&AddPoolsProposal{}, "comdex/AddPoolsProposal")
	govtypes.RegisterProposalType(ProposalAddAssetToPair)
	govtypes.RegisterProposalTypeCodec(&AddAssetToPairProposal{}, "comdex/AddAssetToPairProposal")
	govtypes.RegisterProposalType(ProposalAddMultipleAssetToPair)
	govtypes.RegisterProposalTypeCodec(&AddMultipleAssetToPairProposal{}, "comdex/AddMultipleAssetToPairProposal")
	govtypes.RegisterProposalType(ProposalAddAssetRatesParams)
	govtypes.RegisterProposalTypeCodec(&AddAssetRatesParams{}, "comdex/AddAssetRatesParams")
	govtypes.RegisterProposalType(ProposalAddAuctionParams)
	govtypes.RegisterProposalTypeCodec(&AddAuctionParamsProposal{}, "comdex/AddAuctionParamsProposal")
	govtypes.RegisterProposalType(ProposalAddPoolPairs)
	govtypes.RegisterProposalTypeCodec(&AddPoolPairsProposal{}, "comdex/AddPoolPairsProposal")
	govtypes.RegisterProposalType(ProposalAddAssetRatesPoolPairs)
	govtypes.RegisterProposalTypeCodec(&AddAssetRatesPoolPairsProposal{}, "comdex/AddAssetRatesPoolPairsProposal")
	govtypes.RegisterProposalType(ProposalDepreciatePool)
	govtypes.RegisterProposalTypeCodec(&AddPoolDepreciateProposal{}, "comdex/AddPoolDepreciateProposal")
}

var (
	_ govtypes.Content = &LendPairsProposal{}
	_ govtypes.Content = &AddPoolsProposal{}
	_ govtypes.Content = &AddAssetToPairProposal{}
	_ govtypes.Content = &AddAssetRatesParams{}
	_ govtypes.Content = &AddAuctionParamsProposal{}
	_ govtypes.Content = &MultipleLendPairsProposal{}
	_ govtypes.Content = &AddMultipleAssetToPairProposal{}
	_ govtypes.Content = &AddPoolPairsProposal{}
	_ govtypes.Content = &AddAssetRatesPoolPairsProposal{}
	_ govtypes.Content = &AddPoolDepreciateProposal{}
)

func NewAddLendPairsProposal(title, description string, pairs Extended_Pair) govtypes.Content {
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

	if err := p.Pairs.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddMultipleLendPairsProposal(title, description string, pairs []Extended_Pair) govtypes.Content {
	return &MultipleLendPairsProposal{
		Title:       title,
		Description: description,
		Pairs:       pairs,
	}
}

func (p *MultipleLendPairsProposal) ProposalRoute() string { return RouterKey }

func (p *MultipleLendPairsProposal) ProposalType() string { return ProposalAddMultipleLendPairs }

func (p *MultipleLendPairsProposal) ValidateBasic() error {
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

func NewAddMultipleAssetToPairProposal(title, description string, AssetToPairMapping []AssetToPairSingleMapping) govtypes.Content {
	return &AddMultipleAssetToPairProposal{
		Title:                    title,
		Description:              description,
		AssetToPairSingleMapping: AssetToPairMapping,
	}
}

func (p *AddMultipleAssetToPairProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddMultipleAssetToPairProposal) ProposalType() string {
	return ProposalAddAssetToPair
}

func (p *AddMultipleAssetToPairProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if len(p.AssetToPairSingleMapping) == 0 {
		return ErrorEmptyProposalAssets
	}
	for _, pair := range p.AssetToPairSingleMapping {
		if err := pair.Validate(); err != nil {
			return err
		}
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

func NewAddassetRatesParams(title, description string, AssetRatesParams AssetRatesParams) govtypes.Content {
	return &AddAssetRatesParams{
		Title:            title,
		Description:      description,
		AssetRatesParams: AssetRatesParams,
	}
}

func (p *AddAssetRatesParams) ProposalRoute() string {
	return RouterKey
}

func (p *AddAssetRatesParams) ProposalType() string {
	return ProposalAddAssetRatesParams
}

func (p *AddAssetRatesParams) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if err = p.AssetRatesParams.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddAuctionParams(title, description string, AddAuctionParams AuctionParams) govtypes.Content {
	return &AddAuctionParamsProposal{
		Title:         title,
		Description:   description,
		AuctionParams: AddAuctionParams,
	}
}

func (p *AddAuctionParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddAuctionParamsProposal) ProposalType() string {
	return ProposalAddAuctionParams
}

func (p *AddAuctionParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewAddPoolPairsProposal(title, description string, pool PoolPairs) govtypes.Content {
	return &AddPoolPairsProposal{
		Title:       title,
		Description: description,
		PoolPairs:   pool,
	}
}

func (p *AddPoolPairsProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddPoolPairsProposal) ProposalType() string {
	return ProposalAddPool
}

func (p *AddPoolPairsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	pool := p.PoolPairs
	if err := pool.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddassetRatesPoolPairs(title, description string, AssetRatesPoolPairs AssetRatesPoolPairs) govtypes.Content {
	return &AddAssetRatesPoolPairsProposal{
		Title:               title,
		Description:         description,
		AssetRatesPoolPairs: AssetRatesPoolPairs,
	}
}

func (p *AddAssetRatesPoolPairsProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddAssetRatesPoolPairsProposal) ProposalType() string {
	return ProposalAddAssetRatesPoolPairs
}

func (p *AddAssetRatesPoolPairsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if err = p.AssetRatesPoolPairs.Validate(); err != nil {
		return err
	}

	return nil
}

func NewAddDepreciatePool(title, description string, PoolDepreciate PoolDepreciate) govtypes.Content {
	return &AddPoolDepreciateProposal{
		Title:          title,
		Description:    description,
		PoolDepreciate: PoolDepreciate,
	}
}

func (p *AddPoolDepreciateProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddPoolDepreciateProposal) ProposalType() string {
	return ProposalDepreciatePool
}

func (p *AddPoolDepreciateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if err = p.PoolDepreciate.Validate(); err != nil {
		return err
	}

	return nil
}
