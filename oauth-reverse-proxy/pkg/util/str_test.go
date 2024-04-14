package util

import (
	"testing"
)

func TestBytes(t *testing.T) {
	got := Bytes(5)
	if len(got) != 5 {
		t.Errorf("%v", got)
	}
}

func TestRandomString(t *testing.T) {
	got := RandomString(5)
	if len(got) != 5 {
		t.Errorf("%v", got)
	}
}
