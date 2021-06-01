package main

import (
	"database/sql"

	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB
var err error
var db *sql.DB
var updateId string

type list struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Text   string `json:"Text"`
	Date   string `json:"Date"`
	Status string `json:"Status"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT * from lists")
	var lists []list
	for rows.Next() {
		var tem list
		rows.Scan(&tem.Id, &tem.Title, &tem.Text, &tem.Date, &tem.Status)
		lists = append(lists, tem)
	}
	data := map[string]interface{}{
		"lists": lists,
	}
	temp, _ := template.ParseFiles("webapp/index.html")
	temp.Execute(w, data)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("webapp/insert.html")
	temp.Execute(w, nil)
}

func ProcessInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	text := r.Form.Get("text")
	dobj := time.Now()
	date := dobj.Format("01-02-2006")
	DB.Exec("INSERT INTO lists(title, text, date, status) VALUES (?,?,?,?)", title, text, date, "pending")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	temp := params["id"]
	DB.Exec("DELETE FROM lists WHERE id=?", temp)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateId = params["id"]
	temp, _ := template.ParseFiles("webapp/edit.html")
	temp.Execute(w, nil)
}
func Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	text := r.Form.Get("text")
	DB.Exec("UPDATE lists SET title = ?, text =? WHERE id=?", title, text, updateId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Search(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("webapp/search.html")
	temp.Execute(w, nil)
}
func ProcessSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	search_text := r.Form.Get("searchtext")
	rows, _ := db.Query("SELECT * FROM lists WHERE INSTR(text, ?)", search_text)
	var lists []list
	for rows.Next() {
		var tem list
		rows.Scan(&tem.Id, &tem.Title, &tem.Text, &tem.Date, &tem.Status)
		lists = append(lists, tem)
	}
	data := map[string]interface{}{
		"lists": lists,
	}
	temp, _ := template.ParseFiles("webapp/results.html")
	temp.Execute(w, data)

}

func Status(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateId = params["id"]
	DB.Exec("UPDATE lists SET status ='completed' WHERE id=?", updateId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func ChangeStatusResult(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateId = params["id"]
	DB.Exec("UPDATE lists SET status ='completed' WHERE id=?", updateId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todolist")
	DB, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/todolist"))
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/insert", Insert)
	r.HandleFunc("/process", ProcessInsert)
	r.HandleFunc("/delete/{id}", Delete)
	r.HandleFunc("/edit/{id}", Edit)
	r.HandleFunc("/update", Update)
	r.HandleFunc("/search", Search)
	r.HandleFunc("/processSearch", ProcessSearch)
	r.HandleFunc("/changeStatus/{id}", Status)
	r.HandleFunc("/changeStatusResult/{id}", ChangeStatusResult)
	log.Fatal(http.ListenAndServe(":9000", r))

}
