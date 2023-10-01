package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}

// insert a new user to db
func (u *UserModel) Insert(name, email, password string) (int, error) {
	return 0, nil
}

// authenticate a user with email and password
func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// check if user with id exists
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
