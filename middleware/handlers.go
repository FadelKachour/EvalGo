package middleware

import (
	"database/sql"
	"encoding/json" // package pour encoder et décoder le json en struct et vice versa
	"fmt"
	"log"
	"net/http"        // utilisé pour accéder à l'objet de requête et de réponse de l'API
	"os"              // utilisé pour lire la variable d'environnement
	"projetgo/models" // package de modèles où le schéma utilisateur est défini
	"strconv"         // package utilisé pour convertir une chaîne en type int

	"github.com/gorilla/mux" // utilisé pour obtenir les paramètres de la route

	"github.com/joho/godotenv" // package utilisé pour lire le fichier .env
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// créer une connexion avec la base de données postgres
func createConnection() *sql.DB {

	// charger le fichier .env
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env")
	}

	// Ouvrir la connexion
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	// continuer les go routines
	if err != nil {
		panic(err)
	}

	// vérifier la connexion
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connecté avec succès !")
	// retourner la connexion
	return db
}

// crée un utilisateur dans la base de données postgres
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// créer un utilisateur vide de type models.User
	var user models.User

	// décoder la requête json à l'utilisateur
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Printf("Impossible de décoder le corps de la requête.  %v", err)
	}

	// appelez la fonction utilisateur d'insertion et passez l'utilisateur
	insertID := insertUser(user)

	// formater un objet de réponse
	res := response{
		ID:      insertID,
		Message: "Utilisateur créé avec succès",
	}

	// envoyer la réponse
	json.NewEncoder(w).Encode(res)
}

// renverra un seul utilisateur par son identifiant
func GetUser(w http.ResponseWriter, r *http.Request) {

	// obtenir l'ID utilisateur à partir des paramètres de la requête, la clé est "id"
	params := mux.Vars(r)

	// convertir le type d'identifiant de chaîne en int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Impossible de convertir la chaîne en int.  %v", err)
	}

	// appelez la fonction getUser avec l'identifiant de l'utilisateur pour récupérer un seul utilisateur
	user, err := getUser(int64(id))

	if err != nil {
		log.Printf("Impossible d'obtenir l'utilisateur. %v", err)
	}

	// envoyer la réponse
	json.NewEncoder(w).Encode(user)

}

// renverra tous les utilisateurs
func GetAllUser(w http.ResponseWriter, r *http.Request) {

	// obtenir tous les utilisateurs de la base de données
	users, err := getAllUsers()

	if err != nil {
		log.Printf("Impossible d'obtenir tous les utilisateurs. %v", err)
	}

	// envoyer tous les utilisateurs en réponse
	json.NewEncoder(w).Encode(users)
}

// mettre à jour les détails de l'utilisateur dans la base de données postgres
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// obtenir l'ID utilisateur à partir des paramètres de la requête, la clé est "id"
	params := mux.Vars(r)

	// convertir le type d'identifiant de chaîne en int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Impossible de convertir la chaîne en int.  %v", err)
	}

	// créer un utilisateur vide de type models.User
	var user models.User

	// décoder la requête json à l'utilisateur
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Printf("Impossible de décoder le corps de la requête.  %v", err)
	}

	// appeler update user pour mettre à jour l'utilisateur
	updatedRows := updateUser(int64(id), user)

	// formater la chaîne de message
	msg := fmt.Sprintf("L'utilisateur a été mis à jour avec succès. Nombre total de lignes/enregistrement concernés %v", updatedRows)

	// formater le message de réponse
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// envoyer la réponse
	json.NewEncoder(w).Encode(res)

}

// supprimer les détails de l'utilisateur dans la base de données postgres
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	// obtenir l'ID utilisateur à partir des paramètres de la requête, la clé est "id"
	params := mux.Vars(r)

	// convertir l'id dans la chaîne en int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Impossible de convertir la chaîne en int.  %v", err)
	}

	// appelez le deleteUser, convertissez l'int en int64
	deletedRows := deleteUser(int64(id))

	// format the message string
	msg := fmt.Sprintf("L'utilisateur a été mis à jour avec succès. Nombre total de lignes/enregistrement concernés %v", deletedRows)

	// formater le message de réponse
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// envoyer la réponse
	json.NewEncoder(w).Encode(res)

}

//------------------------- functions BDD ----------------

