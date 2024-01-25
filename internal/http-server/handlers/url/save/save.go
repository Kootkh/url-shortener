package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
	//"url-shortener/internal/storage"
)

// --------------------------------------------------------------------------------------

// Описание структуры поступающих запросов (json объект).
type Request struct {

	// Struct tag 'validate' - даёт пакету валидатора информацию о:
	// required - обязательности поля
	// url - тип валидации
	URL string `json:"url" validate:"required,url"`
	// omitempty установлено		- если значение отсутствует, то в итоговом json объекте строку не пишем
	// omitempty не установлено	- если значение отсутствует, то в итоговом json объекте строку пишем, но с пустым значением
	Alias string `json:"alias,omitempty"`
}

// --------------------------------------------------------------------------------------

// Описание структуры возвращаемого ответа (json объект).
type Response struct {
	// Подключаем структуру Response из модуля response
	resp.Response
	// omitempty установлено		- если значение отсутствует, то в итоговом json объекте строку не пишем
	// omitempty не установлено	- если значение отсутствует, то в итоговом json объекте строку пишем, но с пустым значением
	Alias string `json:"alias,omitempty"`
}

// --------------------------------------------------------------------------------------

const aliasLength = 6 // длина алиаса при генерации

// --------------------------------------------------------------------------------------

// Генерация мока для URLSaver
// go install github.com/vektra/mockery/v2@v2.40.1
//go:generate go run github.com/vektra/mockery/v2@v2.40.1 --name=URLSaver

// URLSaver - Интерфейс стораджа.
// Описываем по месту использования.
type URLSaver interface {
	// Сигнатура этого метода должна дублировать сигнатуру метода в сторадже.
	// internal > storage > sqlite > sqlite.go
	SaveURL(urlToSave string, alias string) (int64, error)
}

// --------------------------------------------------------------------------------------

// New - Функция "конструктор" для обработчика. Возвращает функцию http.HandlerFunc, которая обрабатывает запросы на запись URL.
//
// Принимает: указатель на slog.Logger , URLSaver
//
// Возвращает: функцию http.HandlerFunc.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {

	// Возвращаем функцию http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {

		// Инициализируем log
		const operation = "handlers.url.save.New"

		// Создаем лог для хэндлера
		log = log.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Инициализируем структуру Request JSON объекта входящего запроса
		var req Request

		// Декодируем JSON объект входящего запроса
		err := render.DecodeJSON(r.Body, &req)

		// Обрабатываем, при наличии, ошибку получения запроса с пустым телом
		if errors.Is(err, io.EOF) {

			// Пишем ошибку в log.Error
			log.Error("request body is empty")

			// Рендерим ответ - Возвращаем ошибку клиенту в виде JSON объекта
			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		// Обрабатываем, при наличии, ошибку декодирования входящего запроса
		if err != nil {

			// Пишем ошибку в log.Error
			log.Error("failed to decode request body", sl.Err(err))

			// Рендерим ответ - Возвращаем ошибку клиенту в виде JSON объекта
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		// Добавляем в log.Info информацию об удачном декодировании объекта Request
		log.Info("request body decoded", slog.Any("request", req))

		// Валидируем декодированный объект Request входящего запроса через проверку на ошибку создания нового объекта валидатора
		// с помощью пакета "github.com/go-playground/validator/v10"
		if err := validator.New().Struct(req); err != nil {

			// Приводим ошибку валидации к структуре ValidationError
			validateErr := err.(validator.ValidationErrors)

			// Пишем ошибку (в чистом виде, без изменений) в log.Error
			log.Error("invalid request body", sl.Err(err))

			// Рендерим ответ - Возвращаем ошибку клиенту в виде JSON объекта сформированного функцией ValidationError
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// проверяем на наличие в запросе параметра "alias"
		alias := req.Alias

		if alias == "" {
			// если параметр "alias" отсутствует, то генерируем его из случайных символов с помощь пакета
			alias = random.NewRandomString(aliasLength)
		}

		// Сохраняем URL через переданный нам urlSaver в хранилище
		id, err := urlSaver.SaveURL(req.URL, alias)

		// Обрабатываем ошибку когда url уже существует
		if errors.Is(err, storage.ErrURLExists) {

			// Добавляем в log.Info информацию об ошибке
			log.Info("url already exists", slog.String("url", req.URL))

			// Рендерим ответ - Возвращаем ошибку клиенту в виде JSON объекта
			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		// Обрабатываем остальные ошибки
		if err != nil {

			// Пишем ошибку (в чистом виде, без изменений) в log.Error
			log.Error("failed to save url", sl.Err(err))

			// Рендерим ответ - Возвращаем ошибку клиенту в виде JSON объекта
			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		// Если ошибок нет, то пишем сообщение об успешном сохранении в log.Info
		log.Info("saved url", slog.Int64("id", id))

		// И возвращаем ответ клиенту
		responseOK(w, r, alias)

	}
}

// --------------------------------------------------------------------------------------

// responseOK renders a JSON response with the provided alias.
//
// Parameters:
//
//	w - http.ResponseWriter
//	r - *http.Request
//	alias - string
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
