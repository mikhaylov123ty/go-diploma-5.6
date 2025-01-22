package transactions

type Handler interface {
	Begin() error
	Commit() error
	Rollback() error
}
