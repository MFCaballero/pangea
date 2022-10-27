package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MFCaballero/pangea/appdomain"
	"github.com/MFCaballero/pangea/repository"
	"github.com/gorilla/mux"
)

type UsersService struct {
	DB      appdomain.UserDatastore
	JwtAuth appdomain.UserAuth
}

var usersService *UsersService

func GetUserService() *UsersService {
	if usersService == nil {
		usersService = &UsersService{DB: repository.GetUsersDataStore(), JwtAuth: GetAuthenticator()}
		return usersService
	}
	return usersService
}

func (us *UsersService) Login(w http.ResponseWriter, r *http.Request) {
	user := &appdomain.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currUser, err := us.DB.FindUser(user.Email, user.Password)

	if err != nil {
		log.Print("error occued FindUser ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tokenString, err := us.JwtAuth.GetTokenForUser(currUser)
	if err != nil {
		log.Print("error occued FindUser ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       TokenName,
		Value:      tokenString,
		Path:       "/",
		RawExpires: "0",
	})

	var resp = map[string]interface{}{"status": true, "access-token": tokenString, "user": currUser}
	json.NewEncoder(w).Encode(resp)

}

func (us *UsersService) CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &appdomain.User{}
	json.NewDecoder(r.Body).Decode(user)

	_, err := us.DB.FindUser(user.Email, user.Password)

	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := us.DB.CreateUser(user); err != nil {
		log.Print("error occued CreateUser ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	tokenString, err := us.JwtAuth.GetTokenForUser(user)
	var resp = map[string]interface{}{"status": true, "user": user, "access-token": tokenString}
	json.NewEncoder(w).Encode(resp)
	return
}

func (us *UsersService) FetchUsers(w http.ResponseWriter, r *http.Request) {

	theUsers, err := us.DB.GetAllUsers()
	if err != nil {
		log.Print("error occued during FetchUsers ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(theUsers)
}

func (us *UsersService) UpdateUser(w http.ResponseWriter, r *http.Request) {

	user := appdomain.User{}
	params := mux.Vars(r)
	var id = params["id"]

	json.NewDecoder(r.Body).Decode(&user)

	if err := us.DB.UpdateUser(id, user); err != nil {
		log.Print("error occued during user update ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func (us *UsersService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	if err := us.DB.DeleteUser(id); err != nil {
		log.Print("error occued during user delete ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("User deleted")
}

func (us *UsersService) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	user, err := us.DB.GetUser(id)

	if err != nil {
		log.Print("error occued during user delete ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&user)
}
