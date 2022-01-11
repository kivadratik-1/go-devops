package main

import (
	//"fmt"
	//"os"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	//"bufio"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type metricValue struct {
	val       [8]byte
	isCounter bool
}

var (
	metricMap = make(map[string]metricValue)
	mI        int64
)

func int64ToBytes(value int64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(value))
	return buf
}

func int64FromBytes(bytes []byte) int64 {
	i := binary.LittleEndian.Uint64(bytes)
	return int64(i)
}

func float64ToBytes(value float64) [8]byte {
	bits := math.Float64bits(value)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], bits)
	return bytes
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.RequestURI()
	log.Println(q)
	reqMethod := r.Method
	log.Println(reqMethod)

	//subpath := strings.Split(q, "/")
	//fmt.Println(subpath)
	//fmt.Println(reqMethod)
	if reqMethod == "POST" {
		var m1 metricValue
		switch chi.URLParam(r, "metricType") {
		case "gauge":
			f, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64)
			check(err)
			m1.val = float64ToBytes(f)
			fmt.Println(float64FromBytes(m1.val[:]))
			m1.isCounter = false
			metricMap[chi.URLParam(r, "metricName")] = m1
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		case "counter":
			i, err := strconv.ParseInt(chi.URLParam(r, "metricValue"), 10, 64)
			check(err)
			mI = mI + i
			m1.val = int64ToBytes(mI)
			m1.isCounter = true
			fmt.Println(int64FromBytes(m1.val[:]))
			metricMap[chi.URLParam(r, "metricName")] = m1
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		default:
			fmt.Println("Type", chi.URLParam(r, "metricType"), "wrong")
			outputMessage := "Type " + chi.URLParam(r, "metricType") + " not supported, only [counter/gauge]"
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(outputMessage))
			r.Body.Close()
		}

		//fmt.Println(metricMap)

		options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
		file, err := os.OpenFile("metrics.data", options, os.FileMode(0600))
		check(err)
		_, err = fmt.Fprintln(file, metricMap)
		check(err)
		err = file.Close()
		check(err)
	} else {
		fmt.Println("Method is wrong")
		outputMessage := "Only POST method is alload"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))

	}
}

func handlerGetMetrics(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	q := r.URL.RequestURI()
	log.Println(q)
	reqMethod := r.Method
	log.Println(reqMethod)

	if reqMethod == "GET" {
		switch chi.URLParam(r, "metricType") {
		case "gauge":
			v, ok := metricMap[metricName]

			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%v\n", float64FromBytes(v.val[:]))

			r.Body.Close()

		case "counter":
			v, ok := metricMap[metricName]

			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%v\n", int64FromBytes(v.val[:]))

			r.Body.Close()

		default:
			fmt.Println("Type", chi.URLParam(r, "metricType"), "wrong")
			outputMessage := "Type " + chi.URLParam(r, "metricType") + " not supported, only [counter/gauge]"
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(outputMessage))
			r.Body.Close()
		}

		fmt.Println(metricMap)
		options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
		file, err := os.OpenFile("metrics.data", options, os.FileMode(0600))
		check(err)
		_, err = fmt.Fprintln(file, metricMap)
		check(err)
		err = file.Close()
		check(err)
	} else {
		fmt.Println("Method is wrong")
		outputMessage := "Only POST method is alload"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))

	}
}

func main() {
	port := ":8080"
	r := chi.NewRouter()
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler)
	r.Get("/value/{metricType}/{metricName}", handlerGetMetrics)
	r.Get("/", httpPrintAllMetrics)

	// http.HandleFunc("/", handler)
	// err := http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Serving on " + port)
	err := http.ListenAndServe("127.0.0.1"+port, r)
	check(err)
}
