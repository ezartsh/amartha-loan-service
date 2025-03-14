package router

import "github.com/gorilla/mux"

type Router struct {
	container *mux.Router
}

func NewRouter() *Router {
	muxRouter := mux.NewRouter()
	return &Router{
		container: muxRouter,
	}
}

func (r *Router) Group(path string, cb func(route GroupRoute)) {
	groupRoute := GroupRoute{
		path:   path,
		router: r,
	}
	cb(groupRoute)
}

func (r *Router) Get(path string, handler Request) {
	r.container.HandleFunc(path, requestHandler(handler)).Methods("GET")
}

func (r *Router) Post(path string, handler Request) {
	r.container.HandleFunc(path, requestHandler(handler)).Methods("POST")
}

func (r *Router) Put(path string, handler Request) {
	r.container.HandleFunc(path, requestHandler(handler)).Methods("PUT")
}

func (r *Router) Delete(path string, handler Request) {
	r.container.HandleFunc(path, requestHandler(handler)).Methods("DELETE")
}

func (r *Router) GetContainer() *mux.Router {
	return r.container
}
