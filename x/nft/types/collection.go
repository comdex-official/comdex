package types

import (
	"github.com/comdex-official/comdex/x/nft/exported"
)

func NewCollection(denom Denom, nfts []exported.NFT) (c Collection) {
	c.Denom = denom
	for _, nft := range nfts {
		c = c.AddNFT(nft.(NFT))
	}
	return c
}

func (c Collection) AddNFT(nft NFT) Collection {
	c.NFTs = append(c.NFTs, nft)
	return c
}

func (c Collection) Supply() int {
	return len(c.NFTs)
}

func NewCollections(c ...Collection) []Collection {
	return append([]Collection{}, c...)
}
