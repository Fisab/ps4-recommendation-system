package config

import (
	"encoding/json"
	"io/ioutil"
)

type Credentials struct {
	Mysql struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		IP       string `json:"ip"`
		Port     int    `json:"port"`
	} `json:"mysql"`
}

type Mysql struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		IP       string `json:"ip"`
		Port     int    `json:"port"`
}

func GetMysqlConfig() (Mysql) {
	file, _ := ioutil.ReadFile("config/credentials.json")

	var credentials Credentials

	json.Unmarshal([]byte(file), &credentials)

	return credentials.Mysql
}
