package api

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/eaigner/hood"
)

// Tag Struct
type Tag struct {
	Id      hood.Id
	Name    string       `json:"name" validate:"presence"`
	Weight  int64        `json:"weight" sql:"default(20)"`
	Created hood.Created `json:"created"`
	Updated hood.Updated `json:"updated"`
}

func (table *Tag) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("tname_index", "name") // params: indexName, unique, columns...
}

// GetTag confirm tag existence
func (a *Api) GetTag(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	var tags []Tag
	err := a.Db.Where("name", "=", name).OrderBy("weight").OrderBy("name").Limit(1).Find(&tags)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}
	if len(tags) == 0 {
		rest.NotFound(w, r)
		return
	}

	w.WriteJson(tags[0])
}

// GetAllTags get all tags
func (a *Api) GetAllTags(w rest.ResponseWriter, r *rest.Request) {

	tags := []Tag{}
	err := a.Db.OrderBy("name").Find(&tags)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&tags)
}

// PostTag add new tag
func (a *Api) PostTag(w rest.ResponseWriter, r *rest.Request) {
	log.Println("POST")
	tag := Tag{}
	err := r.DecodeJsonPayload(&tag)
	if err != nil {
		log.Println(err)
		rest.Error(w, "Data error", http.StatusInternalServerError)
		return
	}

	// Start a transaction
	tx := a.Db.Begin()
	_, err = tx.Save(&tag)
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
	log.Println("POST")
	log.Printf("%#v", &tag)
	w.WriteJson(&tag)
}

// UpdateTag add new tag
func (a *Api) UpdateTag(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")
	log.Println("POST")

	// tag := Tag{}

	var tags []Tag
	err := a.Db.Where("name", "=", name).OrderBy("name").Limit(1).Find(&tags)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	err = r.DecodeJsonPayload(&tags[0])
	if err != nil {
		log.Println(err)
		rest.Error(w, "Data error", http.StatusInternalServerError)
		return
	}

	if len(tags) == 1 {
		// Start a transaction
		tx := a.Db.Begin()
		_, err = tx.Save(&tags[0])
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

	w.WriteJson(&tags[0])
}

// DeleteTag delete tag
func (a *Api) DeleteTag(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	var tags []Tag
	err := a.Db.Where("name", "=", name).OrderBy("name").Limit(1).Find(&tags)
	if err != nil {
		log.Println(err)
		rest.NotFound(w, r)
		return
	}

	if len(tags) == 1 {
		// Start a transaction
		tx := a.Db.Begin()

		_, err := tx.Delete(&tags[0])
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
