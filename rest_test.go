package main

import (
	"log"
	"net/http"
	"testing"

	locApi "github.com/DimShadoWWW/capturer/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/eaigner/hood"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func TestTagApi(t *testing.T) {

	// db, err := sql.Open(dbconf.Driver, dbconf.Source)
	db, err := hood.Open("postgres", "user=postgres dbname=indata_test host=192.168.0.10 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	api := locApi.Api{Db: db}
	rapi := rest.NewApi()
	rapi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/tags", api.GetAllTags),
		rest.Post("/tags", api.PostTag),
		rest.Get("/tags/:name", api.GetTag),
		rest.Post("/tags/:name", api.UpdateTag),
		rest.Delete("/tags/:name", api.DeleteTag),
	)
	if err != nil {
		log.Fatal(err)
	}
	rapi.SetApp(router)

	// Start a transaction
	tx := db.Begin()

	tx.DropTableIfExists(locApi.Tag{})
	tx.DropTableIfExists(locApi.Contact{})
	tx.DropTableIfExists(locApi.TagContact{})

	tx.CreateTableIfNotExists(locApi.Tag{})
	tx.CreateTableIfNotExists(locApi.Contact{})
	tx.CreateTableIfNotExists(locApi.TagContact{})

	// Commit changes
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	handler := http.StripPrefix("/api", rapi.MakeHandler())

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()
	recorded.BodyIs(`[]`)

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/api/tags", &map[string]string{"name": "a"}))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()

	var tags []locApi.Tag
	err = recorded.DecodeJsonPayload(&tags)

	if len(tags) != 1 {
		t.Errorf(`Expected one element and got %d`, len(tags))
	}

	if tags[0].Name != "a" {
		t.Error(`Expected tag name "a"`)
	}

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("DELETE", "http://localhost/api/tags/"+tags[0].Name, nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()
	recorded.BodyIs(`[]`)
}
