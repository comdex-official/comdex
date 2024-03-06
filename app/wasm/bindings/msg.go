package bindings

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/math"
)

type ComdexMessages struct {
	MsgWhiteListAssetLocker              *MsgWhiteListAssetLocker              `json:"msg_white_list_asset_locker,omitempty"`
	MsgWhitelistAppIDVaultInterest       *MsgWhitelistAppIDVaultInterest       `json:"msg_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIDLockerRewards       *MsgWhitelistAppIDLockerRewards       `json:"msg_whitelist_app_id_locker_rewards,omitempty"`
	MsgAddExtendedPairsVault             *MsgAddExtendedPairsVault             `json:"msg_add_extended_pairs_vault,omitempty"`
	MsgSetCollectorLookupTable           *MsgSetCollectorLookupTable           `json:"msg_set_collector_lookup_table,omitempty"`
	MsgSetAuctionMappingForApp           *MsgSetAuctionMappingForApp           `json:"msg_set_auction_mapping_for_app,omitempty"`
	MsgUpdatePairsVault                  *MsgUpdatePairsVault                  `json:"msg_update_pairs_vault,omitempty"`
	MsgUpdateCollectorLookupTable        *MsgUpdateCollectorLookupTable        `json:"msg_update_collector_lookup_table,omitempty"`
	MsgRemoveWhitelistAssetLocker        *MsgRemoveWhitelistAssetLocker        `json:"msg_remove_whitelist_asset_locker,omitempty"`
	MsgRemoveWhitelistAppIDVaultInterest *MsgRemoveWhitelistAppIDVaultInterest `json:"msg_remove_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIDLiquidation         *MsgWhitelistAppIDLiquidation         `json:"msg_whitelist_app_id_liquidation,omitempty"`
	MsgRemoveWhitelistAppIDLiquidation   *MsgRemoveWhitelistAppIDLiquidation   `json:"msg_remove_whitelist_app_id_liquidation,omitempty"`
	MsgAddAuctionParams                  *MsgAddAuctionParams                  `json:"msg_add_auction_params,omitempty"`
	MsgBurnGovTokensForApp               *MsgBurnGovTokensForApp               `json:"msg_burn_gov_tokens_for_app,omitempty"`
	MsgAddESMTriggerParams               *MsgAddESMTriggerParams               `json:"msg_add_e_s_m_trigger_params,omitempty"`
	MsgEmissionRewards                   *MsgEmissionRewards                   `json:"msg_emission_rewards,omitempty"`
	MsgFoundationEmission                *MsgFoundationEmission                `json:"msg_foundation_emission,omitempty"`
	MsgRebaseMint                        *MsgRebaseMint                        `json:"msg_rebase_mint,omitempty"`
	MsgGetSurplusFund                    *MsgGetSurplusFund                    `json:"msg_get_surplus_fund,omitempty"`
	MsgEmissionPoolRewards               *MsgEmissionPoolRewards               `json:"msg_emission_pool_rewards,omitempty"`
}

type MsgWhiteListAssetLocker struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgWhitelistAppIDVaultInterest struct {
	AppID uint64 `json:"app_id"`
}

type MsgWhitelistAppIDLockerRewards struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgAddExtendedPairsVault struct {
	AppID               uint64  `json:"app_id"`
	PairID              uint64  `json:"pair_id"`
	StabilityFee        sdk.Dec `json:"stability_fee"`
	ClosingFee          sdk.Dec `json:"closing_fee"`
	LiquidationPenalty  sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee         sdk.Dec `json:"draw_down_fee"`
	IsVaultActive       bool    `json:"is_vault_active"`
	DebtCeiling         sdk.Int `json:"debt_ceiling"`
	DebtFloor           sdk.Int `json:"debt_floor"`
	IsStableMintVault   bool    `json:"is_stable_mint_vault"`
	MinCr               sdk.Dec `json:"min_cr"`
	PairName            string  `json:"pair_name"`
	AssetOutOraclePrice bool    `json:"asset_out_oracle_price"`
	AssetOutPrice       uint64  `json:"asset_out_price"`
	MinUsdValueLeft     uint64  `json:"min_usd_value_left"`
}

type MsgSetCollectorLookupTable struct {
	AppID            uint64  `json:"app_id"`
	CollectorAssetID uint64  `json:"collector_asset_id"`
	SecondaryAssetID uint64  `json:"secondary_asset_id"`
	SurplusThreshold sdk.Int `json:"surplus_threshold"`
	DebtThreshold    sdk.Int `json:"debt_threshold"`
	LockerSavingRate sdk.Dec `json:"locker_saving_rate"`
	LotSize          sdk.Int `json:"lot_size"`
	BidFactor        sdk.Dec `json:"bid_factor"`
	DebtLotSize      sdk.Int `json:"debt_lot_size"`
}

type MsgSetAuctionMappingForApp struct {
	AppID                uint64 `json:"app_id"`
	AssetIDs             uint64 `json:"asset_id"`
	IsSurplusAuctions    bool   `json:"is_surplus_auction"`
	IsDebtAuctions       bool   `json:"is_debt_auction"`
	IsDistributor        bool   `json:"is_distributor"`
	AssetOutOraclePrices bool   `json:"asset_out_oracle_price"`
	AssetOutPrices       uint64 `json:"asset_out_price"`
}

type MsgUpdatePairsVault struct {
	AppID              uint64  `json:"app_id"`
	ExtPairID          uint64  `json:"ext_pair_id"`
	StabilityFee       sdk.Dec `json:"stability_fee"`
	ClosingFee         sdk.Dec `json:"closing_fee"`
	LiquidationPenalty sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee        sdk.Dec `json:"draw_down_fee"`
	IsVaultActive      bool    `json:"is_vault_active"`
	MinCr              sdk.Dec `json:"min_cr"`
	DebtCeiling        sdk.Int `json:"debt_ceiling"`
	DebtFloor          sdk.Int `json:"debt_floor"`
	MinUsdValueLeft    uint64  `json:"min_usd_value_left"`
}

