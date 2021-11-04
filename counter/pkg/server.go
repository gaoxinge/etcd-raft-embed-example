package pkg

import (
	"fmt"
	"net/http"
	"net/url"
)

type Handler struct {
	counter *Counter
}

func NewHandler(counter *Counter) *Handler {
	handler := &Handler{
		counter: counter,
	}
	return handler
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodGet:
		counter := handler.counter.Get()
		w.Write([]byte(fmt.Sprintf("%d", counter)))
	case http.MethodPut:
		handler.counter.Inc()
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Server struct {
	addr    string
	counter *Counter
}

func NewServer(addr string, counter *Counter) *Server {
	server := &Server{
		addr:    addr,
		counter: counter,
	}
	go server.Run()
	return server
}

func (server *Server) Run() {
	url, _ := url.Parse(server.addr)
	http.ListenAndServe(url.Host, server.Handler())
}

func (server *Server) Handler() *Handler {
	return NewHandler(server.counter)
}
