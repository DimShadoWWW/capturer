package main

import (
	"github.com/DimShadoWWW/capturer/api"
	"github.com/eaigner/hood"
)

type TagTable api.Tag

func (m *M) CreateTagsTable_1445047167_Up(hd *hood.Hood) {
	hd.CreateTable(&TagTable{})
}

func (m *M) CreateTagsTable_1445047167_Down(hd *hood.Hood) {
	hd.DropTableIfExists(&TagTable{})
}

func (table *TagTable) Indexes(indexes *hood.Indexes) {
	indexes.Add("tname_index", "name") // params: indexName, unique, columns...
}
