package db

import (
	"github.com/DimShadoWWW/hood"
)

type Tag struct {
	Id      hood.Id
	Name    string       `json:"name" validate:"presence"`
	Weight  int64        `json:"weight" sql:"default(20)"`
	Created hood.Created `json:"created"`
	Updated hood.Updated `json:"updated"`
}

func (table *Tag) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("tname_index", "name")
}

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
	indexes.AddUnique("cname_index", "name")
	indexes.AddUnique("title_index", "title")
}

type TagContact struct {
	Id        hood.Id
	ContactId int64
	TagId     int64
	Created   hood.Created `json:"created"`
	Updated   hood.Updated `json:"updated"`
}

func (table *TagContact) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("tag_contact_index", "contact_id", "tag_id")
}
