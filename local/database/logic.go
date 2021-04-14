package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/aufheben/mutuals-server/local/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	postgresDriver dbType = iota
	sqliteDriver
)

// Init initializes and connects to the database which will be used by later
// functions.
func Init(dbSystem dbType) {
	var err error

	switch dbSystem {

	case postgresDriver:
		db, err = gorm.Open(postgres.Open("users.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err.Error())
		}

	case sqliteDriver:
		db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err.Error())
		}

	}

	primeDb()

}

func primeDb() {
	db.AutoMigrate(&models.User{}) // , &Mutual{})
}

// Actual beginning of database queries

// SubmitUser uploads the user to the database ensuring a unique id
func SubmitUser(user *models.User) {
	db.Save(user)
}

// FetchUser retrieves a user without its mutual list
func FetchUser(screenName string) models.User {
	var user models.User
	var queryUser models.User
	queryUser.ScreenName = screenName
	db.Where(&queryUser).First(&user)
	return user
}

// FetchUserWithMutuals fetches the user and all it's Mutuals
func FetchUserWithMutuals(userID int64) (models.User, error) {
	var user models.User
	user.UserID = userID
	err := db.Preload("Mutuals").First(&user).Error
	return user, err

}

// begin
type metric int

const (
	reply metric = iota
	like
	retweet
)

// FetchTop retrieves the top x number of mutuals
func FetchTop(top int, userID int64) {

}

// end

// FetchMutuals returns a list from the database of the mutuals of the User
// with userID as its key.
// func FetchMutuals(userID int64) ([]Mutual, error) {
// 	// query the db for all the mutuals pertaining to user with userID
// 	return make([]Mutual, 0), nil
// }
