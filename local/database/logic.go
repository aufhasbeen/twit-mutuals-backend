package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // database definition
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // database definition
)

var db *gorm.DB // test

type config struct {
	Host     string
	Port     string
	User     string
	Dbname   string
	Password string
}

// specific to postgres
func (c config) toString() string {

	var builder strings.Builder
	modelString := "host={{.Host}} " +
		"port={{.Port}} " +
		"user={{.User}} " +
		"dbname={{.Dbname}} " +
		"password={{.Password}}"

	tmpl, err := template.New("config").Parse(modelString)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = tmpl.Execute(&builder, c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return builder.String()
}

// specific to postgres
func configFromFile() config {
	var parsedConfig config

	configFile, err := ioutil.ReadFile(".config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	configString := string(configFile)
	decoder := json.NewDecoder(strings.NewReader(configString))
	err = decoder.Decode(&parsedConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	return parsedConfig
}

type dbType int

const (
	postgres dbType = iota
	sqlite
)

// Init initializes and connects to the database which will be used by later
// functions.
func Init(dbSystem dbType) {
	var err error

	switch dbSystem {

	case postgres:
		db, err = gorm.Open("postgres", configFromFile().toString())
		if err != nil {
			log.Fatal(err.Error())
		}

	case sqlite:
		db, err = gorm.Open("sqlite", "/file/location.db")
		if err != nil {
			log.Fatal(err.Error())
		}

	}

}

// Actual beginning of database queries

// FetchUser retrieves a user and their mutual list
func FetchUser(userID int64) {
}

// belong
type metric int

const (
	reply metric = iota
	like
	retweet
)

// FetchTop retrieves the top x number of mutuals
func FetchTop(top int, userID int64) {

}

// belong
