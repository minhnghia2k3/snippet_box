package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID           int
	Name         string
	Email        string
	HashPassword []byte
	Created      time.Time
}

// Wrap connection pool
type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

// Add a new record to the "user" table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	sql := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(sql, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		// Check for duplicate email error.
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.unique_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate() method to verify user exists with valid credentials?
func (m *UserModel) Authenticate(email, password string) (int, error) {
	query := `SELECT id, hashed_password from users WHERE email = ?`
	u := &User{}
	// Get & Check email address from users table
	err := m.DB.QueryRow(query, email).Scan(&u.ID, &u.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Compare hashed password with the plain-text password.
	err = bcrypt.CompareHashAndPassword(u.HashPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Return user id.
	return u.ID, nil
}

// Exists method to check if user exists with specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	var exist bool

	query := `SELECT EXISTS(SELECT true FROM users WHERE id =?)`

	err := m.DB.QueryRow(query, id).Scan(&exist)

	return exist, err
}
