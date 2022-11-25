package calculator_test

import (
	"fmt"
	"pivottechschool/packages/calculator"
	"testing"
)

func CalculatorTest(t *testing.T) {
	testCases := []struct {
		num1     int
		num2     int
		operator string
		solution int
	}{
		{12, 21, "+", 33},
		{987, 789, "+", 1776},
		{7, 1, "+", 8},
		{17, 7, "-", 10},
		{752, 892, "-", -140},
		{9, 4, "-", 5},
		{21, 12, "*", 252},
		{100, 1, "*", 100},
		{4, 0, "*", 0},
		{39, 13, "/", 3},
		{900, 150, "/", 6},
		{8, 0, "/", 0},
	}

	for _, cases := range testCases {
		t.Run(fmt.Sprintf("%d%s%d", cases.num1, cases.operator, cases.num2), func(t *testing.T) {
			switch cases.operator {
			case "+":
				if result := calculator.Add(cases.num1, cases.num2); result != cases.solution {
					t.Errorf("result: %d - solution: %d", result, cases.solution)
				}
			case "-":
				if result := calculator.Subtract(cases.num1, cases.num2); result != cases.solution {
					t.Errorf("result: %d - solution: %d", result, cases.solution)
				}
			case "*":
				if result := calculator.Multiply(cases.num1, cases.num2); result != cases.solution {
					t.Errorf("result: %d - solution: %d", result, cases.solution)
				}
			case "/":
				if result, err := calculator.Divide(cases.num1, cases.num2); err != nil {
					if cases.num2 != 0 {
						t.Errorf("%d is not a number", cases.num2)
					}
				} else if result != cases.solution {
					t.Errorf("result: %d - solution: %d", result, cases.solution)
				}
			case "^":
				if result := calculator.Pow(float64(cases.num1), float64(cases.num2)); result != cases.solution {
					t.Errorf("result: %f - solution: %f", result, cases.solution)
				}
			default:
				t.Errorf("invalid: %s", cases.operator)
			}
		})
	}
}
