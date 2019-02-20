package main

import (
	"database/sql"
	"fmt"
)

//Thes operations are for a table that has the following structure.
//Make changes according to table structure.
// TABLE users (
// 	id SERIAL PRIMARY KEY,
// 	age INT,
// 	first_name TEXT,
// 	last_name TEXT,
// 	email TEXT UNIQUE NOT NULL
// )

//Person object to represent table structure
type person struct {
	ID        int    `json:"id"`
	Age       int    `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

//QueryAllUsers retrives all users in DB
func QueryAllUsers(db *sql.DB) ([]person, error) {
	rows, err := db.Query("SELECT id, age, first_name, last_name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	persons := []person{}

	for rows.Next() {
		var p person
		err = rows.Scan(&p.ID, &p.Age, &p.FirstName, &p.LastName, &p.Email)
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return persons, err
}

//QueryUser querys 1 user
func QueryUser(db *sql.DB, id int) (person, error) {
	sqlStatement := `SELECT id, age, first_name, last_name, email FROM users WHERE id=$1;`
	var p person
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&p.ID, &p.Age, &p.FirstName, &p.LastName, &p.Email)
	switch err {
	case sql.ErrNoRows:
		return p, err
	case nil:
		//fmt.Println(user)
		return p, err
	default:
		return p, err
	}
}

//CreateNewUser Insert user into table
func CreateNewUser(db *sql.DB, age int, email string, firstName string, lastName string) (int, error) {
	//Insert user in table
	sqlStatement := `
	INSERT INTO users (age, email, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int
	err := db.QueryRow(sqlStatement, age, email, firstName, lastName).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//UpdateUser edits a single user
func UpdateUser(db *sql.DB, id int, firstName string, lastName string) (string, error) {
	//Update user in table
	sqlStatement := `
	UPDATE users
	SET first_name = $2, last_name = $3
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, id, firstName, lastName)
	if err != nil {
		return "Could not update table", err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return "Could not update table2", err
	}
	fmt.Println(count)
	return "Table updated successfully", err

}

//DeleteUser deletes a user from table
func DeleteUser(db *sql.DB, id int) (string, error) {
	//Delete user in table
	sqlStatement := `
	DELETE FROM users
	WHERE id=$1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return "User cannot be deleted, big error", err
	}

	return "User deleted.", err
}
