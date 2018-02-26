package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleRequest Creates a handler to handle web requests
func HandleRequest(RequestHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RequestHandler.ServeHTTP(w, r)
	})
}

// HomePageHandler handles homepage requests
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html><body>Welcome</body></html>")
	fmt.Println("Served Home Page")
}

// LoginPageHandler handles login events
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	if r.Method == "POST" {
		if vars["username"] == "" || vars["password"] == "" {
			PrintLogin("You must supply a username and valid password")
			fmt.Println("Failed login")
		}
		if v, ok := UserTable[vars["username"]]; ok {
			// found a user
			fmt.Println("User found testing password")
		}
	}
	{

		fmt.Fprint(w, "<html><body>Welcome</body></html>")
		fmt.Println("Served Home Page")
	}
}
