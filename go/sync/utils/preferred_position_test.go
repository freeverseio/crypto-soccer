package utils_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/utils"
)

func TestPreferredPosition(t *testing.T) {
	expected := "GK"
	result, _ := utils.PreferredPosition(0x0, 0x0)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "D CR"
	result, _ = utils.PreferredPosition(0x1, 0x3)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "MD LCR"
	result, _ = utils.PreferredPosition(0x5, 0x7)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
}

func TestForwardnessToString(t *testing.T) {
	expected := "GK"
	result, _ := utils.ForwardnessToString(0)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "D"
	result, _ = utils.ForwardnessToString(1)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "M"
	result, _ = utils.ForwardnessToString(2)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "F"
	result, _ = utils.ForwardnessToString(3)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "MF"
	result, _ = utils.ForwardnessToString(4)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "MD"
	result, _ = utils.ForwardnessToString(5)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	_, err := utils.ForwardnessToString(6)
	if err == nil {
		t.Fatal("Expected error with forwardness 6")
	}
}

func TestLeftishnessToString(t *testing.T) {
	expected := ""
	result, _ := utils.LeftishnessToString(0x0)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "R"
	result, _ = utils.LeftishnessToString(0x1)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "C"
	result, _ = utils.LeftishnessToString(0x2)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "CR"
	result, _ = utils.LeftishnessToString(0x3)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "L"
	result, _ = utils.LeftishnessToString(0x4)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	expected = "LCR"
	result, _ = utils.LeftishnessToString(0x7)
	if result != expected {
		t.Fatalf("Expected %v but got %v", expected, result)
	}
	_, err := utils.LeftishnessToString(0x8)
	if err == nil {
		t.Fatal("Expected error with leftishness 8")
	}
}
