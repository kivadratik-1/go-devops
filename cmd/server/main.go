package main

import (
	//"fmt"
	//"os"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	//"bufio"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type metricValue struct {
	val       string
	isCounter bool
}

var (
	metricMap = make(map[string]metricValue)
)

var mI int64

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
			m1.val = chi.URLParam(r, "metricValue")
			m1.isCounter = false
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		case "counter":

			v, err := strconv.ParseInt(chi.URLParam(r, "metricValue"), 10, 64)
			check(err)
			mI = mI + v
			m1.val = strconv.FormatInt(mI, 10)
			m1.isCounter = true
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		default:
			fmt.Println("Type", chi.URLParam(r, "metricType"), "wrong")
			outputMessage := "Type " + chi.URLParam(r, "metricType") + " not supported, only [counter/gauge]"
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(outputMessage))
			r.Body.Close()
		}
		metricMap[chi.URLParam(r, "metricName")] = m1
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
	// http.HandleFunc("/", handler)
	// err := http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Serving on " + port)
	err := http.ListenAndServe("127.0.0.1"+port, r)
	check(err)
}
