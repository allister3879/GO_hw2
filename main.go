package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template
var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/appointmentInfo?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func init() {
	tmpl = template.Must((template.ParseFiles("registerForm.html")))
}

type appointmentInfo struct {
	Id      int
	Phone   string
	Fio     string
	Code    string
	Comment string
}

func crudHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	client := appointmentInfo{
		Phone:   r.FormValue("phone"),
		Fio:     r.FormValue("fio"),
		Code:    r.FormValue("code"),
		Comment: r.FormValue("comment"),
	}

	if r.FormValue("submit") == "Insert" {
		_, err := db.Exec("insert into appointments (phone, fio, code, comment) values (?, ?, ?, ?)",
			client.Phone, client.Fio, client.Code, client.Comment)

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: false, Message: err.Error()})
			return
		} else {

			row := db.QueryRow("SELECT * FROM appointments ORDER BY id DESC LIMIT 1")
			var insertedClient appointmentInfo
			err = row.Scan(&insertedClient.Id, &insertedClient.Phone, &insertedClient.Fio, &insertedClient.Code, &insertedClient.Comment)

			if err != nil {
				log.Println("Error retrieving last inserted record:", err)
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: false, Message: "Error retrieving last inserted record"})
				return
			}

			tmpl.Execute(w, struct {
				Success bool
				Message string
				Client  appointmentInfo
			}{Success: true, Message: "Inserted", Client: insertedClient})
		}
	}
	// else if r.FormValue("submit") == "Update" {

	// } else if r.FormValue("submit") == "Delete" {

	// }

	fmt.Println(client)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appointmentID := r.FormValue("appointment_id")
	_, err := db.Exec("DELETE FROM appointments WHERE id=?", appointmentID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", crudHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":8181", nil)
	fmt.Println("Server running on :8181")
}
