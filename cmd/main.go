package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "harsh"
	dbName   = "student"
)

type Users struct {
	Id   int
	Age  int
	Name string
}

func main() {
	// Creating a router

	r := mux.NewRouter()

	r.HandleFunc("/user", createStudentHandler).Methods("POST")
	r.HandleFunc("/user/{id}", getStudentHandler).Methods("GET")
	r.HandleFunc("/user/{id}", updateStudentHandler).Methods("PUT")
	r.HandleFunc("/user/{id}", deleteStudentHandler).Methods("DELETE")

	// Start the HTTP server on port 8090
	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}

// Function to create a student and put into in DB

func createStudentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Parse json data from the request body
	var student Users
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received student data: ID=%d, Name=%s, Age=%d\n", student.Id, student.Name, student.Age)
	if err := CreateStudent(db, student.Id, student.Name, student.Age); err != nil {
		http.Error(w, "Failed to create a student user", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Failed to create a student user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
}

func CreateStudent(db *sql.DB, id int, name string, age int) error {
	query := "INSERT INTO users (id, name, age) VALUES (?, ?, ?)"
	fmt.Println("user created and inserted is", id, name, age)
	_, err := db.Exec(query, id, name, age)
	return err
}

// Reading Data
func getStudentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the ID from the given URL by route
	vars := mux.Vars(r)
	idStr := vars["id"]

	rollNo, err := strconv.Atoi(idStr)

	// Call get user function to fetch the user data from the database
	student, err := GetStudent(db, rollNo)

	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// convert the object to JSON and sent it in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func GetStudent(db *sql.DB, rollNo int) (*Users, error) {
	query := "SELECT * FROM users WHERE id = ?"

	row := db.QueryRow(query, rollNo)

	student := &Users{}

	fmt.Println("user id get is ", rollNo)
	err := row.Scan(&student.Id, &student.Name, &student.Age)
	fmt.Printf("Fetched user: ID=%d, Name=%s, Age=%d\n", student.Id, student.Name, student.Age)

	if err != nil {
		return nil, err
	}

	return student, nil
}

// code for update
func updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get ID and from url and then update it
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert Id into interger
	studentRoll, err := strconv.Atoi(idStr)

	var student Users
	err = json.NewDecoder(r.Body).Decode(&student)
	UpdateStudent(db, studentRoll, student.Age, student.Name)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User updated successfully")
}

func UpdateStudent(db *sql.DB, rollNo, age int, name string) error {
	query := "UPDATE users SET name = ?, age = ? WHERE id = ?"
	_, err := db.Exec(query, name, age, rollNo)
	if err != nil {
		return err
	}
	return nil
}

func deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert Id into interger
	userRoll, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	user := DeleteStudent(db, userRoll)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteStudent(db *sql.DB, userRoll int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, userRoll)
	if err != nil {
		return err
	}
	return nil
}
