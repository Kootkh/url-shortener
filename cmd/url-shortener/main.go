package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

// Запускать командой:
// CONFIG_PATH=./config/local.yaml go run ./cmd/url-shortener/main.go

// --------------------------------------------------------------------------------------

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// --------------------------------------------------------------------------------------

// main - Точка входа программы.
//
// Принимает: None.
//
// Возвращает: None.
func main() {

	// ------------------------
	// INIT CONFIG
	// ------------------------
	// DONE: init config: cleanenv
	cfg := config.MustLoad()

	/**
	*	// !!![DEBUG] DELETE AFTER DEBUG COMPLETE
	*	fmt.Println(cfg)
	*	// !!![DEBUG] DELETE AFTER DEBUG COMPLETE
	**/

	// ------------------------
	// LOGGER
	// ------------------------
	// DONE: init logger: "slog"
	// TODO: change to "Zerolog"
	log := setupLogger(cfg.Env)

	log.Info(
		"Starting url-shortener service...",
		slog.String("env", cfg.Env),
		// slog.String("version", cfg.Version),
		slog.String("version", "123"),
	)
	log.Debug("Debug messages are enabled")
	// log.Error("Error messages are enabled")

	// ------------------------
	// STORAGE
	// ------------------------
	// DONE: init storage: "sqlite"
	// TODO: change to Postgres
	storage, err := sqlite.New(cfg.StoragePath)

	// Проверяем на ошибки
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	// ------------------------
	// ROUTER
	// ------------------------
	// DONE: init router: chi
	// TODO: "chi render"
	// TODO: change to "Fasthttp / mux / gin"
	router := chi.NewRouter()

	// Подключаем к роутеру middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// router.Use(middleware.Logger) // стандартный встроенный логгер golang
	router.Use(mwLogger.New(log)) // написанный нами логгер
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {

		r.Use(middleware.BasicAuth("url-shortener", map[string]string{

			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		// TODO: r.Delete("/url/{alias}", redirect.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Starting server...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// ------------------------
	// RUN SERVER
	// ------------------------
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to run server")
	}

	log.Error("Server stopped")
}

// --------------------------------------------------------------------------------------

// setupLogger - инициализирует и возвращает логгер.
//
// (Вынесли в отдельную функцию, потому что установка зависит от параметра "env" окружения.)
//
// Принимает: env (string, прописанный в конфигурационном файле).
//
// Возвращает: указатель на slog.Logger.
func setupLogger(env string) *slog.Logger {

	// объявляем сам логгер
	var log *slog.Logger

	// создаём логгер в зависимости от параметра "env"
	switch env {

	case envLocal:
		log = setupPrettySlog()
		// log = slog.New(
		//	 slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log

}

// --------------------------------------------------------------------------------------

// setupPrettySlog - инициализирует логгер с помощью slogpretty.
//
// Принимает: None.
//
// Возвращает: указатель на slog.Logger.
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)

}
