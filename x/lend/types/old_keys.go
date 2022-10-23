package types

var (
	PoolKeyPrefix         = []byte{0x13}
	LendUserPrefix        = []byte{0x15}
	LendHistoryIDPrefix   = []byte{0x16}
	PoolIDPrefix          = []byte{0x17}
	LendPairIDKey         = []byte{0x18}
	LendPairKeyPrefix     = []byte{0x19}
	BorrowHistoryIDPrefix = []byte{0x25}
	BorrowPairKeyPrefix   = []byte{0x26}
	LendsKey              = []byte{0x32}
	BorrowsKey            = []byte{0x33}
	BorrowStatsPrefix     = []byte{0x40}
	AuctionParamPrefix    = []byte{0x41}

	AssetToPairMappingKeyPrefix           = []byte{0x20}
	LendForAddressByAssetKeyPrefix        = []byte{0x22}
	UserLendsForAddressKeyPrefix          = []byte{0x23}
	BorrowForAddressByPairKeyPrefix       = []byte{0x24}
	UserBorrowsForAddressKeyPrefix        = []byte{0x27}
	LendIDToBorrowIDMappingKeyPrefix      = []byte{0x28}
	AssetStatsByPoolIDAndAssetIDKeyPrefix = []byte{0x29}
	AssetRatesStatsKeyPrefix              = []byte{0x30}
	LendByUserAndPoolPrefix               = []byte{0x34}
	BorrowByUserAndPoolPrefix             = []byte{0x35}
	DepositStatsPrefix                    = []byte{0x36}
	UserDepositStatsPrefix                = []byte{0x37}
	ReserveDepositStatsPrefix             = []byte{0x38}
	BuyBackDepositStatsPrefix             = []byte{0x39}
	ReservePoolRecordsForBorrowKeyPrefix  = []byte{0x42}
	LendRewardsTrackerKeyPrefix           = []byte{0x43}
	BorrowInterestTrackerKeyPrefix        = []byte{0x44}
)
