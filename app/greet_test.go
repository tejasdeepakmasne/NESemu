package main

import "testing"

func TestGreet(t *testing.T) {
	got := greet("Gopher")
	want := "Hello, Gopher!"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}