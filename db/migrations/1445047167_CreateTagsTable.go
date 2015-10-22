package main

import (
	"github.com/DimShadoWWW/capturer/api"
	"github.com/DimShadoWWW/hood"
)

func (m *M) CreateTagsTable_1445047167_Up(hd *hood.Hood) {
	hd.CreateTable(&api.Tag{})
}

func (m *M) CreateTagsTable_1445047167_Down(hd *hood.Hood) {
	hd.DropTableIfExists(&api.Tag{})
}
