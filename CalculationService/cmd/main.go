package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResponseSuccess struct {
	Result string `json:"result"`
}

type ResponseError struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ResponseError{Error: "Method not allowed"})
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ResponseError{Error: "Expression is not valid"})
			fmt.Println(err)
			return
		}

		result, err := Calc(req.Expression)
		if err != nil {
			switch err.Error() {
			case "недопустимый символ", "пустое выражение", "несоответствующая скобка", "ошибка в выражении", "ошибка - недостаточно данных для вычисления", "ошибка - деление на ноль", "ошибка - неизвестная операция":
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(ResponseError{Error: "Expression is not valid"})
				fmt.Println(err)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ResponseError{Error: "Internal server error"})
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseSuccess{Result: floatToString(result)})
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
