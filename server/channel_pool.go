package rrserver

import (
	"fmt"
	"net"
	"sync"
)

var (
	ErrClosed = fmt.Errorf("pool not created or closed")
)

type channelPool struct {
	mu    sync.Mutex
	conns chan *TCPConnection
}

// Get a new reference of conns
func (c *channelPool) getConns() chan *TCPConnection {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

func (c *channelPool) factory(addr string) (error, *TCPConnection) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err, nil
	}
	return nil, NewTCPConnection(conn)
}

func (c *channelPool) get(addr string) (error, *TCPConnection) {
	conns := c.getConns()
	if conns == nil {
		return ErrClosed, nil
	}

	select {
	case conn, ok := <-conns:
		if !ok || conn == nil {
			return ErrClosed, nil
		}
		return nil, conn
	default:
		// create a new connection
		return c.factory(addr)
	}
}

func (c *channelPool) add(conn *TCPConnection) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		// pool is closed, cann't put it into pool
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

func (c *channelPool) closePool() {
	// used when remote server down
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) lenOf() int { return len(c.getConns()) }
