package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RequestTime struct {
	Time string
}

type Header struct {
	Name  string
	Value string
}

func getHeaders(r *http.Request) []Header {
	var headers []Header

	for k, v := range r.Header {
		header := Header{
			Name:  k,
			Value: v[0],
		}

		headers = append(headers, header)
	}

	return headers
}

func getBodyAsString(r *http.Request) string {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	return bodyString
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	bodyString := getBodyAsString(r)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", bodyString)
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	headers := getHeaders(r)

	response, _ := json.Marshal(headers)
	fmt.Fprintf(w, "%s", response)
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr= %q\n", r.RemoteAddr)
	fmt.Fprintf(w, "Method = %q\n", r.Method)
	fmt.Fprintf(w, "Protocol = %q \n", r.Proto)

	fmt.Fprintf(w, "Headers:\n")
	headers := getHeaders(r)
	for _, v := range headers {
		fmt.Fprintf(w, "\t%q = %q\n", v.Name, v.Value)
	}

	body := getBodyAsString(r)
	fmt.Fprintf(w, "Body = %s\n", body)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format("15:04:05")

	reqTimeStruct := RequestTime{
		Time: requestTime,
	}

	tmpl := template.Must(template.ParseFiles("home.html"))
	tmpl.Execute(w, reqTimeStruct)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)
	r.HandleFunc("/all", allHandler)
	r.HandleFunc("/body", bodyHandler).Methods("POST")
	r.HandleFunc("/headers", headerHandler).Methods("GET", "POST")

	log.Println("Starting to serve at port 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
