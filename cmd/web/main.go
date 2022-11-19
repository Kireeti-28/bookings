package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kireeti-28/bookings/pkg/config"
	"github.com/kireeti-28/bookings/pkg/handlers"
	"github.com/kireeti-28/bookings/pkg/render"
)

const port = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	//change this to true when we are in Production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)


	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	log.Println("Starting App at port 8080")

	// err = http.ListenAndServe(port, nil)

	srv := &http.Server{
		Addr: port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
	
}