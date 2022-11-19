package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kireeti-28/bookings/pkg/config"
	"github.com/kireeti-28/bookings/pkg/models"
)

var app *config.AppConfig
// NewTemplate sets the config for the tamplate package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) (*models.TemplateData) {
	return td
}
// RenderTemplates renders html templates
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc = make(map[string]*template.Template)
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}


	// get requested template from cache
	//log.Println(tc,tmpl)
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	//render the template

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// parsedTemp, parseErr := template.ParseFiles("./templates/" + tmpl,  "./templates/base.layout.html")
	// if parseErr != nil {
	// 	fmt.Println("Error Parsing", parseErr)
	// 	return
	// }
	// err := parsedTemp.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println(fmt.Sprintf("error parsing template %s with msg %s", tmpl, err))
	// 	return
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range thru all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err		
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err		
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}




// Simple Cache - Method 1

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// check to see if we already have the template ain our cache
// 	_, inMap := tc[t]

// 	if !inMap {
// 		// need to create template
// 		log.Println("Creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// we have template in cache
// 		log.Println("Using cached template")
// 	}

// 	tmpl =  tc[t]

// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string {
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.go.html",
// 	}

// 	//parse the template
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache (map)
// 	tc[t] = tmpl

// 	return nil
// }