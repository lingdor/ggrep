package util

import (
	"testing"
)

func TestFullIntBinary(t *testing.T) {

	binary := FullIntBinary(1)
	if binary != 1 {
		t.Failed()
	}

	binary = FullIntBinary(2)
	if binary != 4 {
		t.Failed()
	}
}
