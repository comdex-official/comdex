package locker

type LockerMessages struct {
	WhiteListAssetLocker WhiteListAssetLocker `json:"white-list-asset-locker,omitempty"`
}

type WhiteListAssetLocker struct {
	AppMappingId uint64 `json:"app_mapping_id"`
	AssetId      uint64 `json:"asset_id"`
}
