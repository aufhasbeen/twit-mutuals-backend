package models

// User contains the highest model of user for use withing the whole application
type User struct {
	DBUser `json:"-"`
	jsonUser
}
