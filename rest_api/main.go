package main

import (
	"./db-worker"
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

type httpResponseFavGenres struct {
	status_code int
	status_msg  string
	result      []string
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
	router.HandleFunc("/getFavoriteGenres", getFavoriteGenres).Methods("GET")

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
		payload := httpResponse{500, "U must give me ur login and sha256 of pass"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, payload)
		return
	} else {
		userExist, uid := database.ValidateUser(login, password)
		if userExist == true {
			hasher := sha256.New()
			hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10) + login))
			cookie := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			addCookie(w, "session_key", cookie)

			database.CleanOldCookie(uid)

			payload := httpResponse{200, "Welcome"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, payload)
		} else {
			payload := httpResponse{403, "Wrong pass or login ;("}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, payload)
			return
		}
	}
}

func getFavoriteGenres(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_key")
	if err != nil {
		payload := httpResponse{403, "What are junkie doing there?! Maybe u finally pass the cookies?"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, payload)
		return
	} else {
		// validate cookie
		validUser, uid := database.CheckAuthedUser(c.Value)
		if validUser == true {
			// retrieve favorite genres for this cookie(get uid and then get fv genres)
			status, genres := database.GetFavGenresByUid(uid)
			if status == true {
				payload := httpResponseFavGenres{200, "Welcome", genres}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, payload)
			}
		} else {
			payload := httpResponse{403, "Stop this... Let's just chill(wrong cookie)"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, payload)
			return
		}
	}

}
