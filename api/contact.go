package api

import (
	"time"

	"github.com/eaigner/hood"
)

// Contact Struct
type Contact struct {
	Id          hood.Id
	Name        string    `json:"name"`
	Weight      int64     `json:"weight",sql:"default(20)"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_on"`
}
