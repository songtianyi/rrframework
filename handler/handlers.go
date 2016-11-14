package rrhandler

import (
	"time"
)

// for request & response

type HandlerFuncWrapper func(interface{}, interface{})

type Handler struct { 
	wrapper HandlerFuncWrapper
	timeout time.Duration
}

func (h Handler) Run(c interface{}, msg interface{}) {
	h.wrapper(c, msg)
}

// for timer
