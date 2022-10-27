package repository

import (
	"github.com/MFCaballero/pangea/appdomain"
	"github.com/jmoiron/sqlx"
)

type UsersDB struct {
	*sqlx.DB
}

func GetUsersDataStore() appdomain.UserDatastore {
	return &UsersDB{DBGet()}
}

func (db *UsersDB) CreateUser(user *appdomain.User) error {
	return nil
}

func (db *UsersDB) GetAllUsers() ([]appdomain.User, error) {
	var users []appdomain.User
	return users, nil
}

func (db *UsersDB) FindUser(email, password string) (*appdomain.User, error) {
	user := &appdomain.User{}
	return user, nil
}

func (db *UsersDB) UpdateUser(id string, user appdomain.User) error {
	return nil
}

func (db *UsersDB) DeleteUser(id string) error {
	return nil
}

func (db *UsersDB) GetUser(id string) (appdomain.User, error) {
	user := appdomain.User{}
	return user, nil
}
