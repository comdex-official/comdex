package bindings

type ComdexMessages struct {
	WhiteListAssetLocker *WhiteListAssetLocker `json:"mint_tokens,omitempty"`
}

type WhiteListAssetLocker struct {
	AppMapId uint64 `json:"app_map_id"`
	AssetId  uint64 `json:"asset_id"`
}
