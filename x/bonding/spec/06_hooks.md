<!--
order: 6
-->

# Hooks

In this section we describe the "hooks" that `bonding` module provide for other modules.

## Tokens Locked

On lock/unlock events, bonding module execute hooks for other modules to make following actions.

```go
  OnTokenLocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time)
  OnTokenUnlocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time)
```
