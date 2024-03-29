package appdomain

import "net/http"

const (
	Admin = iota + 1
	Priviledge
	Regular
)

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type UserDatastore interface {
	CreateUser(user *User) error
	GetAllUsers() ([]User, error)
	FindUser(email, password string) (*User, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
	GetUser(id string) (User, error)
}

type UserAuth interface {
	IsTokenExists(r *http.Request) (bool, string)
	IsUserTokenValid(token string) bool
	UserFromToken(tokenString string) (*User, error)
	GetTokenForUser(user *User) (string, error)
}
