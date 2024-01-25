package random

import (
	"math/rand"
	"time"
)

// --------------------------------------------------------------------------------------

// NewRandomString - генерирует строку случайных символов заданной длины.
//
// Принимает:	size (integer) - размер генерируемой строки.
//
// Возвращает	: string - строку случайных символов заданной длины.
func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Создаём слайс из символов используемых при генерации
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	// Создаём слайс-буфер для символов
	b := make([]rune, size)

	// Проходим по слайс-буферу и заполняем его символами из слайса символов
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
