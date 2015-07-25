package main

import (
	"net/http"

	"github.com/Xe/middleware"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/yosssi/ace"
)

func main() {
	/*
		sl, err := gurren.New([]string{"http://172.17.42.1:9200"}, "test", runtime.NumCPU())
		if err != nil {
			panic(err)
		}
	*/

	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		tpl, err := ace.Load("views/layout", "views/index", nil)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tpl.Execute(rw, nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.Get("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/favicon.png", 301)
	})

	n := negroni.Classic()

	middleware.Inject(n)
	//n.Use(sl)
	n.UseHandler(mux)

	n.Run(":3000")
}
