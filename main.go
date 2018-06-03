package main

import (
    "log"
    "net/http"
		"html/template"
		"fmt"
		"database/sql"

		_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Id	int
	Title	string
	Body	string
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
    panic(err)
	}

	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
  statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, title TEXT, body TEXT)")
	 _, err :=statement.Exec()

	if err != nil {
			panic(err)
	}
}

func insertTask(x, y string) {
	db := initDB("./app.db")
	statement, _ := db.Prepare("INSERT INTO tasks(title, body) values(?,?)")
	statement.Exec(x, y)
}

func getTasks()  {
	db := initDB("./app.db")
	// tasks := []Task

	rows, err := db.Query("SELECT * FROM tasks")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &title, &body)

		task := Task {Id: id, Title: title, Body: body}
		fmt.Println(task);
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
	}

	switch r.Method {
	case "GET":
	case "POST":
		r.ParseForm()
		x := r.Form.Get("title")
		y := r.Form.Get("body")
		insertTask(x, y)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
  data := "tests"
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, data)
}

func main() {
	db := initDB("./app.db")
	migrate(db)

	getTasks()

	http.HandleFunc("/", indexHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}