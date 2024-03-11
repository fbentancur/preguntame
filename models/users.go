package models

import (
	"database/sql"

	// The postgres driver
	_ "github.com/lib/pq"

	"github.com/preguntame/preguntame-backend/databases"
)

type UserID = string

type User struct {
	Id       UserID
	Name     string
	Email    string
	Password string
}

func FindUserByEmail(email string) (*User, error) {
	user := User {}

	row := databases.DbPool.QueryRow("SELECT id, name, email, password FROM Users WHERE email = $1", email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		// No user found
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func InsertUser(user User) error {
	stmt := "INSERT INTO Users(id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := databases.DbPool.Exec(stmt, user.Id, user.Name, user.Email, user.Password)
	return err
}