type MsgUpdateCollectorLookupTable struct {
	AppID            uint64  `json:"app_id"`
	AssetID          uint64  `json:"asset_id"`
	DebtThreshold    sdk.Int `json:"debt_threshold"`
	SurplusThreshold sdk.Int `json:"surplus_threshold"`
	LotSize          sdk.Int `json:"lot_size"`
	DebtLotSize      sdk.Int `json:"debt_lot_size"`
	BidFactor        sdk.Dec `json:"bid_factor"`
	LSR              sdk.Dec `json:"lsr"`
}

type MsgRemoveWhitelistAssetLocker struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgRemoveWhitelistAppIDVaultInterest struct {
	AppMappingID uint64 `json:"app_id"`
}

type MsgWhitelistAppIDLiquidation struct {
	AppID uint64 `json:"app_id"`
}

type MsgRemoveWhitelistAppIDLiquidation struct {
	AppID uint64 `json:"app_id"`
}

type MsgAddAuctionParams struct {
	AppID                  uint64  `json:"app_id"`
	AuctionDurationSeconds uint64  `json:"auction_duration_seconds"`
	Buffer                 sdk.Dec `json:"buffer"`
	Cusp                   sdk.Dec `json:"cusp"`
	Step                   uint64  `json:"step"`
	PriceFunctionType      uint64  `json:"price_function_type"`
	SurplusID              uint64  `json:"surplus_id"`
	DebtID                 uint64  `json:"debt_id"`
	DutchID                uint64  `json:"dutch_id"`
	BidDurationSeconds     uint64  `json:"bid_duration_seconds"`
}

type MsgBurnGovTokensForApp struct {
	AppID  uint64         `json:"app_id"`
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coin       `json:"amount"`
}

type MsgAddESMTriggerParams struct {
	AppID         uint64   `json:"app_id"`
	TargetValue   sdk.Coin `json:"target_value"`
	CoolOffPeriod uint64   `json:"cool_off_period"`
	AssetID       []uint64 `json:"asset_id"`
	Rates         []uint64 `json:"rates"`
}

type MsgEmissionRewards struct {
	AppID          uint64    `json:"app_id"`
	Amount         sdk.Int   `json:"amount"`
	EmissionAmount uint64    `json:"emission_amount"`
	ExtendedPair   []uint64  `json:"extended_pair"`
	VotingRatio    []sdk.Int `json:"voting_ratio"`
}

type MsgFoundationEmission struct {
	AppID             uint64   `json:"app_id"`
	Amount            sdk.Int  `json:"amount"`
	FoundationAddress []string `json:"foundation_address"`
}

type MsgRebaseMint struct {
	AppID        uint64         `json:"app_id"`
	Amount       sdk.Int        `json:"amount"`
	ContractAddr sdk.AccAddress `json:"contract_addr"`
}

type MsgGetSurplusFund struct {
	AppID        uint64         `json:"app_id"`
	AssetID      uint64         `json:"asset_id"`
	ContractAddr sdk.AccAddress `json:"contract_addr"`
	Amount       sdk.Coin       `json:"amount"`
}

type MsgEmissionPoolRewards struct {
	AppID       uint64    `json:"app_id"`
	CswapAppID  uint64    `json:"cswap_app_id"`
	Amount      sdk.Int   `json:"amount"`
	Pools       []uint64  `json:"pools"`
	VotingRatio []sdk.Int `json:"voting_ratio"`
}

type TokenFactoryMsg struct {
	/// Contracts can create denoms, namespaced under the contract's address.
	/// A contract may create any number of independent sub-denoms.
	CreateDenom *CreateDenom `json:"create_denom,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeAdmin *ChangeAdmin `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	MintTokens *MintTokens `json:"mint_tokens,omitempty"`
	/// Contracts can burn native tokens for an existing factory denom
	/// that they are the admin of.
	/// Currently, the burn from address must be the admin contract.
	BurnTokens *BurnTokens `json:"burn_tokens,omitempty"`
	/// Sets the metadata on a denom which the contract controls
	SetMetadata *SetMetadata `json:"set_metadata,omitempty"`
	/// Forces a transfer of tokens from one address to another.
	ForceTransfer *ForceTransfer `json:"force_transfer,omitempty"`
}

// CreateDenom creates a new factory denom, of denomination:
// factory/{creating contract address}/{Subdenom}
// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
// The (creating contract address, subdenom) pair must be unique.
// The created denom's admin is the creating contract address,
// but this admin can be changed using the ChangeAdmin binding.
type CreateDenom struct {
	Subdenom string    `json:"subdenom"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// ChangeAdmin changes the admin for a factory denom.
// If the NewAdminAddress is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom           string `json:"denom"`
	NewAdminAddress string `json:"new_admin_address"`
}

type MintTokens struct {
	Denom         string   `json:"denom"`
	Amount        math.Int `json:"amount"`
	MintToAddress string   `json:"mint_to_address"`
}

type BurnTokens struct {
	Denom           string   `json:"denom"`
	Amount          math.Int `json:"amount"`
	BurnFromAddress string   `json:"burn_from_address"`
}

type SetMetadata struct {
	Denom    string   `json:"denom"`
	Metadata Metadata `json:"metadata"`
}

type ForceTransfer struct {
	Denom       string   `json:"denom"`
	Amount      math.Int `json:"amount"`
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
}
