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
	"github.com/rs/cors"
)

// Data type to store contact in address book
type Contact struct {
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	Organisation string `json:"Organisation"`
}

var Contacts []Contact

// Function that is triggered when / endpoint is hit
func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

// Function to list all contacts
func contactsList(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: List all contacts")
	fmt.Println("{}", Contacts)
	json.NewEncoder(response).Encode(Contacts)
}

// Function to list a single contact by Name
// R in CRUD
func readContactByID(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: Return contact by ID")
	id := mux.Vars(request)["ID"]
	for _, contact := range Contacts {
		if contact.ID == id {
			fmt.Println("{}", &contact)
			json.NewEncoder(response).Encode(&contact)
		}
	}
}

// Function to create a new contact
// C in CRUD
func createNewContact(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: Create new contact by ID")
	var contact Contact
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &contact)
	//Add to contact database
	Contacts = append(Contacts, contact)
	json.NewEncoder(response).Encode(contact)
}

// Function to update a contact
// U in CRUD
func updateContactByID(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: Update contact by ID")
	id := mux.Vars(request)["ID"]
	var updatedContact Contact
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &updatedContact)
	for index, contact := range Contacts {
		if contact.ID == id {
			fmt.Println("{}", &contact)
			Contacts = append(Contacts[:index], Contacts[index+1:]...)
			Contacts = append(Contacts, updatedContact)
		}
		fmt.Fprintf(response, "Successfully updated contact")
	}
}

// Function to delete a contact
// D in CRUD
func deleteContactByID(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: Delete Contact by ID")
	id := mux.Vars(request)["ID"]
	for index, contact := range Contacts {
		if contact.ID == id {
			fmt.Println("{}", &contact)
			Contacts = append(Contacts[:index], Contacts[index+1:]...)
		}
		fmt.Fprintf(response, "Successfully deleted contact")
	}
}

// Function that handles API requests
func poll() {
	//Create new mux router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/contacts", contactsList)
	router.HandleFunc("/contact", createNewContact).Methods("POST")
	router.HandleFunc("/contact/{ID}", updateContactByID).Methods("PUT")
	router.HandleFunc("/contact/{ID}", deleteContactByID).Methods("DELETE")
	router.HandleFunc("/contact/{ID}", readContactByID)
	//Specify allowed contacts
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":10000", handler))
}

// Main function
func main() {
	poll()
}
