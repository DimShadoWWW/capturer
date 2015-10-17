package api

import (
	"time"

	"github.com/eaigner/hood"
)

// Contact Struct
type Contact struct {
	ID          hood.Id
	Name        string    `json:"name"`
	Weight      int64     `json:"weight",sql:"default(20)"`
	Title       string    `json:"title"`
	Tags        []*Tag    `json:"tags"`
	Description string    `json:"description"`
	Created     time.Time `json:"created_on"`
}
