package db

import (
	"github.com/eaigner/hood"
	"time"
)

type TagTable struct {
	Id        hood.Id
	Name      string    `json:"name"`
	Weight    int64     `json:"weight",sql:"default(20)"`
	CreatedOn time.Time `json:"created_on"`
}

func (table *TagTable) Indexes(indexes *hood.Indexes) {
	indexes.Add("tname_index", "name")
}

type ContactTable struct {
	Id          hood.Id
	Name        string    `json:"name"`
	Weight      int64     `json:"weight",sql:"default(20)"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_on"`
}

func (table *ContactTable) Indexes(indexes *hood.Indexes) {
	indexes.Add("cname_index", "name")
	indexes.Add("title_index", "title")
}

type TagsContactTable struct {
	Id        hood.Id
	ContactId int64
	TagId     int64
	CreatedOn time.Time
}

func (table *TagsContactTable) Indexes(indexes *hood.Indexes) {
	indexes.Add("contact_index", "contact_id")
	indexes.Add("tag_index", "tag_id")
	indexes.AddUnique("tag_contact", "contact_id", "tag_id")
}
