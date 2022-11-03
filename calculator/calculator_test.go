package calculator

import (
	"github.com/alexhcole95/pivottechschool/calculator"
	"testing"
)

func TestAdd(t *testing.T) {
	result := calculator.Add(1, 2)
	if result != 3 {
		t.Error("Expected 3, got ", result)
	}
}

func TestSubtract(t *testing.T) {
	result := calculator.Subtract(5, 3)
	if result != 2 {
		t.Error("Expected 2, got ", result)
	}
}

func TestMultiply(t *testing.T) {
	result := calculator.Multiply(5, 3)
	if result != 15 {
		t.Error("Expected 15, got ", result)
	}
}

func TestDivide(t *testing.T) {
	result, err := calculator.Divide(6, 3)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if result != 2 {
		t.Error("Expected 2, got ", result)
	}
}
