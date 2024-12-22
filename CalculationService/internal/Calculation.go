package internal

import (
	"errors"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	if expression == "" {
		return 0, errors.New("пустое выражение")
	}

	var numbers []float64
	var ops []rune

	calculate := func() error {
		if len(numbers) < 2 || len(ops) == 0 {
			return errors.New("ошибка - недостаточно данных для вычисления")
		}
		b := numbers[len(numbers)-1]
		a := numbers[len(numbers)-2]
		op := ops[len(ops)-1]

		numbers = numbers[:len(numbers)-2]
		ops = ops[:len(ops)-1]

		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return errors.New("ошибка - деление на ноль")
			}
			result = a / b
		default:
			return errors.New("ошибка - неизвестная операция")
		}
		numbers = append(numbers, result)
		return nil
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if char >= '0' && char <= '9' {
			start := i
			for i < len(expression) && (expression[i] >= '0' && expression[i] <= '9' || expression[i] == '.') {
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, err
			}
			numbers = append(numbers, num)
			i--
		} else if char == '(' {
			ops = append(ops, char)
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if err := calculate(); err != nil {
					return 0, err
				}
			}
			if len(ops) == 0 {
				return 0, errors.New("несоответствующая скобка")
			}
			ops = ops[:len(ops)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(char) {
				if err := calculate(); err != nil {
					return 0, err
				}
			}
			ops = append(ops, char)
		} else {
			return 0, errors.New("недопустимый символ")
		}
	}

	for len(ops) > 0 {
		if err := calculate(); err != nil {
			return 0, err
		}
	}

	if len(numbers) != 1 {
		return 0, errors.New("ошибка в выражении")
	}

	return numbers[0], nil
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}
