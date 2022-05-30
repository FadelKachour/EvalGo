package router

import (
    "projetgo/middleware"

    "github.com/gorilla/mux"
)

// Le routeur est exporté et utilisé dans main.go
func Router() *mux.Router {

    router := mux.NewRouter()

    router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET")
    router.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET")
    router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST")
    router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT")
    router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE")

    return router
}