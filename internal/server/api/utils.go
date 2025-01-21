package api

type transactionsHandler interface {
	Begin() error
	Commit() error
	Rollback() error
}
