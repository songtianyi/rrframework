package storage

// Gerneral storage wrapper
type StorageWrapper interface {
	Save([]byte, string) error // do save binary
	Fetch() ([]byte, error)
}
