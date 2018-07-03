package apiserver

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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

func (this *Router) handleFinal(r *BaseRequest) {
	if r.responseObj != nil {
		// Json headers
		r.Response.Header().Set("Content-Type", "application/json")

		// Errors?
		if r.responseObj.errored == false {
			// Nope
			r.responseObj.OK()
		}
		fmt.Fprintf(r.Response, "%s", r.responseObj.ToString(true))
	}
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
