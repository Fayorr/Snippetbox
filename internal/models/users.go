package models

import (
	"database/sql"
	"time"
)

type User struct {
	id  int
	name string
	email string
	hashed_password []byte
	created  time.Time
}

type UserModel struct {
	DB *sql.DB
}
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
func (m *UserModel) Exist(id int) (bool, error) {
	return false, nil
}