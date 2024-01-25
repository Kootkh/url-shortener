package save_test

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

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

// --------------------------------------------------------------------------------------

// TestSaveHandler - функция тестирования для обработчика SaveHandler.
//
// Тестирует различные сценарии сохранения URL и alias.
func TestSaveHandler(t *testing.T) {

	// Прописываем структуру и таблицу тестовых кейсов.
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}

	// проходим циклом по таблице тестовых кейсов
	for _, tc := range cases {
		tc := tc

		// Запускаем тестовый сценарий
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			// Создаём объект mock для URLSaver из сгенерированных нами mock-ов
			urlSaverMock := mocks.NewURLSaver(t)

			// Заполняем объект mock
			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")). // Мокаем метод SaveURL
													Return(int64(1), tc.mockError). // Возвращаем ошибку
													Once()                          // Одна проверка
			}

			// Создаем обработчик с указанием нашего "не логирующего логгера" и созданного выше mock объекта.
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// Заполняем содержимое запроса из табличных сценариев описанных выше.
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// Создаем запрос с содержимым
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))

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
			var resp save.Response

			// Через функцию NoError пакета github.com/stretchr/testify проверяем ошибку декодирования JSON ответа
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Сравниваем ошибку ответа с ошибкой тестового табличного кейса
			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}
