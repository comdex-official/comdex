package bindings

import sdk "github.com/cosmos/cosmos-sdk/types"

type ComdexQuery struct {
	AppData                                *AppData                                `json:"get_app,omitempty"`
	AssetData                              *AssetData                              `json:"get_asset_data,omitempty"`
	State                                  *State                                  `json:"state,omitempty"`
	MintedToken                            *MintedToken                            `json:"total_supply,omitempty"`
	RemoveWhiteListAssetLocker             *RemoveWhiteListAssetLocker             `json:"remove_white_list_asset,omitempty"`
	WhitelistAppIDVaultInterest            *WhitelistAppIDVaultInterest            `json:"whitelist_app_id_vault_interest,omitempty"`
	WhitelistAppIDLockerRewards            *WhitelistAppIDLockerRewards            `json:"whitelist_app_id_locker_rewards,omitempty"`
	ExternalLockerRewards                  *ExternalLockerRewards                  `json:"external_locker_rewards,omitempty"`
	ExternalVaultRewards                   *ExternalVaultRewards                   `json:"external_vault_rewards,omitempty"`
	CollectorLookupTableQuery              *CollectorLookupTableQuery              `json:"collector_lookup_table_query,omitempty"`
	ExtendedPairsVaultRecordsQuery         *ExtendedPairsVaultRecordsQuery         `json:"extended_pairs_vault_records_query,omitempty"`
	AuctionMappingForAppQuery              *AuctionMappingForAppQuery              `json:"auction_mapping_for_app_query,omitempty"`
	WhiteListedAssetQuery                  *WhiteListedAssetQuery                  `json:"white_listed_asset_query,omitempty"`
	UpdatePairsVaultQuery                  *UpdatePairsVaultQuery                  `json:"update_pairs_vault_query,omitempty"`
	UpdateCollectorLookupTableQuery        *UpdateCollectorLookupTableQuery        `json:"update_collector_lookup_table_query,omitempty"`
	RemoveWhitelistAssetLockerQuery        *RemoveWhitelistAssetLockerQuery        `json:"remove_whitelist_asset_locker_query,omitempty"`
	RemoveWhitelistAppIDVaultInterestQuery *RemoveWhitelistAppIDVaultInterestQuery `json:"remove_whitelist_app_id_vault_interest_query,omitempty"`
	WhitelistAppIDLiquidationQuery         *WhitelistAppIDLiquidationQuery         `json:"whitelist_app_id_liquidation_query,omitempty"`
	RemoveWhitelistAppIDLiquidationQuery   *RemoveWhitelistAppIDLiquidationQuery   `json:"remove_whitelist_app_id_liquidation_query,omitempty"`
	AddESMTriggerParamsForAppQuery         *AddESMTriggerParamsForAppQuery         `json:"add_e_s_m_trigger_params_for_app_query,omitempty"`
	ExtendedPairByApp                      *ExtendedPairByApp                      `json:"extended_pair_by_app,omitempty"`
	CheckSurplusReward                     *CheckSurplusReward                     `json:"check_surplus_reward,omitempty"`
	CheckWhitelistedAsset                  *CheckWhitelistedAsset                  `json:"check_whitelisted_asset,omitempty"`
	CheckVaultCreated                      *CheckVaultCreated                      `json:"check_vault_created,omitempty"`
	CheckBorrowed                          *CheckBorrowed                          `json:"check_borrowed,omitempty"`
	CheckLiquidityProvided                 *CheckLiquidityProvided                 `json:"check_liquidity_provided,omitempty"`
}

type AppData struct {
	AppID uint64 `json:"app_id"`
}

type AppDataResponse struct {
	MinGovDeposit    string `json:"min_gov_deposit"`
	GovTimeInSeconds int64  `json:"gov_time_in_seconds"`
	AssetID          uint64 `json:"gov_token_id"` // only when isGovToken true
}

type AssetData struct {
	AssetID uint64 `json:"asset_id"`
}

type AssetDataResponse struct {
	Denom string `json:"denom"`
}

type State struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
	Height  string `json:"height"`
	Target  string `json:"target"`
}

type StateResponse struct {
	Amount sdk.Coin `json:"amount"`
}

