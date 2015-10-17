package api

import (
	"database/sql"
	"time"

	"github.com/eaigner/hood"
)

// Tag Struct
type Api struct {
	Db *sql.DB
}

// TagContact Struct
type TagContact struct {
	Id        hood.Id
	ContactId int64
	TagId     int64
	CreatedOn time.Time
}
