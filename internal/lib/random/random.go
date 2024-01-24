package random

import (
	"math/rand"
	"time"
)

// NewRandomString - генерирует строку случайных симоволов заданной длины.
// На входе	: принимает integer - размер генерируемой строки.
// На выходе	: возвращает строку случайных симоволов заданной длины.
func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Создаём слайс из символов используемых при генерации
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	// Создаём слайс-буффер для символов
	b := make([]rune, size)

	// Проходим во слайсу-буфферу и заполняем его символами из слайса символов
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
