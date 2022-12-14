package utils

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	got := GetEnv("NotExistingEnv", "default")
	want := "default"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestSetGetEnv(t *testing.T) {
	os.Setenv("Balls", "BIG")
	got := GetEnv("Balls", "SMOL")
	want := "BIG"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
