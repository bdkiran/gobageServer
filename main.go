package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//Database credentials, used to connect to db
const (
	host     = "host_name"
	port     = 1111
	user     = "postgres"
	password = "password"
	dbname   = "name_of_database"
)

//gloabl variable so that db can be used accross other files
var db *sql.DB

func getUsers(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Request Recived")
	persons, err := QueryAllUsers(db)
	if err != nil {
		return
	}
	response, _ := json.Marshal(persons)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Recived")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Invalid user id")
		return
	}

	p, err := QueryUser(db, id)
	if err != nil {
		fmt.Println("User does not exist")
	}
	response, _ := json.Marshal(p)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	var p person
	fmt.Println("Create user called.")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Println("Invalid POST request")
		return
	}
	id, err := CreateNewUser(db, p.Age, p.FirstName, p.LastName, p.Email)
	if err != nil {
		panic(err)
	}
	// defer r.Body.Close()
	fmt.Println(id)

	response, _ := json.Marshal(p)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Invalid user id")
		return
	}
	fmt.Println("id to delete is: ", id)
	message, err := DeleteUser(db, id)
	if err != nil {
		panic(err)
	}
	response, _ := json.Marshal(message)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update user called.")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Invalid user id")
		return
	}
	var p person
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&p)
	if err != nil {
		fmt.Println("Invalid PUT request")
		return
	}
	message, err := UpdateUser(db, id, p.FirstName, p.LastName)
	if err != nil {
		panic(err)
		return
	}
	response, _ := json.Marshal(message)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//fmt.Println("Conected to database successfully.")

	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/user/{id:[0-9]+}", getUser).Methods("GET")
	r.HandleFunc("/user", createNewUser).Methods("POST")
	r.HandleFunc("/user/{id:[0-9]+}", deleteUser).Methods("DELETE")
	r.HandleFunc("/user/{id:[0-9]+}", updateUser).Methods("PUT")

	log.Fatal(http.ListenAndServe("localhost:8000", r))

}
