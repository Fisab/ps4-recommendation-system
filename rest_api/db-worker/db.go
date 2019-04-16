package database

import (
	"../config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Auth struct {
	Id int			   `db:"id"`
	Session_key string `db:"session_key"`
	Timestamp string   `db:"timestamp"`
	Uid int			   `db:"uid"`
}

func RetrieveAuthedUsers() []Auth {
	conn := Conn()

	var data []Auth
	err := conn.Select(&data, "select * from auth")
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