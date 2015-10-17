package main

import (
	"github.com/DimShadoWWW/capturer/api"
	"github.com/eaigner/hood"
)

func (m *M) CreateTagsContactRelationTable_1445089666_Up(hd *hood.Hood) {
	hd.CreateTable(&api.TagContact{})
}

func (m *M) CreateTagsContactRelationTable_1445089666_Down(hd *hood.Hood) {
	hd.CreateTable(&api.TagContact{})
}
