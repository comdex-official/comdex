package types

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (d Denom) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Equal implements equal interface
func (d Denom) Equal(d1 *Denom) bool {
	return d.BaseDenom == d1.BaseDenom &&
		d.SymbolDenom == d1.SymbolDenom &&
		d.Exponent == d1.Exponent
}

// DenomList is array of Denom
type DenomList []Denom

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	for _, d := range dl {
		out += d.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// Contains checks whether or not a SymbolDenom (e.g. CMDX) is in the DenomList
func (dl DenomList) Contains(symbolDenom string) bool {
	for _, d := range dl {
		if strings.EqualFold(d.SymbolDenom, symbolDenom) {
			return true
		}
	}
	return false
}

// ContainDenoms checks if d is a subset of dl
func (dl DenomList) ContainDenoms(d DenomList) bool {
	contains := make(map[string]struct{})

	for _, denom := range dl {
		contains[denom.String()] = struct{}{}
	}

	for _, denom := range d {
		if _, found := contains[denom.String()]; !found {
			return false
		}
	}

	return true
}

// Normalize updates all the SymbolDenom strings to use all caps.
func (dl DenomList) Normalize() DenomList {
	for i := range dl {
		dl[i].SymbolDenom = strings.ToUpper(dl[i].SymbolDenom)
	}
	return dl
}
