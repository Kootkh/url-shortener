package delete_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

// --------------------------------------------------------------------------------------

// TestDeleteHandler - функция тестирования для обработчика DeleteHandler.
//
// Тестирует различные сценарии удаления alias.
func TestDeleteHandler(t *testing.T) {

	// Прописываем структуру и таблицу тестовых кейсов.
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "field ALIAS is a required field",
		},
		{
			name:      "Invalid ALIAS",
			alias:     "some_invalid_ALIAS",
			respError: "field ALIAS is not a valid ALIAS",
		},
		{
			name:      "DeleteALIAS Error",
			alias:     "test_alias",
			respError: "failed to delete alias",
			mockError: errors.New("unexpected error"),
		},
	}

	// проходим циклом по таблице тестовых кейсов
	for _, tc := range cases {
		tc := tc

		// Запускаем тестовый сценарий
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			// Создаём объект mock для ALIASDeleter из сгенерированных нами mock-ов
			aliasDeleterMock := mocks.NewALIASDeleter(t)

			// Заполняем объект mock
			if tc.respError == "" || tc.mockError != nil {
				aliasDeleterMock.On("DeleteALIAS", tc.alias, mock.AnythingOfType("string")). // Мокаем метод DeleteALIAS
														Return(int64(1), tc.mockError). // Возвращаем ошибку
														Once()                          // Одна проверка
			}

			// Создаем обработчик с указанием нашего "не логирующего логгера" и созданного выше mock объекта.
			handler := delete.New(slogdiscard.NewDiscardLogger(), aliasDeleterMock)

			// Заполняем содержимое запроса из табличных сценариев описанных выше.
			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			// Создаем запрос с содержимым
			req, err := http.NewRequest(http.MethodPost, "/delete", bytes.NewReader([]byte(input)))

			// Обрабатываем ошибку через функцию NoError пакета github.com/stretchr/testify
			// require - сфэйлить по месту вызова и вывалиться, когда нет смысла продолжать обработку
			// assert - добавить информацию о фэйле в вывод и продолжить обработку
			require.NoError(t, err)

			// С помощью пакета httptest из стандартной библиотеки GO создаём объект rr (response recorder) для записи в него ответа нашего обработчика
			rr := httptest.NewRecorder()

			// Запускаем запрос
			handler.ServeHTTP(rr, req)

			// Проверяем код возвращаемого ответа
			require.Equal(t, rr.Code, http.StatusOK)

			// Проверяем то, что было записано в тело ответа
			body := rr.Body.String()

			// Создаём объект resp в который запишем ответ Response
			var resp delete.Response

			// Через функцию NoError пакета github.com/stretchr/testify проверяем ошибку декодирования JSON ответа
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Сравниваем ошибку ответа с ошибкой тестового табличного кейса
			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}
