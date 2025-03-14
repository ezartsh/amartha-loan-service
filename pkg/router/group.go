package router

import (
	"loan-service/utils"
)

type GroupRoute struct {
	path   string
	router *Router
}

func (gr *GroupRoute) Get(path string, handler Request) {
	gr.router.Get(utils.TrimSuffixSlashOnPaths(gr.path, path), handler)
}

func (gr *GroupRoute) Post(path string, handler Request) {
	gr.router.Post(utils.TrimSuffixSlashOnPaths(gr.path, path), handler)
}

func (gr *GroupRoute) Put(path string, handler Request) {
	gr.router.Put(utils.TrimSuffixSlashOnPaths(gr.path, path), handler)
}

func (gr *GroupRoute) Delete(path string, handler Request) {
	gr.router.Delete(utils.TrimSuffixSlashOnPaths(gr.path, path), handler)
}
