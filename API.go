//Simple GO RESTful Address book API to teach myself APIs
//Written by Jasper Grant
//2023-08-14

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//Data type to store contact in address book
type Contact struct{
	ID string `json:"ID"`
	Name string `json:"Name"`
	Organisation string `json:"Organisation"`
	Address string `json: "Address"`
}


//Function that is triggered when / endpoint is hit
func homePage(response http.ResponseWriter, request *http.Request){
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

//Function to list all contacts
func contactsList(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: List all contacts")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()
	var Contacts []Contact
	db.Find(&Contacts)
	fmt.Println("{}", Contacts)
	json.NewEncoder(response).Encode(Contacts)
}

//Function to list a single contact by Name
//R in CRUD
func readContactByID(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Return contact by ID")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	fmt.Println("Endpoint Hit: Delete Contact")
	id := mux.Vars(request)["ID"]
	var contact Contact
	db.Where("ID = ?", id).Find(&contact)
	fmt.Println("{}", &contact)
	json.NewEncoder(response).Encode(&contact)
}

//Function to create a new contact
//C in CRUD
func createNewContact(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Create new contact by ID")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	//Get request body
	reqBody,err := ioutil.ReadAll(request.Body)
	if err != nil{
		panic("Failed to get body from request")
	}
	//Create variable for contact and unmarshal from json
	var contact Contact
	json.Unmarshal(reqBody, &contact)
	//Add to contact database
	db.Create(&contact)
	json.NewEncoder(response).Encode(contact)
}

//Function to update a contact
//U in CRUD
func updateContactByID(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Update contact by ID")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	id := mux.Vars(request)["ID"]
	var contact Contact
	db.Where("ID = ?", id).Find(&contact)
	db.Delete(&contact)
	//Get request body
	reqBody,err := ioutil.ReadAll(request.Body)
	if err != nil{
		panic("Failed to get body from request")
	}
	//Create variable for contact and unmarshal from json
	var updatedContact Contact
	json.Unmarshal(reqBody, &updatedContact)
	db.Save(&updatedContact)
	fmt.Fprintf(response, "Successfully updated contact")
}

//Function to delete a contact
//D in CRUD
func deleteContactByID(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Delete Contact by ID")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	id := mux.Vars(request)["ID"]
	var contact Contact
	db.Where("ID = ?", id).Find(&contact)
	fmt.Fprintf(response, "Successfully deleted contact")
}

//Function that handles API requests
func poll(){
	//Create new mux router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/contacts", contactsList)
	router.HandleFunc("/contact", createNewContact).Methods("POST")
	router.HandleFunc("/contact/{ID}", updateContactByID).Methods("PUT")
	router.HandleFunc("/contact/{ID}", deleteContactByID).Methods("DELETE")
	router.HandleFunc("/contact/{ID}", readContactByID)
	log.Fatal(http.ListenAndServe(":10000", router))
}

//Function to set up initial database
func initialMigration(){
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Contact{})
}

//Main function
func main(){
	initialMigration()
	poll()
}