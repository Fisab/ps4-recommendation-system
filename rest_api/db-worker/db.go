package database

import (
	"../config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type Auth struct {
	Id         int    `db:"id"`
	SessionKey string `db:"session_key"`
	Timestamp  string `db:"timestamp"`
	Uid        int    `db:"uid"`
}

type User struct {
	Timestamp      string `db:"timestamp_creation"`
	Uid            int    `db:"uid"`
	Login          string `db:"login"`
	Password       string `db:"password"`
	Mail           string `db:"mail"`
	Wishlist       string `db:"wishlist"`
	FavoriteGenres string `db:"favorite_genres"`
}

type Game struct {
	GameId        int     `db:"game_id"`
	Genres        string  `db:"genres"`
	Rating        string  `db:"rating"`
	Developer     string  `db:"developer"`
	OfPlayers     int     `db:"ofplayers"`
	Name          string  `db:"name"`
	ImgLink       string  `db:"img_link"`
	Summary       string  `db:"summary"`
	Metascore     int     `db:"metascore"`
	UsersScore    float32 `db:"users_score"`
	ProcessedName string  `db:"processed_name"`
}

func AuthUser(uid int, sessionKey string) bool {
	conn := Conn()

	sqlQuery := `
		INSERT INTO 
			auth 
			(session_key, timestamp, uid)
		VALUES 
			(?, ?, ?);`
	_, err := conn.Exec(sqlQuery, sessionKey, time.Now().Format("2006-01-02T15:04:05"), uid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func CheckAuthedUser(sessionKey string) (bool, int) {
	conn := Conn()

	var data []Auth

	sqlQuery := fmt.Sprintf("SELECT * FROM auth WHERE session_key = '%s'", sessionKey)
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return false, 0
	} else {
		return true, data[0].Uid
	}
}

func GetFavGenresByUid(uid int) (bool, []string) {
	conn := Conn()

	var data []User

	sqlQuery := fmt.Sprintf("SELECT * FROM users WHERE uid = %d", uid)
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return false, nil
	} else {
		if len(data[0].FavoriteGenres) == 0 {
			return true, nil
		} else {
			return true, strings.Split(data[0].FavoriteGenres, ",")
		}
	}
}

func RegisterUser(login string, pass string, email string) (bool, string) {
	conn := Conn()

	userExist := checkUserExist(login)

	if userExist == true {
		return false, "U want keep alien login?"
	}

	sqlQuery := fmt.Sprintf("INSERT INTO users (timestamp_creation, login, password, mail, wishlist, favorite_genres) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');", time.Now().Format("2006-01-02T15:04:05"), login, pass, email, "", "")
	_, err := conn.Exec(sqlQuery)
	if err != nil {
		fmt.Println(err)
		return false, "Server is tired, that's what it said :" + err.Error()
	}

	return true, "Now u're in a gang"
}

func checkUserExist(login string) bool {
	conn := Conn()
	var data []User

	sqlQuery := fmt.Sprintf("SELECT * FROM users WHERE login = '%s'", login)
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}

func RetrieveAuthedUsers() []Auth {
	conn := Conn()

	var data []Auth
	err := conn.Select(&data, "SELECT * FROM auth")
	if err != nil {
		panic(err)
	}

	return data
}

func ValidateUser(login string, password string) (bool, int) {
	conn := Conn()

	var data []User

	sqlQuery := fmt.Sprintf("SELECT * FROM users WHERE login = '%s' AND password = '%s'", login, password)
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return false, 0
	} else {
		return true, data[0].Uid
	}
}

func CleanOldCookie(uid int) {
	conn := Conn()

	sqlQuery := `
		DELETE FROM 
			auth 
		WHERE
		 	uid = ?
	`

	_, err := conn.Exec(sqlQuery, uid)
	if err != nil {
		panic(err)
	}
}

func Conn() *sqlx.DB {
	rawCredentials := config.GetMysqlConfig()
	credentials := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", rawCredentials.Login, rawCredentials.Password, rawCredentials.IP, rawCredentials.Port, rawCredentials.Database)

	conn, err := sqlx.Connect("mysql", credentials)
	if err != nil {
		panic(err)
	}

	return conn
}

func retrieveUidByCookie(cookie string) int {
	conn := Conn()

	var data []Auth

	sqlQuery := fmt.Sprintf("SELECT * FROM auth WHERE session_key = '%s'", cookie)
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return -1
	} else {
		return data[0].Uid
	}
}

func RetrieveTopGames(limitGames int64, gamesGenre string) []Game {
	var sqlQuery string
	if gamesGenre == "" {
		sqlQuery = `
			SELECT
				*,
				REPLACE(name, '-', ' ') as processed_name
			FROM
				games
			ORDER BY
				metascore
			LIMIT %d;
		`
		sqlQuery = fmt.Sprintf(sqlQuery, limitGames)
	} else {
		sqlQuery = `
			SELECT
				*,
				REPLACE(name, '-', ' ') as processed_name
			FROM
				games
			WHERE
				genres LIKE '%%%s%%'
			ORDER BY
				metascore
			LIMIT %d;
		`
		sqlQuery = fmt.Sprintf(sqlQuery, gamesGenre, limitGames)
	}

	conn := Conn()

	var data []Game
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	return data
}

func SearchGames(gameName string) []Game {
	sqlQuery := `
		SELECT
			*,
			REPLACE(name, '-', ' ') as processed_name
		FROM
			games
		WHERE
			name LIKE '%%%s%%'
		ORDER BY
			metascore
		LIMIT 20
	`
	sqlQuery = fmt.Sprintf(sqlQuery, gameName)

	fmt.Println(sqlQuery)

	conn := Conn()

	var data []Game
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	return data
}

func GetGameById(gameId int64) []Game {
	sqlQuery := `
		SELECT
			*,
			REPLACE(name, '-', ' ') as processed_name
		FROM
			games
		WHERE
			game_id = %d
	`
	sqlQuery = fmt.Sprintf(sqlQuery, gameId)

	conn := Conn()

	var data []Game
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	return data
}

func GetGenres() []string {
	sqlQuery := `
		SELECT DISTINCT 
			SUBSTRING_INDEX(genres, ',', 1) as genres
		FROM
			games;
	`
	conn := Conn()

	var data []string
	err := conn.Select(&data, sqlQuery)
	if err != nil {
		panic(err)
	}
	return data
}