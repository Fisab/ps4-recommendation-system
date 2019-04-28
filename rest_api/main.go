package main

import (
	"./db-worker"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type httpResponse struct {
	StatusCode int
	StatusMsg  string
}

type httpResponseFavGenres struct {
	StatusCode int
	StatusMsg  string
	Result      []string
}

type httpResponseTopGames struct {
	StatusCode int
	StatusMsg  string
	Result      []database.Game
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

	router.HandleFunc("/auth", auth).Methods("POST")
	router.HandleFunc("/register", registerUser).Methods("POST")

	router.HandleFunc("/getFavoriteGenres", getFavoriteGenres).Methods("GET")
	router.HandleFunc("/getTopGames", getTopGames).Methods("GET")
	router.HandleFunc("/searchGames", searchGames).Methods("GET")
	router.HandleFunc("/getGameById", getGameById).Methods("GET")

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
		payload := httpResponse{http.StatusBadRequest, "U must give me ur login and sha256 of pass"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		userExist, uid := database.ValidateUser(login, password)
		if userExist == true {
			hasher := sha256.New()
			hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10) + login))
			cookie := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			addCookie(w, "session_key", cookie)

			database.CleanOldCookie(uid)
			status := database.AuthUser(uid, cookie)
			if status == true {
				payload := httpResponse{http.StatusOK, "yo"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			} else {
				payload := httpResponse{http.StatusInternalServerError, "idk what happened, seems something bad"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
		} else {
			payload := httpResponse{http.StatusForbidden, "Wrong pass or login ;("}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")
	password := r.Header.Get("password")
	email := r.Header.Get("email")

	if len(login) == 0 || len(password) == 0 {
		payload := httpResponse{http.StatusBadRequest, "U must give me ur login and sha256 of pass\n P.S. Everything is written in the docs"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		status, msg := database.RegisterUser(login, password, email)
		if status == true {
			payload := httpResponse{http.StatusAccepted, msg}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(string(js))

			fmt.Fprint(w, string(js))

		} else {
			payload := httpResponse{http.StatusConflict, msg}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}
}

func getFavoriteGenres(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_key")
	if err != nil {
		payload := httpResponse{http.StatusBadRequest, "What are junkie doing there?! Maybe u finally pass the cookies?"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		// validate cookie
		validUser, uid := database.CheckAuthedUser(c.Value)
		if validUser == true {
			// retrieve favorite genres for this cookie(get uid and then get fv genres)
			status, genres := database.GetFavGenresByUid(uid)
			if status == true {
				payload := httpResponseFavGenres{http.StatusOK, "Here is ur favorite genres", genres}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
		} else {
			payload := httpResponse{http.StatusForbidden, "Stop this... Let's just chill(wrong cookie)"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}

}

func getTopGames(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_key")
	if err != nil {
		payload := httpResponse{http.StatusBadRequest, "What are junkie doing there?! Maybe u finally pass the cookies?"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		// validate cookie
		validUser, _ := database.CheckAuthedUser(c.Value)
		if validUser == true {
			var gamesLimit []string
			gamesLimit = r.URL.Query()["games_limit"]
			gamesGenre := r.URL.Query()["games_genre"]

			if gamesLimit == nil || len(gamesLimit) == 0 {
				gamesLimit = []string{"30"}
			}
			if gamesGenre == nil || len(gamesGenre) == 0 {
				gamesGenre = []string{""}
			}

			gamesLimitInt, err := strconv.ParseInt(gamesLimit[0], 10, 64)
			if err != nil {
				payload := httpResponse{http.StatusBadRequest, "Stop this... Let's just chill(wrong games_limit value)"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
			data := database.RetrieveTopGames(gamesLimitInt, gamesGenre[0])
			payload := httpResponseTopGames{http.StatusOK, "Here is our top", data}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		} else {
			payload := httpResponse{http.StatusForbidden, "Stop this... Let's just chill(wrong cookie)"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}
}

func searchGames(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_key")
	if err != nil {
		payload := httpResponse{http.StatusBadRequest, "What are junkie doing there?! Maybe u finally pass the cookies?"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		// validate cookie
		validUser, _ := database.CheckAuthedUser(c.Value)
		if validUser == true {
			gameName := r.URL.Query()["game_name"]
			if gameName == nil || len(gameName) == 0 {
				payload := httpResponse{http.StatusBadRequest, "Stop this... Let's just chill(wrong game_name value)"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
			data := database.SearchGames(gameName[0])
			payload := httpResponseTopGames{http.StatusOK, "Welcome", data}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		} else {
			payload := httpResponse{http.StatusForbidden, "Stop this... Let's just chill(wrong cookie)"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}
}

func getGameById(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_key")
	if err != nil {
		payload := httpResponse{http.StatusBadRequest, "What are junkie doing there?! Maybe u finally pass the cookies?"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		js, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(js))
	} else {
		// validate cookie
		validUser, _ := database.CheckAuthedUser(c.Value)
		if validUser == true {
			gameId := r.URL.Query()["game_id"]
			if gameId == nil || len(gameId) == 0 {
				payload := httpResponse{http.StatusBadRequest, "Stop this... Let's just chill(wrong game_id value)"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
			gameIdInt, err := strconv.ParseInt(gameId[0], 10, 64)
			if err != nil {
				payload := httpResponse{http.StatusBadRequest, "Stop this... Let's just chill(wrong games_id value)"}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				js, err := json.Marshal(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprint(w, string(js))
			}
			data := database.GetGameById(gameIdInt)
			payload := httpResponseTopGames{http.StatusOK, "Welcome", data}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		} else {
			payload := httpResponse{http.StatusForbidden, "Stop this... Let's just chill(wrong cookie)"}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			js, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, string(js))
		}
	}
}