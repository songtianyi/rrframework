package rrutils

import (
	"net/http"
	_ "net/http/pprof"
	"rrframework/logs"
)

func StartProfiling() {
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logs.Error("Start profiling fail, %s", err)
			return
		}
	}()
}
