package main

import "testing"

func TestBasicMath(t *testing.T) {
	// A simple unit test to prove the CI pipeline works
	expected := 2
	result := 1 + 1

	if result != expected {
		t.Errorf("Expected %d, but got %d. Math is broken!", expected, result)
	}
}