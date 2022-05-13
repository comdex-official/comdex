package types

// event types.
const (
	TypeEvtLockTokens      = "lock_tokens"
	TypeEvtAddTokensToLock = "add_tokens_to_lock"
	TypeEvtBeginUnlock     = "begin_unlock"

	AttributeLockID       = "lock_id"
	AttributeUnlockID     = "unlock_id"
	AttributeLockOwner    = "owner"
	AttributeLockAmount   = "amount"
	AttributeUnLockAmount = "unlock_amount"
	AttributeLockDuration = "duration"
)
