package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/DimShadoWWW/capturer/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/davecgh/go-spew/spew"
	"github.com/DimShadoWWW/hood"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type TestTag struct {
	Request string             `json:"request"`
	Url     string             `json:"url"`
	Input   *map[string]string `json:"input"`
	Status  int                `json:"status"`
	Body    string             `json:"body"`
	Result  *api.Tag           `json:"result"`
}
type TestTagStruct struct {
	First *TestTag
	Last  *TestTag
}

func TestTagApi(t *testing.T) {

	tests := []TestTagStruct{
		{
			First: &TestTag{
				Request: "GET",
				Url:     "http://localhost/api/tags",
				Status:  200,
				Body:    `[]`,
			},
		},
		{
			First: &TestTag{
				Request: "POST",
				Url:     "http://localhost/api/tags",
				Input: &map[string]string{
					"name": "a",
				},
				Status: 200,
			},
			Last: &TestTag{
				Request: "GET",
				Url:     "http://localhost/api/tags",
				Input:   nil,
				Status:  200,
				Result: &api.Tag{
					Name:   "a",
					Weight: 20,
				},
			},
		},
	}

	// db, err := sql.Open(dbconf.Driver, dbconf.Source)
	db, err := hood.Open("postgres", "user=postgres dbname=indata_test host=192.168.0.10 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	serviceApi := api.Api{Db: db}
	rapi := rest.NewApi()
	rapi.Use(rest.DefaultCommonStack...)
	// rapi.Use(&rest.AccessLogApacheMiddleware{})
	// rapi.Use(&rest.GzipMiddleware{})
	rapi.Use(&rest.ContentTypeCheckerMiddleware{})

	router, err := rest.MakeRouter(
		rest.Get("/tags", serviceApi.GetAllTags),
		rest.Post("/tags", serviceApi.PostTag),
		rest.Get("/tags/:name", serviceApi.GetTag),
		rest.Post("/tags/:name", serviceApi.UpdateTag),
		rest.Delete("/tags/:name", serviceApi.DeleteTag),
	)
	if err != nil {
		log.Fatal(err)
	}
	rapi.SetApp(router)

	for _, testVars := range tests {
		// Start a transaction
		tx := db.Begin()

		tx.DropTableIfExists(api.Tag{})
		tx.DropTableIfExists(api.Contact{})
		tx.DropTableIfExists(api.TagContact{})

		tx.CreateTableIfNotExists(api.Tag{})
		tx.CreateTableIfNotExists(api.Contact{})
		tx.CreateTableIfNotExists(api.TagContact{})

		// Commit changes
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
		hndl := rapi.MakeHandler()
		handler := http.StripPrefix("/api", hndl)

		// recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags/notfound", nil))
		// recorded.CodeIs(404)
		// recorded.ContentTypeIsJson()
		// recorded.BodyIs("{\n  \"Error\": \"Resource not found\"\n}")

		// recorded = test.RunRequest(t, &handler, test.MakeSimpleRequest("GET", "http://1.2.3.4/user-notfound", nil))
		// recorded.CodeIs(404)
		// recorded.ContentTypeIsJson()
		// recorded.BodyIs(`{"Error":"Resource not found"}`)

		log.Printf("%s request to %s ", testVars.First.Request, testVars.First.Url)
		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest(testVars.First.Request, testVars.First.Url, testVars.First.Input))
		recorded.CodeIs(testVars.First.Status)
		recorded.ContentTypeIsJson()
		// recorded.ContentEncodingIsGzip()
		log.Printf("body %#v", recorded.Recorder.Body.String())
		if testVars.First.Body != "" {
			recorded.BodyIs(testVars.First.Body)
		}

		if testVars.Last != nil {
			log.Printf("%s request to %s ", testVars.Last.Request, testVars.Last.Url)

			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest(testVars.Last.Request, testVars.Last.Url, testVars.Last.Input))
			recorded.CodeIs(testVars.Last.Status)
			recorded.ContentTypeIsJson()
			// recorded.ContentEncodingIsGzip()

			if testVars.Last.Result != nil {
				var tags []api.Tag
				err = recorded.DecodeJsonPayload(&tags)

				log.Printf("tags %#v", tags)
				if len(tags) != 1 {
					t.Errorf(`Expected one element and got %d`, len(tags))
				}

				if testVars.Last.Result.Name != "" {
					if tags[0].Name != testVars.Last.Result.Name {
						t.Errorf(`Expected tag name "%s"`, testVars.Last.Result.Name)
					}
				}

				// if ok := testVars.Last.Result.Weight; ok {
				// if testVars.Last.Result.Weight {
				if tags[0].Weight != testVars.Last.Result.Weight {
					t.Errorf(`Expected tag weight "%s"`, testVars.Last.Result.Weight)
				}
				// }
				// }
			}
		}

	}
	//
	// recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags", nil))
	// recorded.CodeIs(200)
	// recorded.ContentTypeIsJson()
	// // recorded.ContentEncodingIsGzip()
	//
	// var tags []api.Tag
	// err = recorded.DecodeJsonPayload(&tags)
	//
	// if len(tags) != 1 {
	// 	t.Errorf(`Expected one element and got %d`, len(tags))
	// }
	//
	// if tags[0].Name != "a" {
	// 	t.Error(`Expected tag name "a"`)
	// }
	//
	// recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("DELETE", "http://localhost/api/tags/"+tags[0].Name, nil))
	// recorded.CodeIs(200)
	// recorded.ContentTypeIsJson()
	//
	// recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/tags", nil))
	// recorded.CodeIs(200)
	// recorded.ContentTypeIsJson()
	// // recorded.ContentEncodingIsGzip()
	// recorded.BodyIs(`[]`)
}

func TestContactApi(t *testing.T) {

	// db, err := sql.Open(dbconf.Driver, dbconf.Source)
	db, err := hood.Open("postgres", "user=postgres dbname=indata_test host=192.168.0.10 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	serviceApi := api.Api{Db: db}
	rapi := rest.NewApi()
	rapi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/contacts", serviceApi.GetAllContacts),
		rest.Post("/contacts", serviceApi.PostContact),
		rest.Get("/contacts/:name", serviceApi.GetContact),
		rest.Post("/contacts/:name", serviceApi.UpdateContact),
		rest.Delete("/contacts/:name", serviceApi.DeleteContact),
	)
	if err != nil {
		log.Fatal(err)
	}
	rapi.SetApp(router)

	// Start a transaction
	tx := db.Begin()

	tx.DropTableIfExists(api.Tag{})
	tx.DropTableIfExists(api.Contact{})
	tx.DropTableIfExists(api.TagContact{})

	tx.CreateTableIfNotExists(api.Tag{})
	tx.CreateTableIfNotExists(api.Contact{})
	tx.CreateTableIfNotExists(api.TagContact{})

	// Commit changes
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	handler := http.StripPrefix("/api", rapi.MakeHandler())

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/contacts", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()
	recorded.BodyIs(`[]`)

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/api/contacts",
		&map[string]string{
			"name":        "a",
			"title":       "a",
			"description": "a",
		}))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/contacts", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()

	var contacts []api.Contact
	err = recorded.DecodeJsonPayload(&contacts)

	if len(contacts) != 1 {
		t.Errorf(`Expected one element and got %d`, len(contacts))
	}

	if contacts[0].Name != "a" {
		t.Error(`Expected contact name "a"`)
	}

	spew.Dump(contacts)
	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("DELETE", "http://localhost/api/contacts/"+contacts[0].Name, nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/api/contacts", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	// recorded.ContentEncodingIsGzip()
	recorded.BodyIs(`[]`)
}
