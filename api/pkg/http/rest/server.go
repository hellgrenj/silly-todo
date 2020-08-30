package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellgrenj/silly-todo/pkg/todo"

	"github.com/gorilla/mux"
	"github.com/hellgrenj/silly-todo/pkg/validation"
)

// Server is the http server struct
type Server struct {
	router      *mux.Router
	todoService todo.Service
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	enableCors(&w) // TODO write as a midleware instead?
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	s.router.ServeHTTP(w, r)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Authorization, X-Requested-With")
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v validation.Ok) error {
	// TODO improve this function further, see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // request body max 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}
	return v.OK()
}

// NewServer returns a new http server
func NewServer(todoService todo.Service) *Server {
	s := &Server{router: mux.NewRouter(), todoService: todoService}
	s.routes()
	fmt.Println("API running on port 8080")
	return s
}
