package rrhandler

import (
	"errors"
	"fmt"
	"time"
)

type HandlerRegister struct {
	hmap map[string]*HandlerWrapper
}

func CreateHandlerRegister() (error, *HandlerRegister) {
	return nil, &HandlerRegister{
		hmap: make(map[string]*HandlerWrapper),
	}
}

func (hr *HandlerRegister) Add(key string, h Handler, t time.Duration) {
	hr.hmap[key] = &HandlerWrapper{
		handle:  h,
		timeout: t,
	}
}

func (hr *HandlerRegister) Get(key string) (error, *HandlerWrapper) {
	if _, ok := hr.hmap[key]; !ok {
		//return errors.New(fmt.Sprintf("value for key [%d] not exist in map", key)), HandlerWrapper{handle: nil, timeout: 0 * time.Second}
		return errors.New(fmt.Sprintf("value for key [%d] not exist in map", key)), nil
	}
	return nil, hr.hmap[key]
}
