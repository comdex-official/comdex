package bindings

type ComdexMessages struct {
	MsgWhiteListAssetLocker        MsgWhiteListAssetLocker        `json:"msg_white_list_asset_locker,omitempty"`
	MsgWhitelistAppIdVaultInterest MsgWhitelistAppIdVaultInterest `json:"msg_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIdLockerRewards MsgWhitelistAppIdLockerRewards `json:"msg_whitelist_app_id_locker_rewards,omitempty"`
}

type MsgWhiteListAssetLocker struct {
	AppMappingId uint64 `json:"app_mapping_id"`
	AssetId      uint64 `json:"asset_id"`
}

type MsgWhitelistAppIdVaultInterest struct {
	AppMappingId uint64 `json:"app_mapping_id"`
}

type MsgWhitelistAppIdLockerRewards struct {
	AppMappingId uint64   `json:"app_mapping_id"`
	AssetId      []uint64 `json:"asset_id"`
}
