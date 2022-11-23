package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pivottechschool/packages/calculator"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("CALCULATOR")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("ENTER CALCULATION (EX: 1+1): ")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		handleInput(str)
	}
}

func handleInput(str string) {
	var sepStr []string
	str = strings.Replace(str, "\n", "", -1)
	operator, err := checkContains(str)
	if checkPrintError(err) {
		return
	}

	sepStr = strings.Split(str, operator)
	if len(sepStr) > 3 {
		fmt.Println(errors.New("error: too many arguments. please try again"))
	}

	for i := range sepStr {
		sepStr[i] = strings.Replace(sepStr[i], " ", "", -1)
		if i > 1 {
			fmt.Println(errors.New("error: only one calculation allowed. please try again"))
			return
		}
	}

	a, err := strconv.Atoi(sepStr[0])
	if checkPrintError(err) {
		return
	}

	b, err := strconv.Atoi(sepStr[1])
	if checkPrintError(err) {
		return
	}

	var result int
	switch operator {
	case "+":
		result = calculator.Add(a, b)
	case "-":
		result = calculator.Subtract(a, b)
	case "*":
		result = calculator.Multiply(a, b)
	case "/":
		result, err = calculator.Divide(a, b)
		if checkPrintError(err) {
			return
		}
	default:
		fmt.Println(errors.New("error: unknown operator. please try again"))
	}
	fmt.Printf("RESULT: %d\n", result)
}

func checkContains(str string) (string, error) {
	operators := []string{"+", "-", "*", "/"}
	var containsOperator []string
	for _, v := range operators {
		if strings.Contains(str, v) {
			containsOperator = append(containsOperator, v)
		}
	}
	switch len(containsOperator) {
	case 0:
		return "", errors.New("error: no operator. please try again")
	case 1:
		return containsOperator[0], nil
	default:
		return "", errors.New("error: more than one operator. please try again")
	}
}

func checkPrintError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
