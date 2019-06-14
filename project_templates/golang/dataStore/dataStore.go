package dataStore

import (
	"database/sql"

	"github.com/twinj/uuid"
)

// List the Data Store's methods
type StoreType interface {

	// userStore.go
	CreateUser(user *User) error
	GetUserByEmail(email string) ([]*User, error)
	GetPositions() ([]*Position, error)
	GetUsers() ([]*UserPublic, error)
}

func init() {
	uuid.Init()
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type DBStore struct {
	DB *sql.DB
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var Store StoreType

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s StoreType) {
	Store = s
}
