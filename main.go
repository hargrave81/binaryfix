package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	//sb := &CurrencyService{}
	//go sb.Start(30)

	amw := authenticationMiddleware{}
	amw.Populate()

	r := mux.NewRouter()
	r.Handle("/css/", ServeStatic("./html"))
	r.Handle("/scripts/", ServeStatic("./html/scripts"))
	r.Handle("/", HandleRequest(HomePageHandler))
	r.Handle("/login", HandleRequest(LoginPageHandler))
	r.Handle("/dashboard", amw.Middleware(HandleRequest(DashboardPageHandler)))
	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
