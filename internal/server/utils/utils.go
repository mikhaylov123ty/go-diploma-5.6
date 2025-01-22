package utils

import (
	"io"
	"net/http"
)

type TransactionsHandler interface {
	Begin() error
	Commit() error
	Rollback() error
}

type ContextKey string

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type WorkersResponse struct {
	OrderID string
	Err     error
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func ArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
