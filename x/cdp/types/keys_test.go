package types

import (
	cmp "github.com/google/go-cmp/cmp"
	"testing"
)

func TestCdpKey(t *testing.T) {
	get := CdpKey(1 )
	want := []byte{2, 0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}

func TestGetCdpIDBytes(t *testing.T) {
	get := GetCdpIDBytes(1)
	want := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}

func TestGetCdpIDFromBytes(t *testing.T) {
	get := GetCdpIDBytes(1)
	want := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}