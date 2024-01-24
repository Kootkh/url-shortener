package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Описываем часть возвращаемого ответа для всех хэндлеров  (json объект)
type Response struct {
	Status string `json:"status"` // { Error | Ok }
	// omitempty установлено		- если значение отсутствует, то в итоговом json объекте параметр не пишем
	// omitempty не установлено	- если значение отсутствует, то в итоговом json объекте параметр пишем, но с пустым значением
	Error string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// ValidationError - Генерация ответа для ошибок валидации.
// На вход принимает список ошибок валидатора.
// Возвращает структуру Response.
func ValidationError(errs validator.ValidationErrors) Response {

	// Собираем сообщения об ошибках валидации
	var errMsgs []string

	// Перебираем ошибки валидации и формируем ответ клиенту
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

	// Возвращаем ответ RESPONSE клиенту
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
