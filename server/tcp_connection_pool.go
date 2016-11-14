package rrserver

import ()

type TCPConnectionPool struct {
	pool map[string]*TCPConnection
}

var (
	connPool = new(TCPConnectionPool)
)

func (s *TCPConnectionPool) Add(key string, c *TCPConnection) {
	s.pool[key] = c
}

func (s *TCPConnectionPool) Del(key string) error {
	if _, ok := s.pool[key]; !ok {
		return nil
	}
	if err := s.pool[key].conn.Close(); err != nil {
		return err
	}
	delete(s.pool, key)
	return nil
}

func (s *TCPConnectionPool) Get(key string) *TCPConnection {
	if _, ok := s.pool[key]; !ok {
		return nil
	}
	return s.pool[key]
}

func (s *TCPConnectionPool) CloseAll(key string) {
}
