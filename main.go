package main

import (
    "fmt"
    "projetgo/router"
    "log"
    "net/http"
)

func main() {
    r := router.Router()
    fmt.Println("Serveur de démarrage sur le port 8080")

    log.Print(http.ListenAndServe(":8080", r))
}