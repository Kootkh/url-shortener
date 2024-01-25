package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// --------------------------------------------------------------------------------------

// URLGetter - интерфейс для получения URL по принятому alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// --------------------------------------------------------------------------------------

// NEW - Функция-конструктор для обработчика редиректа
//
// Принимает: log - как указатель на экземпляр логгера slog.Logger и urlGetter - как интерфейс URLGetter.
//
// Возвращает: http.HandlerFunc - функция-обработчик редиректа
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Принимаем константе имя текущей функции.
		const operation = "handlers.url.redirect.New"

		// Создаём объект для логгера.
		log := log.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// С помощью библиотеки сhi получаем alias из запроса router.Get("/{alias}", redirect.New(log, storage)).
		alias := chi.URLParam(r, "alias")

		// Проверяем alias на пустоту.
		if alias == "" {

			// Выводим в лог сообщение об ошибке.
			log.Info("alias is empty")

			// Отправляем ответ с ошибкой (JSON формат).
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		// Обращаемся к нашему sqliteStorage-у через urlGetter и запрашиваем URL по alias-у
		resURL, err := urlGetter.GetURL(alias)

		// Обрабатываем ошибку ErrURLNotFound (если она есть)
		if errors.Is(err, storage.ErrURLNotFound) {

			// Выводим в лог сообщение об ошибке с указанием alias-а вызвавшего ошибку.
			log.Info("url not found", "alias", alias)

			// Отправляем ответ с ошибкой (JSON формат).
			render.JSON(w, r, resp.Error("not found"))

			return
		}

		// Обрабатываем остальные ошибку
		if err != nil {

			// Выводим в лог сообщение об ошибке.
			log.Error("failed to get url", sl.Err(err))

			// Отправляем ответ с ошибкой (JSON формат) без указания подробностей (security reasons).
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		// Выводим в лог сообщение об успешном получении URL с указанием полученного URL
		log.Info("got url", slog.String("url", resURL))

		// С помощью стандартной библиотеки http делаем переадресацию на найденный URL с указанием статуса 302 (StatusFound)
		// 301 (StatusMovedPermanently) - ресурс перемещён навсегда. Браузеры могут закэшировать ссылку и, если по alias-у поменяется URL или будет удалён, будут редиректить на старый URL.
		// 302 (StatusFound) - по данному статусу браузер не кэширует информацию.
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
