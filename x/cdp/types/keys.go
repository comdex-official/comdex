package types

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
	TypeMsgDepositRequest = ModuleName + ":deposit"
	TypeMsgWithdrawRequest = ModuleName + ":withdraw"
	TypeMsgDrawDebtRequest = ModuleName + ":draw_debt"
	TypeMsgRepayDebtRequest = ModuleName + ":repay_debt"
	TypeMsgLiquidateRequest = ModuleName + ":liquidate"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
