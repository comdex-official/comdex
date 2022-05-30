package bindings

type ComdexMessages struct {
	WhiteListAssetLocker WhiteListAssetLocker `json:"white_list_asset_locker,omitempty"`
}

type WhiteListAssetLocker struct {
	AppMappingId uint64 `json:"app_mapping_id"`
	AssetId      uint64 `json:"asset_id"`
}
