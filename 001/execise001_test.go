package main

import "testing"

func TestEx001(t *testing.T) {
	want := "7,14,21,28"
	got := Ex001(1, 28)

	if got != want {
		t.Errorf("Ex001 = %v, want %v", got, want)
	}

}
