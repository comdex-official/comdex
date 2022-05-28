package tokenmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenMintQuery struct {
	MintedToken *MintedToken `json:"total_supply,omitempty"`
}

type MintedToken struct {
	App_Id   uint64 `json:"app_id"`
	Asset_Id uint64 `json:"asset_id"`
}

type MintedTokenResponse struct {
	MintedTokens sdk.Int `json:"total_supply_response"`
}
