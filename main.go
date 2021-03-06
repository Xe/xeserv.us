package main

import (
	"net/http"
	"runtime"
	"time"

	"github.com/Xe/gurren/middleware/gurren"
	"github.com/Xe/middleware"
	"github.com/Xe/xeserv.us/interop/minecraft"
	"github.com/Xe/xeserv.us/interop/tf2"
	"github.com/Xe/xeserv.us/interop/xonotic"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/unrolled/render"
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

func fetchAndCache(name string, sl *gurren.StatsLog, r *http.Request, doer func() (interface{}, error)) (interface{}, error) {
	now := time.Now()

	if c, ok := caches[name]; ok {
		if now.Before(c.when.Add(time.Second * time.Duration(120))) {
			return c.data, nil
		}
	}

	c := &cache{
		when: now,
	}

	sl.Log(r, "Fetching data for "+name)
	var err error
	c.data, err = doer()
	if err != nil {
		return nil, err
	}

	caches[name] = c

	return c.data, nil
}

func main() {
	sl, err := gurren.New([]string{"http://logging.hyperadmin.yochat.biz:9200"}, "test", runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	re := render.New()

	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/index", rw, r, nil)
	})

	mux.Get("/rules", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/rules", rw, r, nil)
	})

	mux.Get("/chat", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("views/chat", rw, r, nil)
	})

	mux.Get("/tf2", func(rw http.ResponseWriter, r *http.Request) {
		s, err := fetchAndCache("tf2", sl, r, func() (interface{}, error) {
			return tf2.Query("10.0.0.5:27025", "cqcontrol")
		})
		if err != nil {
			handleError(rw, r, err)
		}

		doTemplate("views/tf2", rw, r, s)
	})

	mux.Get("/api/tf2.json", func(rw http.ResponseWriter, r *http.Request) {
		s, err := fetchAndCache("tf2", sl, r, func() (interface{}, error) {
			return tf2.Query("10.0.0.5:27025", "cqcontrol")
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		re.JSON(rw, http.StatusOK, s)
	})

	mux.Get("/minecraft", func(rw http.ResponseWriter, r *http.Request) {
		s, err := fetchAndCache("minecraft", sl, r, func() (interface{}, error) {
			return minecraft.Query("10.0.0.5", 25575, "swag")
		})
		if err != nil {
			handleError(rw, r, err)
		}

		doTemplate("views/minecraft", rw, r, s)
	})

	mux.Get("/api/minecraft.json", func(rw http.ResponseWriter, r *http.Request) {
		s, err := fetchAndCache("minecraft", sl, r, func() (interface{}, error) {
			return minecraft.Query("10.0.0.5", 25575, "swag")
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		re.JSON(rw, http.StatusOK, s)
	})

	mux.Get("/xonotic", func(rw http.ResponseWriter, r *http.Request) {
		c := xonotic.Dial("172.17.0.114", "26000")

		stats, err := fetchAndCache("xonotic", sl, r, func() (interface{}, error) {
			return c.Status()
		})
		if err != nil {
			handleError(rw, r, err)
		}

		doTemplate("views/xonotic", rw, r, stats)
	})

	mux.Get("/api/xonotic.json", func(rw http.ResponseWriter, r *http.Request) {
		c := xonotic.Dial("172.17.0.114", "26000")

		stats, err := fetchAndCache("xonotic", sl, r, func() (interface{}, error) {
			return c.Status()
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		re.JSON(rw, http.StatusOK, stats)
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

	tpl, err := ace.Load("views/layout", "views/error", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
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
