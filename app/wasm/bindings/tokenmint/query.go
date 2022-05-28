package tokenmint

import (
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
)

type TokenMintQuery struct {
	MintedToken *MintedToken `json:"total_supply,omitempty"`
}

type MintedToken struct {
	App_Id   uint64 `json:"app_id"`
	Asset_Id uint64 `json:"asset_id"`
}

type MintedTokenResponse struct {
	MintedTokens tokenminttypes.MintedTokens `json:"total_supply_response"`
}
