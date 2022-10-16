package utils

import "testing"

func TestConcat(t *testing.T) {
	got := ConcatStr("Hello", " ", "World")
	want := "Hello World"

	if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}
