package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	LiveObjects,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	OtherSys,
	StackInuse,
	StackSys,
	Sys,

	PauseTotalNs uint64
	NumGC        uint32
	NumGoroutine int
	PollCount    int
	RandomValue  int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	var PollCount int = 1
	m.PollCount = PollCount
	rand.Seed(time.Now().Unix())
	m.RandomValue = rand.Intn(100) + 1
	for {
		<-time.After(interval)
		PollCount++
		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		m.PollCount = PollCount
		rand.Seed(time.Now().Unix())
		m.RandomValue = rand.Intn(10000) + 1

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))

	}
}

func SimpleFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Something"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func main() {

	var pollInterval int = 2
	cmd := exec.Command("ifconfig")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	fmt.Println(cmd.Stdout)
	content, err := ioutil.ReadFile("../../static/ptp-active-clients")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)

	fileServer := http.FileServer(http.Dir("./static"))
	handler2 := http.HandlerFunc(SimpleFunc)
	fmt.Println(handler2)
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port 8080\n")
	go NewMonitor(pollInterval)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
