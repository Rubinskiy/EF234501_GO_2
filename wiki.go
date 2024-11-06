package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	Title string
	Body  string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password are required.", http.StatusBadRequest)
		return
	}
	userData := username + ";" + password + "\n"
	err := ioutil.WriteFile("users.txt", []byte(userData), 0644)
	if err != nil {
		http.Error(w, "Failed to save user data.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User saved successfully."))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password are required.", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadFile("users.txt")
	if err != nil {
		http.Error(w, "Failed to read user data.", http.StatusInternalServerError)
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ";")
		if fields[0] == username && fields[1] == password {
			w.Write([]byte("Valid user."))
			return
		}
	}
	http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
}

var templates = template.Must(template.ParseFiles("login.html", "register.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/validate", validateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}