package gorez

import "testing"

func TestNothing(t *testing.T) {
	want := "nothing"
	if got := Nothing(); got != want {
		t.Errorf("Nothing() = %q, want %q", got, want)
	}
}
