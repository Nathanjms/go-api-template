package database

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	*sqlx.DB
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (userModel *UserModel) FindUser(id int64) (User, error) {
	u := new(User)

	row := userModel.DB.QueryRow("SELECT id, username, password FROM users WHERE id = ?", id)

	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return User{}, err
	}

	return *u, nil
}

func (userModel *UserModel) GetByUsername(username string) (User, error) {
	u := new(User)

	row :=
		userModel.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)

	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return User{}, err
	}

	return *u, nil
}

func (u *UserModel) Create(username string, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := u.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *UserModel) Delete(id int64) error {
	_, err := u.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
