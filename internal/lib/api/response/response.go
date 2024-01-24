package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// --------------------------------------------------------------------------------------

// Response - Описание части структуры возвращаемого ответа для всех обработчиков (json объект).
type Response struct {
	Status string `json:"status"` // { Error | Ok }
	// omitempty установлено - если значение отсутствует, то в итоговом json объекте параметр не пишем
	// omitempty не установлено - если значение отсутствует, то в итоговом json объекте параметр пишем, но с пустым значением
	Error string `json:"error,omitempty"`
}

// --------------------------------------------------------------------------------------

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// --------------------------------------------------------------------------------------

// OK - Возвращает часть структуры ответа Response (JSON объект) с параметром "Status: StatusOK".
//
// Принимает: None.
//
// Возвращает: Response (JSON объект)
//
//	{
//	   Status: StatusOK,
//	}
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// --------------------------------------------------------------------------------------

// Error - Возвращает часть структуры ответа Response (JSON объект) с параметрами "Status: StatusError" и "Error: msg".
//
// Принимает: msg (string) - сообщение об ошибке.
//
// Возвращает: Response (JSON объект)
//
//	{
//	   Status: StatusError,
//	   Error: msg,
//	}
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// --------------------------------------------------------------------------------------

// ValidationError - функция генерации структуры ответа Response для отправки ошибок валидации.
//
// Принимает: errs (validator.ValidationErrors) - список ошибок валидатора.
//
// Возвращает: часть структуры ответа Response со списком ошибок валидации (JSON объект).
//
//	{
//	   Status: StatusError,
//	   Error:  strings.Join(errMsgs, ", "),
//	}
func ValidationError(errs validator.ValidationErrors) Response {

	// Ошибки валидации будем собирать в массив строк errMsgs.
	var errMsgs []string

	// Перебираем список ошибок валидации errs и формируем ответ клиенту
	for _, err := range errs {

		// Проверяем ошибку по тэгу
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	// Возвращаем ответа Response клиенту (JSON объект)
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
