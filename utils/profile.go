package rrutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"rrframework/logs"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"time"
)

var startTime = time.Now()
var pid int

func init() {
	pid = os.Getpid()
}

// ProcessInput parse input command string
func ProcessInput(input string, w io.Writer) {
	switch input {
	case "goroutine":
		p := pprof.Lookup("goroutine")
		p.WriteTo(w, 2)
	case "heap":
		p := pprof.Lookup("heap")
		p.WriteTo(w, 2)
	case "threadcreate":
		p := pprof.Lookup("threadcreate")
		p.WriteTo(w, 2)
	case "block":
		p := pprof.Lookup("block")
		p.WriteTo(w, 2)
	case "cpuprof":
		GetCPUProfile(w)
	case "memprof":
		MemProf(w)
	case "gc":
		PrintGCSummary(w)
	}
}

// MemProf record memory profile in pprof
func MemProf(w io.Writer) {
	filename := "mem-" + strconv.Itoa(pid) + ".memprof"
	if f, err := os.Create(filename); err != nil {
		fmt.Fprintf(w, "create file %s error %s\n", filename, err.Error())
		logs.Error("record heap profile failed: ", err)
	} else {
		runtime.GC()
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Fprintf(w, "create heap profile %s \n", filename)
		_, fl := path.Split(os.Args[0])
		fmt.Fprintf(w, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)
	}
}

// GetCPUProfile start cpu profile monitor
func GetCPUProfile(w io.Writer) {
	sec := 30
	filename := "cpu-" + strconv.Itoa(pid) + ".pprof"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(w, "Could not enable CPU profiling: %s\n", err)
		logs.Error("record cpu profile failed: ", err)
	}
	pprof.StartCPUProfile(f)
	time.Sleep(time.Duration(sec) * time.Second)
	pprof.StopCPUProfile()

	fmt.Fprintf(w, "create cpu profile %s \n", filename)
	_, fl := path.Split(os.Args[0])
	fmt.Fprintf(w, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)
}

// PrintGCSummary print gc information to io.Writer
func PrintGCSummary(w io.Writer) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	gcstats := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gcstats)

	printGC(memStats, gcstats, w)
}

func printGC(memStats *runtime.MemStats, gcstats *debug.GCStats, w io.Writer) {

	if gcstats.NumGC > 0 {
		lastPause := gcstats.Pause[0]
		elapsed := time.Now().Sub(startTime)
		overhead := float64(gcstats.PauseTotal) / float64(elapsed) * 100
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "NumGC:%d Pause:%s Pause(Avg):%s Overhead:%3.2f%% Alloc:%s Sys:%s Alloc(Rate):%s/s Histogram:%s %s %s \n",
			gcstats.NumGC,
			toS(lastPause),
			toS(avg(gcstats.Pause)),
			overhead,
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)),
			toS(gcstats.PauseQuantiles[94]),
			toS(gcstats.PauseQuantiles[98]),
			toS(gcstats.PauseQuantiles[99]))
	} else {
		// while GC has disabled
		elapsed := time.Now().Sub(startTime)
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "Alloc:%s Sys:%s Alloc(Rate):%s/s\n",
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)))
	}
}

func avg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

// format bytes number friendly
func toH(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}

// short string format
func toS(d time.Duration) string {

	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			return "0"
		case u < uint64(time.Microsecond):
			return fmt.Sprintf("%.2fns", float64(u))
		case u < uint64(time.Millisecond):
			return fmt.Sprintf("%.2fus", float64(u)/1000)
		default:
			return fmt.Sprintf("%.2fms", float64(u)/1000/1000)
		}
	} else {
		switch {
		case u < uint64(time.Minute):
			return fmt.Sprintf("%.2fs", float64(u)/1000/1000/1000)
		case u < uint64(time.Hour):
			return fmt.Sprintf("%.2fm", float64(u)/1000/1000/1000/60)
		default:
			return fmt.Sprintf("%.2fh", float64(u)/1000/1000/1000/60/60)
		}
	}
}

func ProfIndex(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")

	var (
		data   = make(map[string]string)
		result bytes.Buffer
	)
	if command == "" {
		data["Content"] = "you can use command=goroutine,heap,threadcreate,block,cpuprof,memprof,gc"
	} else {
		ProcessInput(command, &result)
		data["Content"] = result.String()
	}

	rw.Header().Set("Content-Type", "application/text")
	rw.Write([]byte(data["Content"]))
	return
}

// for profile file download
func isExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func FileHandler(mux *http.ServeMux, prefix string, staticDir string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if exists, _ := isExists(file); !exists {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, file)
	})
}

func StartProfiling() {
	go func() {
		mux := http.NewServeMux()
		FileHandler(mux, "/file/", ".")
		mux.HandleFunc("/profile", ProfIndex)
		if err := http.ListenAndServe("localhost:6060", mux); err != nil {
			logs.Error("ListenAndServe failed, %s", err)
		}
	}()
}
