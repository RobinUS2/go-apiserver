package apiserver

import (
	"compress/gzip"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Router struct {
	server     *http.Server
	router     *httprouter.Router
	listenConf string
}

func (r *Router) Router() *httprouter.Router {
	return r.router
}

type BaseHandle func(*BaseRequest)

func (this *Router) RegisterController(controller BaseControllerI) {
	// Init
	controller.InitController()

	// Custom action
	if controller.CustomActions() != nil && len(controller.CustomActions()) > 0 {
		for _, customAction := range controller.CustomActions() {
			route := customAction.FullRoute(controller, false)
			if customAction.Method == "GET" {
				this.GET(route, customAction.Handle)
			} else if customAction.Method == "PUT" {
				this.PUT(route, customAction.Handle)
			} else if customAction.Method == "POST" {
				this.POST(route, customAction.Handle)
			} else if customAction.Method == "DELETE" {
				this.DELETE(route, customAction.Handle)
			} else {
				panic("Method not supported")
			}
		}
	}
}

var zippers = sync.Pool{New: func() interface{} {
	return gzip.NewWriter(nil)
}}

func (this *Router) handleFinal(r *BaseRequest) {
	var respBytes []byte
	if r.responseBytes != nil {
		respBytes = r.responseBytes
	} else if r.responseObj != nil {
		// Json headers
		r.Response.Header().Set("Content-Type", "application/json")

		// Errors?
		if r.responseObj.hasError == false {
			// Nope
			r.responseObj.OK()
		} else {
			// Not good
			errCode := http.StatusBadRequest
			if r.responseObj.errorCode > 0 {
				errCode = r.responseObj.errorCode
			}
			r.Response.WriteHeader(errCode)

		}
		pretty := r.GetParam("pretty") == "1"
		respStr := fmt.Sprintf("%s", r.responseObj.ToString(pretty))
		respBytes = []byte(respStr)
	}

	if !strings.Contains(strings.ToLower(r.Request.Header.Get("Accept-Encoding")), "gzip") {
		// NO gzip
		r.Response.Write(respBytes)
	} else {
		// gzip header
		r.Response.Header().Set("Content-Encoding", "gzip")

		// Get a Writer from the Pool
		gz := zippers.Get().(*gzip.Writer)

		// When done, put the Writer back in to the Pool
		defer zippers.Put(gz)

		// We use Reset to set the writer we want to use.
		gz.Reset(r.Response)

		// write to gzip stream
		gz.Write(respBytes)

		// flush & close
		gz.Flush()
		gz.Close()
	}
}

func (this *Router) ServeFiles(path string, filesDir string) {
	this.router.ServeFiles(path, http.Dir(filesDir))
}

func (this *Router) GET(path string, handle BaseHandle) {
	log.Printf("Registered GET %s", path)
	this.router.GET(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wrappedRequest := newBaseRequest(w, r, params)
		handle(wrappedRequest)
		this.handleFinal(wrappedRequest)
	})
}

func (this *Router) PUT(path string, handle BaseHandle) {
	log.Printf("Registered PUT %s", path)
	this.router.PUT(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wrappedRequest := newBaseRequest(w, r, params)
		handle(wrappedRequest)
		this.handleFinal(wrappedRequest)
	})
}

func (this *Router) POST(path string, handle BaseHandle) {
	log.Printf("Registered POST %s", path)
	this.router.POST(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wrappedRequest := newBaseRequest(w, r, params)
		handle(wrappedRequest)
		this.handleFinal(wrappedRequest)
	})
}

func (this *Router) DELETE(path string, handle BaseHandle) {
	log.Printf("Registered DELETE %s", path)
	this.router.DELETE(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wrappedRequest := newBaseRequest(w, r, params)
		handle(wrappedRequest)
		this.handleFinal(wrappedRequest)
	})
}

func (this *Router) Listen(fork bool) {
	start := func() {
		log.Printf("Starting server %s", this.server.Addr)
		log.Fatal(this.server.ListenAndServe())
	}
	if fork {
		go start()
	} else {
		start()
	}
}

func New(listenConf string) *Router {
	r := &Router{}
	r.router = httprouter.New()
	r.server = &http.Server{
		Addr:    listenConf,
		Handler: r.router,
	}
	return r
}
