package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must((template.ParseFiles("registerForm.html")))
}

type appointmentInfo struct {
	Phone   string
	Fio     string
	Code    string
	Comment string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
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
	tmpl.Execute(w, struct {
		Succes bool
		Client appointmentInfo
	}{true, client})

}

func main() {
	http.HandleFunc("/", formHandler)
	http.ListenAndServe(":8181", nil)
}