type MintedToken struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MintedTokenResponse struct {
	MintedTokens int64 `json:"current_supply"`
}

type RemoveWhiteListAssetLocker struct {
	AppID    uint64 `json:"app_id"`
	AssetIDs uint64 `json:"asset_ids"`
}

type RemoveWhiteListAssetResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type WhitelistAppIDVaultInterest struct {
	AppID uint64 `json:"app_id"`
}

type WhitelistAppIDVaultInterestResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type WhitelistAppIDLockerRewards struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type WhitelistAppIDLockerRewardsResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type ExternalLockerRewards struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type ExternalLockerRewardsResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type ExternalVaultRewards struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type ExternalVaultRewardsResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type CollectorLookupTableQuery struct {
	AppID            uint64 `json:"app_id"`
	CollectorAssetID uint64 `json:"collector_asset_id"`
	SecondaryAssetID uint64 `json:"secondary_asset_id"`
}

type CollectorLookupTableQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type ExtendedPairsVaultRecordsQuery struct {
	AppID        uint64  `json:"app_id"`
	PairID       uint64  `json:"pair_id"`
	StabilityFee sdk.Dec `json:"stability_fee"`
	ClosingFee   sdk.Dec `json:"closing_fee"`
	DrawDownFee  sdk.Dec `json:"draw_down_fee"`
	DebtCeiling  sdk.Int `json:"debt_ceiling"`
	DebtFloor    sdk.Int `json:"debt_floor"`
	PairName     string  `json:"pair_name"`
}

type ExtendedPairsVaultRecordsQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type AuctionMappingForAppQuery struct {
	AppID uint64 `json:"app_id"`
}

type AuctionMappingForAppQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type WhiteListedAssetQuery struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type WhiteListedAssetQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type UpdatePairsVaultQuery struct {
	AppID     uint64 `json:"app_id"`
	ExtPairID uint64 `json:"ext_pair_id"`
}

type UpdatePairsVaultQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type UpdateCollectorLookupTableQuery struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type UpdateCollectorLookupTableQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type RemoveWhitelistAssetLockerQuery struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type RemoveWhitelistAssetLockerQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type RemoveWhitelistAppIDVaultInterestQuery struct {
	AppID uint64 `json:"app_id"`
}

type RemoveWhitelistAppIDVaultInterestQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type WhitelistAppIDLiquidationQuery struct {
	AppID uint64 `json:"app_id"`
}

type WhitelistAppIDLiquidationQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type RemoveWhitelistAppIDLiquidationQuery struct {
	AppID uint64 `json:"app_id"`
}

type RemoveWhitelistAppIDLiquidationQueryResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type AddESMTriggerParamsForAppQuery struct {
	AppID uint64 `json:"app_id"`
}

type AddESMTriggerParamsForAppResponse struct {
	Found bool   `json:"found"`
	Err   string `json:"err"`
}

type ExtendedPairByApp struct {
	AppID uint64 `json:"app_id"`
}

type ExtendedPairByAppResponse struct {
	ExtendedPair []uint64 `json:"ext_pair"`
}

type CheckSurplusReward struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type CheckSurplusRewardResponse struct {
	Amount sdk.Coin `json:"amount"`
}

type CheckWhitelistedAsset struct {
	Denom string `json:"denom"`
}

type CheckWhitelistedAssetResponse struct {
	Found bool `json:"found"`
}

type CheckVaultCreated struct {
	Address string `json:"address"`
	AppID   uint64 `json:"app_id"`
}

type VaultCreatedResponse struct {
	IsCompleted bool `json:"is_completed"`
}

type CheckBorrowed struct {
	AssetID uint64 `json:"asset_id"`
	Address string `json:"address"`
}

type BorrowedResponse struct {
	IsCompleted bool `json:"is_completed"`
}

type CheckLiquidityProvided struct {
	AppID   uint64 `json:"app_id"`
	PoolID  uint64 `json:"pool_id"`
	Address string `json:"address"`
}

type LiquidityProvidedResponse struct {
	IsCompleted bool `json:"is_completed"`
}
