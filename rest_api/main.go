package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type httpResponse struct {
	status_code int
	status_msg  string
}

func dir(obj interface{}) {
	fooType := reflect.TypeOf(obj)
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		fmt.Println(method.Name)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/auth", auth).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func addCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: 3600,
	}
	http.SetCookie(w, &cookie)
}

func auth(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")
	password := r.Header.Get("password")

	if len(login) == 0 || len(password) == 0 {
		payload := httpResponse{500, "Go away"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, payload)
		return
	}

	c, err := r.Cookie("session_key")
	if err != nil {
		hasher := sha256.New()
		hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10) + login))
		cookie := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		addCookie(w, "session_key", cookie)

		// remove old cookie

		payload := httpResponse{200, "Welcome"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, payload)

		// w.Write([]byte("error in reading cookie : " + err.Error() + "\n"))

		return
	} else {
		payload := httpResponse{200, "U are already logged in"}

		// check cookie for actual

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, payload)

		// w.Write([]byte("cookie has : " + value + "\n"))

		return
	}

}
