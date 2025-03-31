package database

type DBErrorProcessor interface {
	ProcessError(err error) error
}
