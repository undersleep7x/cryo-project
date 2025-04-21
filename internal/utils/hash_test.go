package utils

import (
	"testing"
)

func TestGenerateRef(t *testing.T) {
	key := "secret-key"
	value := "test-data"

	full := GenerateRef(key, value, "test")
	full2 := GenerateRef(key, value, "")
	if len(full) != 64 {
		t.Errorf("Expected full hash, got %d", len(full))
	}
	if full != full2 {
		t.Errorf("Expected consistent hash values, got %s and %s", full, full2)
	}

	soft := GenerateRef(key, value, "soft")
	if len(soft) != 16 {
		t.Errorf("Expected soft hash with len of 16, got %d", len(soft))
	}

	hard := GenerateRef(key, value, "hard")
	if len(hard) != 32 {
		t.Errorf("Expected hard hash with len of 32, got %d", len(hard))
	}
}

func TestBuildRefString(t *testing.T) {
	result := BuildReferenceString("abc", "123", "XYZ")
	expected := "abc123XYZ"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}