// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// Start Parse the configuration file and starts the ipe server
// It Panic if could not start the HTTP or HTTPS server
func Start(filename string) {
	var conf configFile

	rand.Seed(time.Now().Unix())
	file, err := os.Open(filename)

	if err != nil {
		log.Error(err)
		return
	}

	defer file.Close()

	// Reading config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Error(err)
		return
	}

	StartMongoDBInstance(conf.MongoHost, conf.MongoDatabase)

	// Using a in memory database
	db := newMemdb()

	// Adding applications
	for _, a := range conf.Apps {
		db.AddApp(newAppFromConfig(a))
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/app/{key}", (&websocketHandler{db}).ServeHTTP)
	r.Group(func(r chi.Router) {
		r.Use(checkAppDisabled(db))
		r.Use(authenticationHandler(db))

		r.Post("/apps/{app_id}/events", (&postEventsHandler{db}).ServeHTTP)
		r.Get("/apps/{app_id}/channels", (&getChannelsHandler{db}).ServeHTTP)
		r.Get("/apps/{app_id}/channels/{channel_name}", (&getChannelHandler{db}).ServeHTTP)
		r.Get("/apps/{app_id}/channels/{channel_name}/users", (&getChannelUsersHandler{db}).ServeHTTP)
	})

	// Not working
	if conf.Profiling {
		r.Mount("/debug", middleware.Profiler())
	}

	if conf.SSL {
		go func() {
			log.Infof("Starting HTTPS service on %s ...", conf.SSLHost)
			log.Fatal(http.ListenAndServeTLS(conf.SSLHost, conf.SSLCertFile, conf.SSLKeyFile, r))
		}()
	}
	log.Infoln(r.Routes())
	log.Infof("Starting HTTP service on %s ...", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, r))
}
