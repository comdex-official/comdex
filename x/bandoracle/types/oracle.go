package types

import "encoding/binary"

type (
	OracleScriptID  uint64
	OracleRequestID int64
)

// int64ToBytes convert int64 to a byte slice
func int64ToBytes(num int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(num))
	return result
}
