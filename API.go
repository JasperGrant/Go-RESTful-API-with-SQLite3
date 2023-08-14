//Simple GO RESTful API to teach myself APIs
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
	//"github.com/jinzhu/gorm"
)

//Character data type is generic
//TODO: Replace with something cooler
type Character struct{
	//gorm.Model
	Id string `json:"Id"`
	Name string `json:"Name"`
	Age int16 	`json:"Age"`
	Organisation string `json:"Organisation"`
}

//Global array that stores people
//TODO: Communicate with SQL database instead
var Characters []Character

//Function that is triggered when / endpoint is hit
func homePage(response http.ResponseWriter, request *http.Request){
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

//Function to list all characters
func charactersList(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: List all characters")
	json.NewEncoder(response).Encode(Characters)
}

//Function to list a single character by ID
//R in CRUD
func returnSingleCharacter(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Return single character")
	key := mux.Vars(request)["id"]
	for _,character := range Characters{
		if character.Id == key{
			json.NewEncoder(response).Encode(character)
		}
	}
}

//Function to create a new character
//C in CRUD
func createNewCharacter(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Create new character")
	//Get request body
	reqBody,_ := ioutil.ReadAll(request.Body)
	//Create variable for character and unmarshal from json
	var character Character
	json.Unmarshal(reqBody, &character)
	//Add to character database
	Characters = append(Characters, character)
	json.NewEncoder(response).Encode(character)
}

//Function to update a character
//U in CRUD
//TODO: Current implementation does not check if updated character and old character have same Id
func updateCharacter(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Update character")
	//Get request body
	reqBody,_ := ioutil.ReadAll(request.Body)
	//Create variable for character and unmarshal from json
	var updatedCharacter Character
	json.Unmarshal(reqBody, &updatedCharacter)
	//Add updated character to array with original character gone
	key := mux.Vars(request)["id"]
	for index,character := range Characters{
		if character.Id == key{
			Characters = append(append(Characters[:index], Characters[index+1:]...,), updatedCharacter)
		}
	}
}

//Function to delete a character
//D in CRUD
func deleteCharacter(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Delete Character")
	key := mux.Vars(request)["id"]
	for index,character := range Characters{
		if character.Id == key{
			Characters = append(Characters[:index], Characters[index+1:]...)
		}
	}
}

//Function that handles API requests
func poll(){
	//Create new mux router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/characters", charactersList)
	router.HandleFunc("/character", createNewCharacter).Methods("POST")
	router.HandleFunc("/character/{id}", updateCharacter).Methods("PUT")
	router.HandleFunc("/character/{id}", deleteCharacter).Methods("DELETE")
	router.HandleFunc("/character/{id}", returnSingleCharacter)
	log.Fatal(http.ListenAndServe(":10000", router))
}

//Main function
func main(){
	//Generate dummy data
	Characters = []Character{
		{Id: "1", Name: "Spongebob Squarepants", Age: 49, Organisation: "Krusty Krab"},
		{Id: "2", Name: "Bruce Wayne", Age: 21, Organisation: "Wayne Enterprises"},
	}
	poll()
}