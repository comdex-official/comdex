package types

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

func NewGenesisState(auction []Auction, auctionParams AuctionParams, auctionFeesCollectionFromLimitBidTx []AuctionFeesCollectionFromLimitBidTx, params Params, auctionId, userBiddingID uint64) *GenesisState {
	return &GenesisState{
		Auction:                             auction,
		AuctionParams:                       auctionParams,
		AuctionFeesCollectionFromLimitBidTx: auctionFeesCollectionFromLimitBidTx,
		Params:                              params,
		AuctionId:                           auctionId,
		UserBiddingID:                       userBiddingID,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	var auctionId, userBiddingID uint64
	return &GenesisState{
		Auction:                             []Auction{},
		AuctionParams:                       AuctionParams{},
		AuctionFeesCollectionFromLimitBidTx: []AuctionFeesCollectionFromLimitBidTx{},
		Params:                              DefaultParams(),
		AuctionId:                           auctionId,
		UserBiddingID:                       userBiddingID,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
