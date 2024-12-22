package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

type requestBody struct {
	Expression string `json:"expression"`
}

type responseBodySuccess struct {
	Result string `json:"result"`
}

type responseBodyError struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(responseBodyError{Error: "Method not allowed"})
			return
		}

		var req requestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(responseBodyError{Error: "Expression is not valid"})
			fmt.Println(err)
			return
		}

		result, err := Calc(req.Expression)
		if err != nil {
			switch err.Error() {
			case "недопустимый символ", "пустое выражение", "несоответствующая скобка", "ошибка в выражении", "ошибка - недостаточно данных для вычисления", "ошибка - деление на ноль", "ошибка - неизвестная операция":
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(responseBodyError{Error: "Expression is not valid"})
				fmt.Println(err)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(responseBodyError{Error: "Internal server error"})
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBodySuccess{Result: floatToString(result)})
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
