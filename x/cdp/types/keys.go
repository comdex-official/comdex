package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "cdp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cdp"
)

var (
	TypeMsgCreateCDPRequest = ModuleName + ":create_cdp"
	TypeMsgDepositRequest   = ModuleName + ":deposit"
	TypeMsgWithdrawRequest  = ModuleName + ":withdraw"
	TypeMsgDrawDebtRequest  = ModuleName + ":draw_debt"
	TypeMsgRepayDebtRequest = ModuleName + ":repay_debt"
	TypeMsgLiquidateRequest = ModuleName + ":liquidate"
)

var (
	CdpIDIndexKeyPrefix = []byte{0x01}
	CdpKeyPrefix        = []byte{0x02}
	CdpIDKey            = []byte{0x03}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetCdpIDBytes(cdpID uint64) (cdpIDBz []byte) {
	cdpIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(cdpIDBz, cdpID)
	return
}

// GetCdpIDFromBytes returns cdpID in uint64 format from a byte array
func GetCdpIDFromBytes(bz []byte) (cdpID uint64) {
	return binary.BigEndian.Uint64(bz)
}

func CdpKey(cdpID uint64) []byte {
	return createKey(CdpKeyPrefix, GetCdpIDBytes(cdpID))
}

func createKey(bytes ...[]byte) (r []byte) {
	for _, b := range bytes {
		r = append(r, b...)
	}
	return
}
