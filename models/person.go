package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
)

//Person - This is person type
type Person struct {
	ID        string `json:"id" binding:"exists"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
}

//AddPersonToDB - Save
func (db *DB) AddPersonToDB(person Person) (*Person, error) {
	insForm, err := db.Prepare("INSERT INTO PERSON(ID, Firstname, Lastname) VALUES(?,?,?)")
	if err != nil {
		return nil, err
	}
	defer insForm.Close()
	_, er := insForm.Exec(person.ID, person.Firstname, person.Lastname)
	if er != nil {
		return nil, er
	}
	log.Println("INSERT: PERSON: " + person.ID)
	return &person, nil
}

//GetPerson - get single Person from DB
func (db *DB) GetPerson(id string) (*Person, error) {

	row := db.QueryRow("select * from PERSON where ID=?", id)

	person := new(Person)

	err := row.Scan(&person.ID, &person.Firstname, &person.Lastname)

	if err == sql.ErrNoRows {
		return nil, errors.New("No Person found matching your ID")
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return person, nil
}

//DeletePerson - delete single person
func (db *DB) DeletePerson(id string) (int64, error) {

	result, err := db.Exec("delete from PERSON where ID =?", id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

//UpdatePerson -
func (db *DB) UpdatePerson(person Person) (int64, error) {
	result, err := db.Exec("update PERSON set Firstname = ?, Lastname = ? where ID=? ", person.Firstname, person.Lastname, person.ID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()

}

//Unmarshal --- unmarshal person
func Unmarshal(data []byte, person *Person) error {

	err := json.Unmarshal(data, &person)

	if err != nil {
		return err
	}
	if person.ID == "" {
		return errors.New("json: ID is missing")
	}
	if person.Firstname == "" {
		return errors.New("json: Firstname is missing")
	}
	if person.Lastname == "" {
		return errors.New("json: Lastname is missing")
	}
	return nil
}
