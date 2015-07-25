package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/Xe/gurren/middleware/gurren"
	"github.com/Xe/middleware"
	"github.com/Xe/xeserv.us/interop/minecraft"
	"github.com/Xe/xeserv.us/interop/xonotic"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/yosssi/ace"
)

type cache struct {
	data interface{}
	when time.Time
}

var (
	caches map[string]*cache
)

func init() {
	caches = make(map[string]*cache)
}

func fetchAndCache(name string, doer func() (interface{}, error)) (interface{}, error) {
	now := time.Now()

	if c, ok := caches[name]; ok {
		if now.Before(c.when.Add(time.Second * time.Duration(15))) {
			return c.data, nil
		}
	}

	c := &cache{
		when: now,
	}

	log.Printf("Fetching data for %s", name)
	var err error
	c.data, err = doer()
	if err != nil {
		return nil, err
	}

	caches[name] = c

	return c.data, nil
}

func main() {

	sl, err := gurren.New([]string{"http://172.17.42.1:9200"}, "test", runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/index", rw, r, nil)
	})

	mux.Get("/rules", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/rules", rw, r, nil)
	})

	mux.Get("/tf2", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/tf2", rw, r, nil)
	})

	mux.Get("/chat", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/chat", rw, r, nil)
	})

	mux.Get("/minecraft", func(rw http.ResponseWriter, r *http.Request) {
		s, err := fetchAndCache("minecraft", func() (interface{}, error) {
			return minecraft.Query("10.0.0.5", 25575, "swag")
		})
		if err != nil {
			handleError(rw, r, err)
		}

		doTemplate("views/minecraft", rw, r, s)
	})

	mux.Get("/xonotic", func(rw http.ResponseWriter, r *http.Request) {
		c := xonotic.Dial("10.0.0.18", "26000")

		stats, err := fetchAndCache("xonotic", func() (interface{}, error) {
			return c.Status()
		})
		if err != nil {
			handleError(rw, r, err)
		}

		doTemplate("views/xonotic", rw, r, stats)
	})

	mux.Get("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/favicon.png", 301)
	})

	n := negroni.Classic()

	middleware.Inject(n)
	n.Use(sl)
	n.UseHandler(mux)

	n.Run(":3000")
}

func handleError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHeader(500)

	data := struct {
		Path   string
		Reason string
	}{
		Path:   r.URL.String(),
		Reason: err.Error(),
	}

	doTemplate("views/error", rw, r, data)
}

func doTemplate(name string, rw http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := ace.Load("views/layout", name, nil)
	if err != nil {
		handleError(rw, r, err)
		return
	}

	if err := tpl.Execute(rw, data); err != nil {
		handleError(rw, r, err)
		return
	}
}
