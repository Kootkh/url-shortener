package storage

import "errors"

// Определим общие ошибки для нашего стораджа
var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExists     = errors.New("url already exists")
	ErrALIASNotFound = errors.New("alias not found")
	ErrALIASExists   = errors.New("alias already exists")
)
