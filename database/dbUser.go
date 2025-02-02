package database

import (
	"crypto/sha256"
	"errors"
	"time"

	"OpnLaaS.cyber.unh.edu/lib"
)

var ErrUserExists = errors.New("user already exists")

const USERS_STATEMENT = `CREATE TABLE IF NOT EXISTS users (
	email TEXT PRIMARY KEY NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	privilege INTEGER NOT NULL DEFAULT 0
);`

const INSERT_USER_STATEMENT = `INSERT INTO users (email, first_name, last_name, password_hash) VALUES (?, ?, ?, ?);`
const SELECT_USER_STATEMENT = `SELECT email, first_name, last_name, password_hash, privilege FROM users WHERE email = ?;`
const DELETE_USER_STATEMENT = `DELETE FROM users WHERE email = ?;`
const UPDATE_USER_NAME_STATEMENT = `UPDATE users SET first_name = ?, last_name = ? WHERE email = ?;`
const UPDATE_USER_PASSWORD_STATEMENT = `UPDATE users SET password_hash = ? WHERE email = ?;`
const UPDATE_USER_PRIVILEGE_STATEMENT = `UPDATE users SET privilege = ? WHERE email = ?;`

type DBUser struct {
	Email, FirstName, LastName, PasswordHash string
	CreateTime                               time.Time
	Privilege                                int
}

func UserExists(email string) bool {
	rows, err := QueuedQuery(SELECT_USER_STATEMENT, email)

	if err != nil {
		return false
	}

	defer rows.Close()
	return rows.Next()
}

func CreateUser(email, firstName, lastName, passwordHash string) error {
	if UserExists(email) {
		return ErrUserExists
	}

	return QueuedExec(INSERT_USER_STATEMENT, email, firstName, lastName, passwordHash)
}

func GetUser(email string) (*DBUser, error) {
	rows, err := QueuedQuery(SELECT_USER_STATEMENT, email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var user DBUser
	err = rows.Scan(&user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &user.Privilege)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(email string) error {
	return QueuedExec(DELETE_USER_STATEMENT, email)
}

func UpdateUserName(email, firstName, lastName string) error {
	return QueuedExec(UPDATE_USER_NAME_STATEMENT, firstName, lastName, email)
}

func UpdateUserPassword(email, passwordHash string) error {
	return QueuedExec(UPDATE_USER_PASSWORD_STATEMENT, passwordHash, email)
}

func UpdateUserPrivilege(email string, privilege int) error {
	return QueuedExec(UPDATE_USER_PRIVILEGE_STATEMENT, privilege, email)
}

func HashPassword(rawPassword string) string {
	hash := sha256.New()
	hash.Write([]byte(lib.Config.DBSalt + rawPassword))
	return string(hash.Sum(nil))
}
