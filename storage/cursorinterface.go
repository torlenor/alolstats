package storage

type QueryCursor interface {
	Next() bool
	Decode(interface{}) error
	Close() error
}
