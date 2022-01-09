package main

import (
	//"fmt"
	//"os"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	//"bufio"
	"net/http"
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
	metod := r.Method

	log.Println(q)

	subpath := strings.Split(q, "/")
	//fmt.Println(subpath)
	fmt.Println(metod)
	if metod == "POST" {
		var m1 metricValue

		switch subpath[2] {
		case "gauge":
			m1.val = subpath[4]
			m1.isCounter = false
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		case "counter":

			v, err := strconv.ParseInt(subpath[4], 10, 64)
			check(err)
			mI = mI + v
			m1.val = strconv.FormatInt(mI, 10)
			m1.isCounter = true
			w.WriteHeader(http.StatusOK)
			r.Body.Close()

		default:
			fmt.Println("Type", subpath[2], "wrong")
			outputMessage := "Type " + subpath[2] + " not supported, only [counter/gauge]"
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(outputMessage))
			r.Body.Close()
		}
		metricMap[subpath[3]] = m1
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
	//metricGaugeMap := make(map[string]float64)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:8080", nil)
	check(err)
}
