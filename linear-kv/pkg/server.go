package pkg

import (
	"io"
	"net/http"
	"net/url"
)

type Handler struct {
	kv *KV
}

func NewHandler(kv *KV) *Handler {
	handler := &Handler{
		kv: kv,
	}
	return handler
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := r.RequestURI
	defer r.Body.Close()
	switch r.Method {
	case http.MethodGet:
		s, v, ok := handler.kv.Get(k)
		if !s {
			http.Error(w, "Fail to GET", http.StatusBadRequest)
		}
		if !ok {
			http.Error(w, "Fail to GET", http.StatusNotFound)
		}
		w.Write([]byte(v))
	case http.MethodPut:
		v, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Fail to PUT", http.StatusBadRequest)
		}
		s := handler.kv.Put(k, string(v))
		if !s {
			http.Error(w, "Fail to PUT", http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Server struct {
	addr string
	kv   *KV
}

func NewServer(addr string, kv *KV) *Server {
	server := &Server{
		addr: addr,
		kv:   kv,
	}
	go server.Run()
	return server
}

func (server *Server) Run() {
	url, _ := url.Parse(server.addr)
	http.ListenAndServe(url.Host, server.Handler())
}

func (server *Server) Handler() *Handler {
	return NewHandler(server.kv)
}
