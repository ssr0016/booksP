package routing

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware func(httprouter.Handle) httprouter.Handle

// Router is a http.Handler that wraps httprouter.Router with additional features.
type Router struct {
	filterss []Middleware
	path     string
	router   *httprouter.Router
}

type modifiedFileAccess struct {
	http.Dir
}

func (mfa modifiedFileAccess) Open(path string) (result http.File, err error) {
	file, err := mfa.Dir.Open(path)
	if err != nil {
		return
	}

	info, err := file.Stat()
	if err != nil {
		return
	}
	if info.IsDir() {
		return mfa.Dir.Open("404 page not found")
	}

	return file, nil
}

// New returns *Router with a new initialized *httprouter.Router embedded.
func NewRouter() *Router {
	return &Router{router: httprouter.New()}
}

func (r *Router) joinPath(path string) string {
	if (r.path + path)[0] != '/' {
		panic("path should start with '/' in path '" + path + "'.")
	}

	return r.path + path
}

// Group returns new *Router with given path and filterss.
// It should be used for handles which have same path prefix or common filterss.
func (r *Router) Group(path string, m ...Middleware) *Router {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return &Router{
		filterss: append(m, r.filterss...),
		path:     r.joinPath(path),
		router:   r.router,
	}
}

// Use appends new filters to current Router.
func (r *Router) Use(m ...Middleware) *Router {
	r.filterss = append(m, r.filterss...)
	return r
}

// Handle registers a new request handle combined with filterss.
func (r *Router) Handle(method, path string, handle httprouter.Handle, middleware ...Middleware) {
	for i := len(middleware) - 1; i >= 0; i-- {
		handle = middleware[i](handle)
	}
	for _, v := range r.filterss {
		handle = v(handle)
	}
	r.router.Handle(method, r.joinPath(path), handle)
}

func (r *Router) GET(path string, handle httprouter.Handle, middleware ...Middleware) {
	r.Handle(http.MethodGet, path, handle, middleware...)
}

func (r *Router) PUT(path string, handle httprouter.Handle, middleware ...Middleware) {
	r.Handle(http.MethodPut, path, handle, middleware...)
}

func (r *Router) POST(path string, handle httprouter.Handle, middleware ...Middleware) {
	r.Handle(http.MethodPost, path, handle, middleware...)
}

func (r *Router) DELETE(path string, handle httprouter.Handle, middleware ...Middleware) {
	r.Handle(http.MethodDelete, path, handle, middleware...)
}

// Handler is an adapter for http.Handler.
func (r *Router) Handler(method, path string, handler http.Handler) {
	handle := func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		handler.ServeHTTP(w, req)
	}
	r.Handle(method, path, handle)
}

// HandlerFunc is an adapter for http.HandlerFunc.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

// Static serves files from given root directory.
func (r *Router) Static(path, root string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path should end with '/*filepath' in path '" + path + "'.")
	}

	base := r.joinPath(path[:len(path)-9])
	fileServer := http.StripPrefix(base, http.FileServer(modifiedFileAccess{http.Dir(root)}))

	r.Handler("GET", path, fileServer)
}

// File serves the named file.
func (r *Router) File(path, name string) {
	r.HandlerFunc("GET", path, func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, name)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
