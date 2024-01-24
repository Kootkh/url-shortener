package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		// Создаём копию логгера, добавляя подсказку что это "компонент - middleware/logger"
		// Это параметр который будет выводиться с каждой строчкой логов.
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		// Информируем о включении хэндлера middleware логгера
		log.Info("Logger middleware enabled")

		// Функция обработки запроса
		fn := func(w http.ResponseWriter, r *http.Request) {

			// Выполняется до обработки запроса
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote-addr", r.RemoteAddr),
				slog.String("user-agent", r.UserAgent()),
				slog.String("request-id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			// Выполняется после обработки запроса
			// (После всех цепочек "next" middleware)
			// Информируем об окончании обработки
			defer func() {
				entry.Info("Request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		// Возвращаем обработчик, приводя к нужному типу fn
		return http.HandlerFunc(fn)

	}

}
