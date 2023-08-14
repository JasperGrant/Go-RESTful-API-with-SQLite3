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
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//Character data type is generic
//TODO: Replace with something cooler
type Character struct{
	gorm.Model
	Name string `json:"Name"`
	Age int16 	`json:"Age"`
	Organisation string `json:"Organisation"`
}


//Function that is triggered when / endpoint is hit
func homePage(response http.ResponseWriter, request *http.Request){
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

//Function to list all characters
func charactersList(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: List all characters")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()
	var Characters []Character
	db.Find(&Characters)
	fmt.Println("{}", Characters)
	json.NewEncoder(response).Encode(Characters)
}

//Function to list a single character by Name
//R in CRUD
func readCharacterByName(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Return character by name")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	fmt.Println("Endpoint Hit: Delete Character")
	name := mux.Vars(request)["Name"]
	var character Character
	db.Where("Name = ?", name).Find(&character)
	fmt.Println("{}", &character)
	json.NewEncoder(response).Encode(&character)
}

//Function to create a new character
//C in CRUD
func createNewCharacter(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Create new character by name")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	//Get request body
	reqBody,_ := ioutil.ReadAll(request.Body)
	//Create variable for character and unmarshal from json
	var character Character
	json.Unmarshal(reqBody, &character)
	//Add to character database
	db.Create(&character)
	json.NewEncoder(response).Encode(character)
}

//Function to update a character
//U in CRUD
func updateCharacterByName(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Update character by name")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	name := mux.Vars(request)["Name"]
	var character Character
	db.Where("Name = ?", name).Find(&character)
	db.Delete(&character)
	//Get request body
	reqBody,_ := ioutil.ReadAll(request.Body)
	//Create variable for character and unmarshal from json
	var updatedCharacter Character
	json.Unmarshal(reqBody, &updatedCharacter)
	db.Save(&updatedCharacter)
	fmt.Fprintf(response, "Successfully updated user")
}

//Function to delete a character
//D in CRUD
func deleteCharacterByName(response http.ResponseWriter, request *http.Request){
	fmt.Println("Endpoint Hit: Delete Character by name")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil{
		panic("Failed to connect to database")
	}
	defer db.Close()
	fmt.Println("Endpoint Hit: Delete Character")
	name := mux.Vars(request)["Name"]
	var character Character
	db.Where("Name = ?", name).Find(&character)
	fmt.Fprintf(response, "Successfully deleted user")
}

//Function that handles API requests
func poll(){
	//Create new mux router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/characters", charactersList)
	router.HandleFunc("/character", createNewCharacter).Methods("POST")
	router.HandleFunc("/character/{Name}", updateCharacterByName).Methods("PUT")
	router.HandleFunc("/character/{Name}", deleteCharacterByName).Methods("DELETE")
	router.HandleFunc("/character/{Name}", readCharacterByName)
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

	db.AutoMigrate(&Character{})
}

//Main function
func main(){
	initialMigration()
	poll()
}