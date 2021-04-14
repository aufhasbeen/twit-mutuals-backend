package database

import (
	"log"
	"os"
	"testing"

	"github.com/aufheben/mutuals-server/local/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup() {
	// 1 touch a file in the filesystem. name it test.db.
	os.Create("test.db")

	// 2 initialize the db to sqlite to open the file.
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {

		log.Fatal(err.Error())
	}
	// make relevant tables
	db.AutoMigrate(&models.User{})
}

func teardown() {
	// 1 delete the file named test.db
	os.Remove("test.db")
}

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()

	os.Exit(result)
}

func TestSubmitUser(t *testing.T) {
	var outUser models.User

	user := models.User{}
	user.UserID = 2
	SubmitUser(&user)
	defer db.Delete(&models.User{}, 0)

	db.First(&outUser, 2)
	if outUser.UserID != user.UserID {
		t.Error("user not submitted")
	}
	// consider for addition of mutuals later on db.Model(&User{}).Association("Mutual").Append()
}

func TestFetchUserWithMutuals(t *testing.T) {
	var dummyUserA, dummyUserB, dbUser models.User

	dummyUserA.UserID = 37
	dummyUserB.UserID = 47
	dummyUserA.Mutuals = append(dummyUserA.Mutuals, dummyUserB.DBUser)
	dummyUserB.Mutuals = append(dummyUserB.Mutuals, dummyUserA.DBUser)

	SubmitUser(&dummyUserA)
	SubmitUser(&dummyUserB)

	dbUser, _ = FetchUserWithMutuals(dummyUserA.UserID)
	aPresent := dummyUserA.UserID == dbUser.UserID
	bLinkPresent := dummyUserB.UserID == dbUser.Mutuals[0].UserID

	if !(aPresent && bLinkPresent) {
		t.Error("users not linked as mutuals")
	}

	// confirmed to not recurse
	// bLinkToAPresent := !(len(dummyUserA.Mutuals[0].Mutuals) == 0)
	// if !bLinkToAPresent {
	// 	t.Error("Does not recurse")
	// }
}
