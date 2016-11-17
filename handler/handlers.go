package rrhandler

import (
	"time"
)

// handler for requests
type Handler func(interface{}, chan []byte)

type HandlerWrapper struct {
	handle  Handler
	timeout time.Duration
}

func (h HandlerWrapper) Run(req interface{}, res chan []byte) {
	h.handle(req, res)
}

// handler for timer jobs
