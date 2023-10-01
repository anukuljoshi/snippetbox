package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, nil
	}
	query := `
		INSERT INTO users (name, email, hashed_password, created)
		VALUES (?, ?, ?, UTC_TIMESTAMP())
	`
	result, err := u.DB.Exec(query, name, email, hashedPassword)
	if err!=nil {
		// check if error is due to duplicate email
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number==1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}
	id, err := result.LastInsertId()
	if err!=nil {
		return 0, nil
	}
	return int(id), nil
}

// authenticate a user with email and password
func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// check if user with id exists
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
