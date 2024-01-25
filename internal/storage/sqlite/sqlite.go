package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3" // Импортируем библиотеку для работы с SQLite
)

// --------------------------------------------------------------------------------------

// Описание структуры объекта Storage
type Storage struct {
	db *sql.DB
}

// --------------------------------------------------------------------------------------

// New - Функция-конструктор, которая соберёт и вернёт объект Storage для базы
//
// Принимает: storagePath - "путь к базе данных" (string)
//
// Возвращает: *Storage - указатель на объект Storage и error - ошибку (nil - если ошибки нет)
func New(storagePath string) (*Storage, error) {

	const operation = "storage.sqlite.New"

	// Подключаемся к базе
	db, err := sql.Open("sqlite3", storagePath)

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Подготавливаем запрос
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
	`)

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return nil, fmt.Errorf("#{operation}: #{err}")
	}

	// Выполняем запрос
	_, err = stmt.Exec()

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Возвращаем объект Storage
	return &Storage{db: db}, nil
}

// --------------------------------------------------------------------------------------

// SaveURL - функция-метод сохранения URL в базу по указателю *Storage
//
// Принимает: urlToSave - URL для сохранения (string) и alias - алиас для URL (string)
//
// Возвращает: int64 - ID сохраненного URL и error - ошибку (nil - если ошибки нет)
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {

	const operation = "storage.sqlite.SaveURL"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос
	res, err := stmt.Exec(urlToSave, alias)

	// Проверяем на ошибки
	if err != nil {

		// Если есть ошибка SQLite то, сперва, приводим ошибку к внутреннему типу sqlite3.Error, затем смотрим её extended код
		// и, если он равен ErrConstraintUnique, то возвращаем ошибку ErrURLExists
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrURLExists)
		}

		// В противном случае просто возвращаем ошибку
		return 0, fmt.Errorf("%s: %w", operation, err)

	}

	// Получаем ID
	id, err := res.LastInsertId()
	// LastInsertId поддерживается не всеми БД
	// PostgreSQL	- поддерживается
	// MySQL		- не поддерживается

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return 0, fmt.Errorf("%s: Failed to get last insert id: %w", operation, err)
	}

	// Возвращаем ID и нулевую ошибку
	return id, nil
}

// --------------------------------------------------------------------------------------

// GetURL - функция-метод получения URL из базы по указателю *Storage
//
// Принимает: alias - алиас для URL (string)
//
// Возвращает: сохраненный URL (string) и error - ошибку (nil - если ошибки нет)
func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.sqlite.GetURL"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		// return "", fmt.Errorf("%s: %w", operation, err)
		return "", fmt.Errorf("%s: Prepare statement: %w", operation, err)
	}

	// Подготавливаем переменную в которую положим возвращаемый URL
	var resURL string

	// Выполняем запрос
	err = stmt.QueryRow(alias).Scan(&resURL)

	// Проверяем на ошибки
	if err != nil {

		// Если ничего не получили в ответ - возвращаем ошибку ErrURLNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		// Если какая-то другая ошибка - просто возвращаем её
		return "", fmt.Errorf("%s: Execute statement: %w", operation, err)
	}

	// Возвращаем переменную с полученным URL и нулевую ошибку
	return resURL, nil
}

// --------------------------------------------------------------------------------------

// SaveURL - функция-метод удаления URL из базы по указателю *Storage по принятому alias-у
//
// Принимает: alias - алиас для URL (string)
//
// Возвращает: error - ошибку (nil - если ошибки нет)
func (s *Storage) DeleteURL(alias string) error {

	const operation = "storage.sqlite.DeleteURL"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")

	// Если что-то не так - возвращаем ошибку
	if err != nil {
		return fmt.Errorf("%s: Prepare statement: %w", operation, err)
	}

	// Выполняем запрос
	_, err = stmt.Exec(alias)

	// Проверяем на ошибки
	if err != nil {

		// Если ничего не получили в ответ - возвращаем ошибку ErrURLNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrURLNotFound
		}

		// Если какая-то другая ошибка - просто возвращаем её
		return fmt.Errorf("%s: Execute statement: %w", operation, err)
	}

	// Возвращаем нулевую ошибку
	return nil

}
