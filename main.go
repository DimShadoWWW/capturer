package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	locApi "github.com/DimShadoWWW/capturer/api"
	"github.com/ant0ine/go-json-rest/rest"
	// "github.com/coreos/go-semver/semver"
	"github.com/kardianos/osext"
	_ "github.com/lib/pq"
)

// type semVerMiddleware struct {
// 	MinVersion string
// 	MaxVersion string
// }
//
// func (mw *semVerMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
//
// 	minVersion, err := semver.NewVersion(mw.MinVersion)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	maxVersion, err := semver.NewVersion(mw.MaxVersion)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	return func(writer rest.ResponseWriter, request *rest.Request) {
// 		version, err := semver.NewVersion(strings.TrimLeft(request.PathParam("version"), "v"))
// 		if err != nil {
// 			rest.Error(
// 				writer,
// 				"Invalid version: "+err.Error(),
// 				http.StatusBadRequest,
// 			)
// 			return
// 		}
//
// 		if version.LessThan(*minVersion) {
// 			rest.Error(
// 				writer,
// 				"Min supported version is "+minVersion.String(),
// 				http.StatusBadRequest,
// 			)
// 			return
// 		}
//
// 		if maxVersion.LessThan(*version) {
// 			rest.Error(
// 				writer,
// 				"Max supported version is "+maxVersion.String(),
// 				http.StatusBadRequest,
// 			)
// 			return
// 		}
//
// 		request.Env["VERSION"] = version
// 		handler(writer, request)
// 	}
// }

func load(path, env string) (config, error) {
	envs := make(map[string]config)
	if env == "" {
		env = "development"
	}
	f, err := os.Open(path)
	if err != nil {
		return config{}, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&envs)
	if err != nil {
		return config{}, err
	}
	fmt.Printf("%#v\n", envs)
	c, ok := envs[env]
	if !ok {
		return c, fmt.Errorf("config entry for specified environment '%s' not found", env)
	}
	return c, nil
}

type config struct {
	Driver string `json:"driver"`
	Source string `json:"source"`
}

var (
	db      *sql.DB
	flagEnv = flag.String("env", "development", "which DB environment to use")
)

func main() {

	flag.Parse()

	p, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	dbconf, err := load(path.Join(p, "db", "config.json"), *flagEnv)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(dbconf.Driver, dbconf.Source)
	if err != nil {
		log.Fatal(err)
	}

	// svmw := semVerMiddleware{
	// 	MinVersion: "1.0.0",
	// 	MaxVersion: "1.0.0",
	// }

	api := locApi.Api{Db: db}
	rapi := rest.NewApi()
	rapi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		rest.Get("/tags", api.GetAllTags),
		rest.Post("/tags", api.PostTag),
		rest.Get("/tags/:name", api.GetTag),
		rest.Delete("/tags/:name", api.DeleteTag),

		// rest.Get("/#version/message", svmw.MiddlewareFunc(
		// 	func(w rest.ResponseWriter, req *rest.Request) {
		// 		version := req.Env["VERSION"].(*semver.Version)
		// 		if version.Major == 2 {
		// 			// http://en.wikipedia.org/wiki/Second-system_effect
		// 			w.WriteJson(map[string]string{
		// 				"Body": "Hello broken World!",
		// 			})
		// 		} else {
		// 			w.WriteJson(map[string]string{
		// 				"Body": "Hello World!",
		// 			})
		// 		}
		// 	},
		// )),
	)
	if err != nil {
		log.Fatal(err)
	}
	rapi.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", rapi.MakeHandler()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
