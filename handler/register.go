package rrhandler

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type HandlerRegister struct {
	lock *sync.RWMutex
	hmap map[string]*HandlerWrapper
}

func CreateHandlerRegister() (error, *HandlerRegister) {
	return nil, &HandlerRegister{
		lock: new(sync.RWMutex),
		hmap: make(map[string]*HandlerWrapper),
	}
}

func (hr *HandlerRegister) Add(key string, h Handler, t time.Duration) {
	hr.lock.Lock()
	defer hr.lock.Unlock()
	hr.hmap[key] = &HandlerWrapper{
		handle:  h,
		timeout: t,
	}
}

func (hr *HandlerRegister) Get(key string) (error, *HandlerWrapper) {
	hr.lock.RLock()
	defer hr.lock.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return errors.New(fmt.Sprintf("value for key [%d] does not exist in map", key)), nil
}
