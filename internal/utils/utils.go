package utils

import (
	"encoding/base64"
	"io"
	"net/http"
)

type ContextKey string

// Структура обертки компрессии gzip для интерфейса writer
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type WorkersResponse struct {
	WorkerID int
	Err      error
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

func EncodeString(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func DecodeString(s string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil

}
