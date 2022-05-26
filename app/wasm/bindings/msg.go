package bindings

type ComdexMessages struct {
	WhiteListAssetLocker *WhiteListAssetLocker `json:"white-list-asset-locker,omitempty"`
}

type WhiteListAssetLocker struct {
	AppMapId uint64 `json:"app_map_id"`
	AssetId  uint64 `json:"asset_id"`
}
