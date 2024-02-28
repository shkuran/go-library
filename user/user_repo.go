package user

import (
	"database/sql"
	"errors"

	"github.com/shkuran/go-library/utils"
)

type Repository interface {
	getById(id int64) (User, error)
	save(user User) error
	getAll() ([]User, error)
	validateCredentials(u *User) error
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) getById(id int64) (User, error) {
	var user User
	query := `
	SELECT * FROM users 
	WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repo) save(user User) error {
	query := `
	INSERT INTO users (name, email, password) 
	VALUES ($1, $2, $3)
	`

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, user.Name, user.Email, hashedPass)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) getAll() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := r.db.Query(query)
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

func (r *Repo) validateCredentials(u *User) error {
	query := `
	SELECT id, password 
	FROM users 
	WHERE email = $1
	`
	row := r.db.QueryRow(query, u.Email)

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
