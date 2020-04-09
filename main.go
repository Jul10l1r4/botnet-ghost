package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/base64"
	"html/template"
	"os"
)

func saveResponse(x []byte) {
	fmt.Printf("Make life")
}
func queryResponse(x []byte) string {
	return "Query db for command"
}

func errorStatment(err error) {
	if err != nil {
                log.Fatal("error:", err)
        }
}
func pageAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API here")
}
func proccessData(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(os.Args[1]))
	tmpl.Execute(w, nil)
	value, err := base64.StdEncoding.DecodeString(string(r.URL.Query()["gclid"][0]))
	errorStatment(err)
	saveResponse(value)
	resp := queryResponse(value)
	fmt.Fprintf(w, "Pay\t%s", resp)
}

func main() {
	http.HandleFunc("/", proccessData)
	http.HandleFunc("/auth", pageAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
