package hasher

import (
	"testing"
)

func TestHasherSHA_Hash(t *testing.T) {
	hasher := NewHasher("some_salt")

	result, err := hasher.Hash("string")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	expected := "736f6d655f73616c74473287f8298dba7163a897908958f7c0eae733e25d2e027992ea2edc9bed2fa8"
	if result != expected {

		t.Errorf("Expected %v, but got %v", expected, result)
	}

	result2, err := hasher.Hash("another_test_string")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if result == result2 {
		t.Errorf("Expected different hashes, but got the same: %v", result)
	}

}
