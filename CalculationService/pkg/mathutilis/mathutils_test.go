package mathutilis

import "testing"

func TestAdd(t *testing.T) {
	result := Add(7, 8)
	expected := 15

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}
