package user

type User struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
}
