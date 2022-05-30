# Golang-api

Utilisation de Gorilla/Mux pour le routage.
Utilisation de ElephantSQL pour un hébergement.
Application de coureur/footing

## STARTUP

Lancer le projet
go run main.go

## architecture
```
├── ProjetGo
│    ├── .env
│    ├── go.mod
│    ├── go.sum
│    ├── main.go
│    ├── middleware
│    │   └──  handlers.go
│    ├── models
│    │   └── Models.go
│    └── router
│        └── Router.go
```

### Routes

User :
- GET /users
- GET /users/:id
- POST /users {name, niveau, kmmax}
- PUT /users/:id {name, niveau, kmmax}
- DELETE /users/:id

## Modele

### User :
- Id
- Name
- Niveau
- KmMax