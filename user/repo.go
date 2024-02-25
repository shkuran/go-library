package user

import (
	"database/sql"
	"errors"
	"github.com/shkuran/go-library/utils"
)

type Repository interface {
	GetUserById(id int64) (User, error)
	SaveUser(user User) error
	GetUsers() ([]User, error)
	ValidateCredentials(u *User) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) MySQLRepository {
	return MySQLRepository{db: db}
}

func (mysql MySQLRepository) GetUserById(id int64) (User, error) {
	var user User

	row := mysql.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (mysql MySQLRepository) SaveUser(user User) error {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES (?, ?, ?)
	`

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = mysql.db.Exec(query, user.Name, user.Email, hashedPass)
	if err != nil {
		return err
	}

	return nil
}

func (mysql MySQLRepository) GetUsers() ([]User, error) {
	rows, err := mysql.db.Query("SELECT * FROM users")
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

func (mysql MySQLRepository) ValidateCredentials(u *User) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := mysql.db.QueryRow(query, u.Email)

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
