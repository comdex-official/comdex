package types

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetShareValue multiplies with truncation by receiving int amount and decimal ratio and returns int result.
func GetShareValue(amount sdk.Int, ratio sdk.Dec) sdk.Int {
	return amount.ToDec().MulTruncate(ratio).TruncateInt()
}

type StrIntMap map[string]sdk.Int

// AddOrSet Set when the key not existed on the map or add existed value of the key.
func (m StrIntMap) AddOrSet(key string, value sdk.Int) {
	if _, ok := m[key]; !ok {
		m[key] = value
	} else {
		m[key] = m[key].Add(value)
	}
}

// DateRangesOverlap returns true if two date ranges overlap each other.
// End time is exclusive and start time is inclusive.
// func DateRangesOverlap(startTimeA, endTimeA, startTimeB, endTimeB time.Time) bool {
// 	return startTimeA.Before(endTimeB) && endTimeA.After(startTimeB)
// }

// DateRangeIncludes returns true if the target date included on the start, end time range.
// End time is exclusive and start time is inclusive.
func DateRangeIncludes(startTime, endTime, targetTime time.Time) bool {
	return endTime.After(targetTime) && !startTime.After(targetTime)
}

// ParseDec is a shortcut for sdk.MustNewDecFromStr.
func ParseDec(s string) sdk.Dec {
	return sdk.MustNewDecFromStr(s)
}

// ParseCoin parses and returns sdk.Coin.
func ParseCoin(s string) sdk.Coin {
	coin, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		panic(err)
	}
	return coin
}

// ParseCoins parses and returns sdk.Coins.
func ParseCoins(s string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}

// ParseTime parses and returns time.Time in time.RFC3339 format.
func ParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// DecApproxEqual returns true if a and b are approximately equal,
// which means the diff ratio is equal or less than 0.1%.
func DecApproxEqual(a, b sdk.Dec) bool {
	if b.GT(a) {
		a, b = b, a
	}
	return a.Sub(b).Quo(a).LTE(sdk.NewDecWithPrec(1, 3))
}

// RandomInt returns a random integer in the half-open interval [min, max).
func RandomInt(r *rand.Rand, min, max sdk.Int) sdk.Int {
	return min.Add(sdk.NewIntFromBigInt(new(big.Int).Rand(r, max.Sub(min).BigInt())))
}

// RandomDec returns a random decimal in the half-open interval [min, max).
func RandomDec(r *rand.Rand, min, max sdk.Dec) sdk.Dec {
	return min.Add(sdk.NewDecFromBigIntWithPrec(new(big.Int).Rand(r, max.Sub(min).BigInt()), sdk.Precision))
}

// SafeMath runs f in safe mode, which means that any panics occurred inside f
// gets caught by recover() and if the panic was an overflow, onOverflow is run.
// Otherwise, if the panic was not an overflow, then SafeMath will re-throw
// the panic.
func SafeMath(f, onOverflow func()) {
	defer func() {
		if r := recover(); r != nil {
			if IsOverflow(r) {
				onOverflow()
			} else {
				panic(r)
			}
		}
	}()
	f()
}

// IsOverflow returns true if the panic value can be interpreted as an overflow.
func IsOverflow(r interface{}) bool {
	switch r := r.(type) {
	case string:
		s := strings.ToLower(r)
		return strings.Contains(s, "overflow") || strings.HasSuffix(s, "out of bound")
	}
	return false
}

// This function lets you run the function f. In case of panic recovery is done
// if error occurs it is logged into the logger.
// further modifications can me made to avoid any state changes in case if error is returned by f -
// eg, revert state change if error returned by f else work as normal
func ApplyFuncIfNoError(ctx sdk.Context, f func(ctx sdk.Context) error) (err error) {
	// Add a panic safeguard
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			PrintPanicRecoveryError(ctx, recoveryError)
			err = errors.New("panic occurred during execution")
		}
	}()
	cacheCtx, writeCache := ctx.CacheContext()
	err = f(cacheCtx)
	if err == nil {
		// write state to the underlying multi-store
		writeCache()
	} else {
		ctx.Logger().Error(err.Error())
	}
	return err
}

// PrintPanicRecoveryError error logs the recoveryError, along with the stacktrace, if it can be parsed.
// If not emits them to stdout.
func PrintPanicRecoveryError(ctx sdk.Context, recoveryError interface{}) {
	errStackTrace := string(debug.Stack())
	switch e := recoveryError.(type) {
	case string:
		ctx.Logger().Error("Recovering from (string) panic: " + e)
	case runtime.Error:
		ctx.Logger().Error("recovered (runtime.Error) panic: " + e.Error())
	case error:
		ctx.Logger().Error("recovered (error) panic: " + e.Error())
	default:
		ctx.Logger().Error("recovered (default) panic. Could not capture logs in ctx, see stdout")
		fmt.Println("Recovering from panic ", recoveryError)
		debug.PrintStack()
		return
	}
	ctx.Logger().Error("stack trace: " + errStackTrace)
}
