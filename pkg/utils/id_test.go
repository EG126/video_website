package utils

import (
	"strings"
	"testing"
)

func TestA_GenerateID_NotInitialized(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("GenerateID() should panic when not initialized")
		}
	}()

	GenerateID()
}

func TestInitSnowflake_Success(t *testing.T) {
	err := InitSnowflake("2024-01-01", 1)
	if err != nil {
		t.Errorf("InitSnowflake() error = %v", err)
	}
}

func TestInitSnowflake_InvalidStartTime(t *testing.T) {
	err := InitSnowflake("invalid-date", 1)
	if err == nil {
		t.Error("InitSnowflake() should return error for invalid date")
	}
}

func TestGenerateID_AfterInit(t *testing.T) {
	err := InitSnowflake("2024-01-01", 1)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}

	id := GenerateID()
	if id == "" {
		t.Error("GenerateID() returned empty string")
	}

	if len(id) == 0 {
		t.Error("GenerateID() returned empty ID")
	}
}

func TestGenerateID_Unique(t *testing.T) {
	err := InitSnowflake("2024-01-01", 1)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}

	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("GenerateID() should return unique IDs")
	}
}

func TestGenerateID_Format(t *testing.T) {
	err := InitSnowflake("2024-01-01", 1)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}

	id := GenerateID()

	if strings.Contains(id, " ") {
		t.Error("GenerateID() should not contain spaces")
	}

	if len(id) < 18 {
		t.Errorf("GenerateID() ID length should be at least 18, got %d", len(id))
	}
}

func TestInitSnowflake_DifferentMachineID(t *testing.T) {
	err := InitSnowflake("2024-01-01", 1)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}
	id1 := GenerateID()

	err = InitSnowflake("2024-01-01", 2)
	if err != nil {
		t.Fatalf("InitSnowflake() error = %v", err)
	}
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("Different machine IDs should generate different IDs")
	}
}
