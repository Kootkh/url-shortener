package slogdiscard

import (
	"context"
	"log/slog"
)

// --------------------------------------------------------------------------------------

// NewDiscardLogger - конструктор, который собирает логгер и возвращает его.
//
// Принимает: None.
//
// Возвращает: Указатель на *slog.Logger
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// --------------------------------------------------------------------------------------

// DiscardHandler - обработчик, который просто игнорирует запись журнала
type DiscardHandler struct{}

// --------------------------------------------------------------------------------------

// NewDiscardHandler - билдер для обработчика DiscardHandler.
//
// Принимает: None.
//
// Возвращает: Указатель на DiscardHandler.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// --------------------------------------------------------------------------------------

// Метод обработчика  DiscardHandler, который просто игнорирует запись журнала.
//
// Принимает: context.Context , slog.Record
//
// Возвращает: nil | error.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

// --------------------------------------------------------------------------------------

// Метод обработчика DiscardHandler, который возвращает тот же обработчик, так как нет атрибутов для сохранения.
//
// Принимает: _ []slog.Attr
//
// Возвращает: slog.Handler
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

// --------------------------------------------------------------------------------------

// Метод обработчика DiscardHandler, который возвращает тот же обработчик, так как нет группы для сохранения.
//
// Принимает: _ string
//
// Возвращает: slog.Handler
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

// --------------------------------------------------------------------------------------

// Метод обработчика DiscardHandler, который всегда возвращает false, так как запись журнала игнорируется.
//
// Принимает: _ context.Context, _ slog.Level
//
// Возвращает: false (bool)
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}
