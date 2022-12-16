package api

import (
	"github.com/MFCaballero/pangea/service"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	us := service.GetUserService()
	jv := service.GetAuthenticator()
	ses := service.NewSessionManager()
	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")
	r.Use(ses.LoadAndSave)
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(jv.JwtVerify)
	s.HandleFunc("/user", us.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", us.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", us.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", us.DeleteUser).Methods("DELETE")
	return r
}
