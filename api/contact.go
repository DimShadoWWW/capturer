package api

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/eaigner/hood"
)

// Contact Struct
type Contact struct {
	Id          hood.Id
	Name        string       `json:"name" validate:"presence"`
	Weight      int64        `json:"weight" sql:"default(20)"`
	Title       string       `json:"title" validate:"presence"`
	Description string       `json:"description" validate:"presence"`
	Created     hood.Created `json:"created"`
	Updated     hood.Updated `json:"updated"`
}

func (table *Contact) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("cname_index", "name")  // params: indexName, unique, columns...
	indexes.AddUnique("title_index", "title") // params: indexName, unique, columns...
}

// GetContact confirm contact existence
func (a *Api) GetContact(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	var contacts []Contact
	err := a.Db.Where("name", "=", name).OrderBy("weight").OrderBy("name").Limit(1).Find(&contacts)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(contacts[0])
}

// GetAllContacts get all contacts
func (a *Api) GetAllContacts(w rest.ResponseWriter, r *rest.Request) {

	contacts := []Contact{}
	err := a.Db.OrderBy("name").Find(&contacts)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&contacts)
}

// PostContact add new contact
func (a *Api) PostContact(w rest.ResponseWriter, r *rest.Request) {
	log.Println("POST")
	contact := Contact{}
	err := r.DecodeJsonPayload(&contact)
	if err != nil {
		log.Println(err)
		rest.Error(w, "Data error", http.StatusInternalServerError)
		return
	}

	// Start a transaction
	tx := a.Db.Begin()
	_, err = tx.Save(&contact)
	if err != nil {
		log.Println(err)
		rest.Error(w, "Failed", http.StatusInternalServerError)
		return
	}

	// Commit changes
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		rest.Error(w, "Failed", http.StatusInternalServerError)
		return
	}

	w.WriteJson(&contact)
}

// UpdateContact add new contact
func (a *Api) UpdateContact(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")
	log.Println("POST")

	// contact := Contact{}

	var contacts []Contact
	err := a.Db.Where("name", "=", name).OrderBy("name").Limit(1).Find(&contacts)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	err = r.DecodeJsonPayload(&contacts[0])
	if err != nil {
		log.Println(err)
		rest.Error(w, "Data error", http.StatusInternalServerError)
		return
	}

	if len(contacts) == 1 {
		// Start a transaction
		tx := a.Db.Begin()
		_, err = tx.Save(&contacts[0])
		if err != nil {
			log.Println(err)
			rest.Error(w, "Failed", http.StatusInternalServerError)
			return
		}

		// Commit changes
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			rest.Error(w, "Failed", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(&contacts[0])
}

// DeleteContact delete contact
func (a *Api) DeleteContact(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	var contacts []Contact
	err := a.Db.Where("name", "=", name).OrderBy("name").Limit(1).Find(&contacts)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	if len(contacts) == 1 {
		// Start a transaction
		tx := a.Db.Begin()

		_, err := tx.Delete(&contacts[0])
		if err != nil {
			log.Println(err)
			rest.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		// Commit changes
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			rest.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

	} else {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
