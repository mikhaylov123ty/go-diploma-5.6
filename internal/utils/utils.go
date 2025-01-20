package utils

import (
	"io"
	"net/http"
)

type ContextKey string

// Структура обертки компрессии gzip для интерфейса writer
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Обертка метода Write для записи компрессированных сообщений
func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Метод проверки строки в массиве строк
func ArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
