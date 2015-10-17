package main

import (
	"github.com/DimShadoWWW/capturer/api"
	"github.com/eaigner/hood"
)

type ContactTable api.Contact

func (m *M) CreateContactTable_1445047176_Up(hd *hood.Hood) {
	hd.CreateTable(&ContactTable{})
}

func (m *M) CreateContactTable_1445047176_Down(hd *hood.Hood) {
	hd.DropTable(&ContactTable{})
}

func (table *ContactTable) Indexes(indexes *hood.Indexes) {
	indexes.Add("cname_index", "name")  // params: indexName, unique, columns...
	indexes.Add("title_index", "title") // params: indexName, unique, columns...
}
