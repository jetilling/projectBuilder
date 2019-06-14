package dataStore

import (
	"database/sql"
	"log"
	"time"

	"github.com/twinj/uuid"
)

type User struct {
	Id            *int       `json:"id"`
	UUID          *string    `json:"uuid"`
	Firstname     string     `json:"firstname"`
	Lastname      string     `json:"lastname"`
	Email         string     `json:"email"`
	PhoneNumber   int        `json:"phone_number"`
	Password      string     `json:"password"`
	StreetAddress string     `json:"street_address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       int        `json:"zipcode"`
	Timezone      string     `json:"timezone"`
	Created       *time.Time `json:"created"`
}

type UserPublic struct {
	Id            *int       `json:"id"`
	UUID          *string    `json:"uuid"`
	Firstname     string     `json:"firstname"`
	Lastname      string     `json:"lastname"`
	Email         string     `json:"email"`
	PhoneNumber   int        `json:"phone_number"`
	StreetAddress string     `json:"street_address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       int        `json:"zipcode"`
	Timezone      string     `json:"timezone"`
	Created       *time.Time `json:"created"`
}

type Position struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

/*
* POST
*
* CREATE USER
 */
func (store *DBStore) CreateUser(user *User) error {

	// NewV4 generates a new RFC4122 version 4 UUID a cryptographically secure random UUID.
	uuid := uuid.NewV4()

	_, err := store.DB.Query(`INSERT INTO users (uuid, first_name, last_name, email, phone_number, 
    password, street_address, city, state, zipcode, timezone) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		uuid, user.Firstname, user.Lastname,
		user.Email, user.PhoneNumber,
		user.Password, user.StreetAddress,
		user.City, user.State, user.ZipCode,
		user.Timezone)
	if err != nil {
		return err
	}

	// Fill in the users default work info - this can be edited later. I thought about making these DEFAULTS for the
	// postgres table, but I think this would be a good start to allowing other companies to set their own user defaults
	_, err = store.DB.Query(`INSERT INTO user_default_work_info (user_uuid, user_default_weight, user_default_rental_fee,
		user_default_in_training, user_default_revenue_bonus) VALUES ($1,$2,$3,$4,$5)`,
		uuid, 80, 10, false, false)
	return err
}

/*
* GET
*
* GET USER BY EMAIL
 */
func (store *DBStore) GetUserByEmail(email string) ([]*User, error) {
	rows, err := store.DB.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	users, err := formatUserQuery_Helper(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

/*
* GET
*
* GET USERS
 */
func (store *DBStore) GetUsers() ([]*UserPublic, error) {
	rows, err := store.DB.Query(`SELECT id, uuid, first_name, last_name, email, phone_number,
															street_address, city, state, zipcode, timezone, created FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*UserPublic{}

	for rows.Next() {
		user := &UserPublic{}

		err := rows.Scan(&user.Id, &user.UUID, &user.Firstname, &user.Lastname,
			&user.Email, &user.PhoneNumber, &user.StreetAddress,
			&user.City, &user.State, &user.ZipCode, &user.Timezone, &user.Created)
		if err != nil {
			panic(err)
			// return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

/*
* GET
*
* GET POSITIONS
 */
func (store *DBStore) GetPositions() ([]*Position, error) {
	// Query the database for all users, and return the result to the
	// `rows` object
	rows, err := store.DB.Query("SELECT * from positions")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of users
	positions := []*Position{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a user,
		position := &Position{}
		// Populate the `Firstname` and `Lastname` attributes of the user,
		// and return incase of an error
		if err := rows.Scan(&position.Id, &position.Name); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		positions = append(positions, position)
	}
	return positions, nil
}

/* ----------------- Helper Functions ------------------ */

func formatUserQuery_Helper(rows *sql.Rows) ([]*User, error) {
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}

		err := rows.Scan(&user.Id, &user.UUID, &user.Firstname, &user.Lastname,
			&user.Email, &user.Password, &user.PhoneNumber, &user.StreetAddress,
			&user.City, &user.State, &user.ZipCode, &user.Timezone, &user.Created)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
