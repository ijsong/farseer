package storage

type Storage interface {
	Connect() (StorageSession, error)
	Close() error
}

type StorageSession interface {
	Close()
}
