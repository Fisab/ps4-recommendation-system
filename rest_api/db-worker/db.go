package database

import (
	"../config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Auth struct {
	Id         int    `db:"id"`
	SessionKey string `db:"session_key"`
	Timestamp  string `db:"timestamp"`
	Uid        int    `db:"uid"`
}

func AuthUser(uid int, sessionKey string) bool {
	conn := Conn()

	sqlQuery := `
		INSERT INTO 
			auth 
			(session_key, timestamp, uid)
		VALUES 
			(?, ?, ?);`
	_, err := conn.Exec(sqlQuery, sessionKey, time.Now().Format(time.RFC3339), uid)
	if err != nil {
		return false
	}

	return true
}

func CheckAuthedUser(sessionKey string) bool {
	conn := Conn()

	var data []Auth

	sqlQuery := fmt.Sprintf("SELECT * FROM auth WHERE session_key = '%s'", sessionKey)
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

func Conn() *sqlx.DB {
	rawCredentials := config.GetMysqlConfig()
	credentials := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", rawCredentials.Login, rawCredentials.Password, rawCredentials.IP, rawCredentials.Port, rawCredentials.Database)

	conn, err := sqlx.Connect("mysql", credentials)
	if err != nil {
		panic(err)
	}

	return conn
}
