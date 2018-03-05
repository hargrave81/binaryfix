package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/////////////////////////////// STATIC FILE SERVE

type FileHandler struct {
	baseFolder string
}

func (ah FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status := 200
	file := path.Join(path.Dir(ah.baseFolder), r.URL.Path)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		status = http.StatusNotFound
	} else {
		w.Write(data)
	}
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page:
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

// ServeStatic serves static files
func ServeStatic(baseDirectory string) FileHandler {
	result := FileHandler{baseDirectory}
	return result
}

// HandleRequest Creates a handler to handle web requests
func HandleRequest(RequestHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RequestHandler.ServeHTTP(w, r)
	})
}

// HomePageHandler handles homepage requests
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	err := writeFile("/index.html", w)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Served Home Page")
}

// DashboardPageHandler handles homepage requests
func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	writeFile("/dashboard.html", w)
	user, _, _ := r.BasicAuth()
	w.Write([]byte("<script>var userName = '" + user + "';</script>"))
	fmt.Println("Served Dashboard Page")
}

// LoginPageHandler handles login events
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	w.WriteHeader(http.StatusOK)
	if r.Method == "POST" {
		if vars["username"] == "" || vars["password"] == "" {
			writeFile("/login.html", w)
			w.Write([]byte("<script>var errorMessageData = 'Failed To Login';</script>"))
			fmt.Println("Failed login")
			return
		}
		if v, ok := UserTable[vars["username"]]; ok {
			// found a user
			fmt.Println("User found testing password")
			if !v.CheckPassword(vars["password"]) {
				fmt.Println("Failed to validate password")
				writeFile("/login.html", w)
				w.Write([]byte("<script>var errorMessageData = 'Failed Login - user or password incorrect';</script>"))
				return
			}
			// they are good, lets redirect to the logged in page
			securityCookie := &http.Cookie{}
			securityCookie.Expires = time.Now().Add(time.Hour * time.Duration(1))
			securityCookie.Domain = r.URL.Host
			securityCookie.Value, _ = Encrypt("SECURE-" + vars["username"] + "/" + fmt.Sprint(time.Now().Second()))
			securityCookie.Name = "session"
			http.SetCookie(w, securityCookie)
			http.Redirect(w, r, "/dashboard.html", http.StatusTemporaryRedirect)
		}
		fmt.Println("Failed to locate user")
		err = writeFile("/login.html", w)
		if err == nil {
			w.Write([]byte("<script>var errorMessageData = 'Failed Login - user or password incorrect';</script>"))
		}
	} else {
		writeFile("/login.html", w)
		w.Write([]byte("<script>var errorMessageData = '';</script>"))
	}
}

func writeFile(filename string, w http.ResponseWriter) error {
	if strings.Contains(strings.ToLower(filename), "html") {
		w.Header().Add("Content-Type", "text/html")
	}
	result, err := ioutil.ReadFile("html" + filename)
	if err != nil {
		w.WriteHeader(404)
		return err
	}
	w.Write(result)
	return nil
}
