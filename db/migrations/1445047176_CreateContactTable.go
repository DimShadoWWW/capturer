package main

import (
	"github.com/DimShadoWWW/capturer/api"
	"github.com/DimShadoWWW/hood"
)

func (m *M) CreateContactTable_1445047176_Up(hd *hood.Hood) {
	hd.CreateTable(&api.Contact{})
}

func (m *M) CreateContactTable_1445047176_Down(hd *hood.Hood) {
	hd.DropTable(&api.Contact{})
}
