package rrhandler

import (
	"time"
	"errors"
	"fmt"
)

type HandleRegister struct {
	hmap map[string]Handler
}

func CreateHandleRegister() (error, *HandleRegister) {
	return nil, &HandleRegister{
		hmap: make(map[string]Handler),
	}
}

func (hr *HandleRegister) Add(key string, w HandlerFuncWrapper, timeout time.Duration) {
	hr.hmap[key] = Handler{
		wrapper: w,
		timeout: timeout,
	}
}

func (hr *HandleRegister) Get(key string) (error, interface{}) {
	if _, ok := hr.hmap[key]; !ok {
		return errors.New(fmt.Sprintf("value for key [%d] not exist in hmap", key)), nil
	}
	return nil, hr.hmap[key]
}
