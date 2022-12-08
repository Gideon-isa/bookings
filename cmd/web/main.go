package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Gideon-isa/bookings/internal/config"
	"github.com/Gideon-isa/bookings/internal/handlers"
	"github.com/Gideon-isa/bookings/internal/helpers"
	"github.com/Gideon-isa/bookings/internal/models"
	"github.com/Gideon-isa/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var app config.AppConfig

var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s\n", portNumber)
	//	http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

	// What am I doing to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create template cache")

		// this is the error returned in case the run function encounters an issue
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	return nil
}
