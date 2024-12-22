package pkg

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 97)
	if result != 99 {
		t.Errorf("Expected 8, got %f", result)
	}
}

func TestSub(t *testing.T) {
	result := Sub(99, 66)
	if result != 33 {
		t.Errorf("Expected 6, got %f", result)
	}
}
