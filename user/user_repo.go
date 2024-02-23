package user

import (
	"errors"

	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/utils"
)

func GetUserById(id int64) (User, error) {
	var user User

	row := db.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func saveUser(user User) error {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES (?, ?, ?)
	`

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(query, user.Name, user.Email, hashedPass)
	if err != nil {
		return err
	}

	return nil
}

func getUsers() ([]User, error) {
	rows, err := db.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func validateCredentials(u *User) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var passFromDB string
	err := row.Scan(&u.ID, &passFromDB)
	if err != nil {
		return errors.New("invalid credentials")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, passFromDB)

	if !passwordIsValid {
		return errors.New("invalid credentials")
	}
	return nil
}
