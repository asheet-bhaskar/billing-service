package utils

import (
	"testing"

	"github.com/google/uuid"
)

func Test_GetNewUUID(t *testing.T) {
	id := GetNewUUID()
	if len(id) != 36 {
		t.Error("length should be 36")
	}
}

func Test_IsValidUUID(t *testing.T) {
	id := GetNewUUID()
	if err := uuid.Validate(id); err != nil {
		t.Error("invalid uuid")
	}
}

func Test_RandomString(t *testing.T) {
	s := RandomString(10)

	if len(s) != 10 {
		t.Error("length should be 10")
	}
}
