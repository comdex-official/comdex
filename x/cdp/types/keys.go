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

var sep = []byte(":")

var (
	CdpIDKeyPrefix             = []byte{0x01}
	CdpKeyPrefix               = []byte{0x02}
	CollateralRatioIndexPrefix = []byte{0x03}
	CdpIDKey                   = []byte{0x04}
	DebtDenomKey               = []byte{0x05}
	GovDenomKey                = []byte{0x06}
	DepositKeyPrefix           = []byte{0x07}
	PrincipalKeyPrefix         = []byte{0x08}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetCdpIDBytes(cdpID uint64) (cdpIDBz []byte){
	cdpIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(cdpIDBz, cdpID)
	return
}

// GetCdpIDFromBytes returns cdpID in uint64 format from a byte array
func GetCdpIDFromBytes(bz []byte) (cdpID uint64) {
	return binary.BigEndian.Uint64(bz)
}

func CdpKey(prefix byte, cdpID uint64) []byte {
	return createKey([]byte{prefix}, sep, GetCdpIDBytes(cdpID))
}

func createKey(bytes ...[]byte) (r []byte) {
	for _, b := range bytes {
		r = append(r, b...)
	}
	return
}
