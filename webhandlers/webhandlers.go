package webhandlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dannylesnik/http-inject-context/models"
	"github.com/gorilla/mux"
)

//GetPerson -
func GetPerson(w http.ResponseWriter, r *http.Request) {
	personID := mux.Vars(r)["id"]
	log.Printf(" Reuqest URI %s", r.RequestURI)

	log.Printf(" Person ID %s", personID)

	db, ok := r.Context().Value(models.SQLKEY).(*models.DB)

	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}
	person, err := db.GetPerson(personID)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't get Person", Message: err.Error(), Code: 404})
	} else if err != nil {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't get Person", Message: err.Error(), Code: 500})
	} else {
		json.NewEncoder(w).Encode(person)
	}
}

//UpdatePerson -
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	db, ok := r.Context().Value(models.SQLKEY).(*models.DB)

	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't read Request", Message: err.Error(), Code: 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(models.Error{Error: "Can't parse JSON Request", Message: err.Error(), Code: 400})
		} else {
			result, err := db.UpdatePerson(person)
			if err != nil {
				json.NewEncoder(w).Encode(models.Error{Error: "Can't Update Person!!", Message: err.Error(), Code: 500})
			} else if result == 0 {
				json.NewEncoder(w).Encode(models.Error{Error: "Person with such ID doesnt exist!!!", Message: errors.New("Query returned 0 affected records").Error(), Code: 404})
			} else {
				json.NewEncoder(w).Encode(person)
			}
		}
	}
}

//CreatePerson - 
func  CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	db, ok := r.Context().Value(models.SQLKEY).(*models.DB)

	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't read Request", Message:  err.Error(),Code: 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(models.Error{Error: "Can't parse JSON Request", Message: err.Error(), Code: 400})
		} else {
			_, err := db.AddPersonToDB(person)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(models.Error{Error: "Can't Save to DB!!",Message:  err.Error(),Code: 500})
			} else {
				json.NewEncoder(w).Encode(person)
			}
		}
	}
}

//DeletePerson - 
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	personID := mux.Vars(r)["id"]
	log.Printf(" Event ID %s", personID)
	db, ok := r.Context().Value(models.SQLKEY).(*models.DB)

	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}
	person, err := db.GetPerson(personID)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't get Person", Message: err.Error(), Code: 404})
	} else if err != nil {
		json.NewEncoder(w).Encode(models.Error{Error: "Can't delete Person", Message: err.Error(), Code: 500})
	} else {
		result, err := db.DeletePerson(personID)
		if err != nil {
			json.NewEncoder(w).Encode(models.Error{Error: "Can't delete Person", Message: err.Error(), Code: 500})
		} else if result == 0 {
			json.NewEncoder(w).Encode(models.Error{Error: "Can't delete Person", Message: "Person does not exist!!!", Code: 404})
		} else {
			json.NewEncoder(w).Encode(person)
		}
	}
}