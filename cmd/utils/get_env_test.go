package utils

import "testing"

func TestGetEnv(t *testing.T) {
	got := GetEnv("NotExistingEnv", "default")
	want := "default"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
