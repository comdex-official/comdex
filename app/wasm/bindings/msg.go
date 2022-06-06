package bindings

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ComdexMessages struct {
	MsgWhiteListAssetLocker              *MsgWhiteListAssetLocker              `json:"msg_white_list_asset_locker,omitempty"`
	MsgWhitelistAppIdVaultInterest       *MsgWhitelistAppIdVaultInterest       `json:"msg_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIdLockerRewards       *MsgWhitelistAppIdLockerRewards       `json:"msg_whitelist_app_id_locker_rewards,omitempty"`
	MsgAddExtendedPairsVault             *MsgAddExtendedPairsVault             `json:"msg_add_extended_pairs_vault,omitempty"`
	MsgSetCollectorLookupTable           *MsgSetCollectorLookupTable           `json:"msg_set_collector_lookup_table,omitempty"`
	MsgSetAuctionMappingForApp           *MsgSetAuctionMappingForApp           `json:"msg_set_auction_mapping_for_app,omitempty"`
	MsgUpdateLsrInPairsVault             *MsgUpdateLsrInPairsVault             `json:"msg_update_lsr_in_pairs_vault,omitempty"`
	MsgUpdateLsrInCollectorLookupTable   *MsgUpdateLsrInCollectorLookupTable   `json:"msg_update_lsr_in_collector_lookup_table,omitempty"`
	MsgRemoveWhitelistAssetLocker        *MsgRemoveWhitelistAssetLocker        `json:"msg_remove_whitelist_asset_locker,omitempty"`
	MsgRemoveWhitelistAppIdVaultInterest *MsgRemoveWhitelistAppIdVaultInterest `json:"msg_remove_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIdLiquidation         *MsgWhitelistAppIdLiquidation         `json:"msg_whitelist_app_id_liquidation,omitempty"`
	MsgRemoveWhitelistAppIdLiquidation   *MsgRemoveWhitelistAppIdLiquidation   `json:"msg_remove_whitelist_app_id_liquidation,omitempty"`
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

type MsgAddExtendedPairsVault struct {
	AppMappingId        uint64  `json:"app_mapping_id"`
	PairId              uint64  `json:"pair_id"`
	LiquidationRatio    sdk.Dec `json:"liquidation_ratio"`
	StabilityFee        sdk.Dec `json:"stability_fee"`
	ClosingFee          sdk.Dec `json:"closing_fee"`
	LiquidationPenalty  sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee         sdk.Dec `json:"draw_down_fee"`
	IsVaultActive       bool    `json:"is_vault_active"`
	DebtCeiling         uint64  `json:"debt_ceiling"`
	DebtFloor           uint64  `json:"debt_floor"`
	IsPsmPair           bool    `json:"is_psm_pair"`
	MinCr               sdk.Dec `json:"min_cr"`
	PairName            string  `json:"pair_name"`
	AssetOutOraclePrice bool    `json:"asset_out_oracle_price"`
	AssetOutPrice       uint64  `json:"asset_out_price"`
}

type MsgSetCollectorLookupTable struct {
	AppMappingId     uint64  `json:"app_mapping_id"`
	CollectorAssetId uint64  `json:"collector_asset_id"`
	SecondaryAssetId uint64  `json:"secondary_asset_id"`
	SurplusThreshold uint64  `json:"surplus_threshold"`
	DebtThreshold    uint64  `json:"debt_threshold"`
	LockerSavingRate sdk.Dec `json:"locker_saving_rate"`
	LotSize          uint64  `json:"lot_size"`
	BidFactor        sdk.Dec `json:"bid_factor"`
}

type MsgSetAuctionMappingForApp struct {
	AppMappingId        uint64   `json:"app_mapping_id"`
	AssetId             []uint64 `json:"asset_id"`
	IsSurplusAuction    []bool   `json:"is_surplus_auction"`
	IsDebtAuction       []bool   `json:"is_debt_auction"`
	AssetOutOraclePrice []bool   `json:"asset_out_oracle_price"`
	AssetOutPrice       []uint64 `json:"asset_out_price"`
}

type MsgUpdateLsrInPairsVault struct {
	AppMappingId       uint64  `json:"app_mapping_id"`
	ExtPairId          uint64  `json:"ext_pair_id"`
	LiquidationRatio   sdk.Dec `json:"liquidation_ratio"`
	StabilityFee       sdk.Dec `json:"stability_fee"`
	ClosingFee         sdk.Dec `json:"closing_fee"`
	LiquidationPenalty sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee        sdk.Dec `json:"draw_down_fee"`
	MinCr              sdk.Dec `json:"min_cr"`
	DebtCeiling        uint64  `json:"debt_ceiling"`
	DebtFloor          uint64  `json:"debt_floor"`
}

type MsgUpdateLsrInCollectorLookupTable struct {
	AppMappingId uint64  `json:"app_mapping_id"`
	AssetId      uint64  `json:"asset_id"`
	LSR          sdk.Dec `json:"lsr"`
}

type MsgRemoveWhitelistAssetLocker struct {
	AppMappingId uint64 `json:"app_mapping_id"`
	AssetId      uint64 `json:"asset_id"`
}

type MsgRemoveWhitelistAppIdVaultInterest struct {
	AppMappingId uint64 `json:"app_mapping_id"`
}

type MsgWhitelistAppIdLiquidation struct {
	AppMappingId uint64 `json:"app_mapping_id"`
}

type MsgRemoveWhitelistAppIdLiquidation struct {
	AppMappingId uint64 `json:"app_mapping_id"`
}
