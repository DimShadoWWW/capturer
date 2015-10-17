package db

import (
	"github.com/eaigner/hood"
	"time"
)

type Tag struct {
	ID      hood.Id
	Name    string    `json:"name"`
	Weight  int64     `json:"weight",sql:"default(20)"`
	Created time.Time `json:"created_on"`
}

type Contact struct {
	ID          hood.Id
	Name        string     `json:"name"`
	Weight      int64      `json:"weight",sql:"default(20)"`
	Title       string     `json:"title"`
	Tags        []*api.Tag `json:"tags"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created_on"`
}
