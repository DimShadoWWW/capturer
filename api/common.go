package api

import "github.com/DimShadoWWW/hood"

// Api Struct
type Api struct {
	Db *hood.Hood
}

// TagContact Struct
type TagContact struct {
	Id        hood.Id
	ContactId int64
	TagId     int64
	Created   hood.Created `json:"created"`
	Updated   hood.Updated `json:"updated"`
}

func (table *TagContact) Indexes(indexes *hood.Indexes) {
	// indexes.Add("contact_index", "contact_id") // params: indexName, unique, columns...
	// indexes.Add("tag_index", "tag_id")         // params: indexName, unique, columns...
	indexes.AddUnique("tag_contact_index", "contact_id", "tag_id")
}
