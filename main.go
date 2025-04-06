package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var users = make(map[string]string)
var loggedInUser string

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/logout", logoutHandler)
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if _, exists := users[username]; exists {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		users[username] = password
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		storedPassword, exists := users[username]
		if !exists || storedPassword != password {
			// Hiển thị lại trang login với thông báo lỗi
			tmpl.Execute(w, map[string]string{
				"Error": "Sai tài khoản hoặc mật khẩu!",
			})
			return
		}

		loggedInUser = username
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}
	tmpl.Execute(w, nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	if loggedInUser == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, loggedInUser)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	loggedInUser = ""
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
