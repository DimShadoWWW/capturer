package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/eaigner/hood"
)

// Tag Struct
type Tag struct {
	ID      hood.Id
	Name    string    `json:"name"`
	Weight  int64     `json:"weight",sql:"default(20)"`
	Created time.Time `json:"created_on"`
}

// GetTag confirm tag existence
func (a *Api) GetTag(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	var tag *Tag
	if err := a.Db.QueryRow(`SELECT * FROM tags WHERE name IN = $1`, name).Scan(&tag); err != nil {
		log.Fatal(err)
	}

	if tag == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(tag)
}

// GetAllTags get all tags
func (a *Api) GetAllTags(w rest.ResponseWriter, r *rest.Request) {

	var tags []Tag
	if err := a.Db.QueryRow(`SELECT * FROM tags`).Scan(&tags); err != nil {
		log.Fatal(err)
	}

	if tags == nil {
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
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stmt, err := a.Db.Prepare("INSERT INTO tags(name,weight,created_on) VALUES($1,$2,$3);")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(tag.Name, tag.Weight, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.LastInsertId())

	w.WriteJson(&tag)
}

// DeleteTag delete tag
func (a *Api) DeleteTag(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")

	stmt, err := a.Db.Prepare("DELETE FROM tags WHERE name=$1;")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := stmt.Exec(name); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}
