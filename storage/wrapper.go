package storage

// Gerneral storage wrapper
type StorageWrapper interface {
	Save([]byte) (error) // do save binary
	Fetch() ([]byte, error)
}
