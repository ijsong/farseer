package storage

type Storage interface {
	Connect() (StorageSession, error)
}

type StorageSession interface {
	Close()
}