// insérer un utilisateur dans la BD
func insertUser(user models.User) int64 {

	// créer la connexion à la base de données postgres
	db := createConnection()

	// fermer la connexion à la base de données
	defer db.Close()

	// créer la requête sql d'insertion
	// renvoyer l'identifiant de l'utilisateur renverra l'identifiant de l'utilisateur inséré
	sqlStatement := `INSERT INTO users (name, kmmax, niveau) VALUES ($1, $2, $3) RETURNING userid`

	// l'identifiant inséré sera stocké dans cet identifiant
	var id int64

	// exécuter l'instruction sql
	// La fonction de numérisation enregistrera l'identifiant d'insertion dans l'identifiant
	err := db.QueryRow(sqlStatement, user.Name, user.KmMax, user.Niveau).Scan(&id)

	if err != nil {
		log.Printf("Impossible d'exécuter la requête. %v", err)
	}

	fmt.Printf("Insertion d'un seul enregistrement %v", id)

	// renvoie l'identifiant inséré
	return id
}

// obtenir un utilisateur de la base de données par son ID utilisateur
func getUser(id int64) (models.User, error) {
	// créer la connexion à la base de données postgres
	db := createConnection()

	// fermer la connexion à la base de données
	defer db.Close()

	// créer un utilisateur de modèles. Type d'utilisateur
	var user models.User

	// créer la requête select sql
	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	// exécuter l'instruction sql
	row := db.QueryRow(sqlStatement, id)

	// décoder l'objet de ligne à l'utilisateur
	err := row.Scan(&user.ID, &user.Name, &user.KmMax, &user.Niveau)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Aucune ligne n'a été renvoyée !")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Printf("Impossible d'analyser la ligne. %v", err)
	}

	// renvoie un utilisateur vide en cas d'erreur
	return user, err
}

// obtenir un utilisateur de la base de données par son ID utilisateur
func getAllUsers() ([]models.User, error) {
	// créer la connexion à la base de données postgres
	db := createConnection()

	// fermer la connexion à la base de données
	defer db.Close()

	var users []models.User

	// créer la requête select sql
	sqlStatement := `SELECT * FROM users`

	// exécuter l'instruction sql
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Printf("Impossible d'exécuter la requête. %v", err)
	}

	// fermer la déclaration
	defer rows.Close()

	// itérer sur les lignes
	for rows.Next() {
		var user models.User

		// décoder l'objet de ligne à l'utilisateur
		err = rows.Scan(&user.ID, &user.Name, &user.KmMax, &user.Niveau)

		if err != nil {
			log.Printf("Impossible d'analyser la ligne. %v", err)
		}

		// ajouter l'utilisateur dans la tranche des utilisateurs
		users = append(users, user)

	}

	// renvoie un utilisateur vide en cas d'erreur
	return users, err
}

// mettre à jour l'utilisateur dans la base de données
func updateUser(id int64, user models.User) int64 {

	// créer la connexion à la base de données postgres
	db := createConnection()

	// fermer la connexion à la base de données
	defer db.Close()

	// créer la requête sql de mise à jour
	sqlStatement := `UPDATE users SET name=$2, niveau=$3, kmmax=$4 WHERE userid=$1`

	// exécuter l'instruction sql
	res, err := db.Exec(sqlStatement, id, user.Name, user.Niveau, user.KmMax)

	if err != nil {
		log.Printf("Impossible d'exécuter la requête. %v", err)
	}

	// vérifier le nombre de lignes affectées
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Erreur lors de la vérification des lignes concernées. %v", err)
	}

	fmt.Printf("Nombre total de lignes/enregistrement concernés %v", rowsAffected)

	return rowsAffected
}

// supprimer l'utilisateur dans la base de données
func deleteUser(id int64) int64 {

	// créer la connexion à la base de données postgres
	db := createConnection()

	// fermer la connexion à la base de données
	defer db.Close()

	// créer la requête de suppression sql
	sqlStatement := `DELETE FROM users WHERE userid=$1`

	// exécuter l'instruction sql
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Printf("Impossible d'exécuter la requête. %v", err)
	}

	// vérifier le nombre de lignes affectées
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Erreur lors de la vérification des lignes concernées. %v", err)
	}

	fmt.Printf("Nombre total de lignes/enregistrement concernés %v", rowsAffected)

	return rowsAffected
}
