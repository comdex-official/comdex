package types

import "testing"

type Test struct {
	out error
}

var tests = []Test{
	{ Spot.IsValid()},
	{liquidation.IsValid()},
}

func TestPricefeedType_IsValid(t *testing.T) {
	for _, test := range tests {
	if test.out != nil {
			t.Error()
		}
	}
}