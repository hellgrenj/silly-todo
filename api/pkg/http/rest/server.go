package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hellgrenj/silly-todo/pkg/todo"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/hellgrenj/silly-todo/pkg/validation"
)

// Server is the http server struct
type Server struct {
	router      *mux.Router
	todoService todo.Service
}

// Middleware setting up CORS-headers and moving to router
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Authorization, X-Requested-With")
	s.router.ServeHTTP(w, r)
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v validation.Ok) error {
	// Improved this function further, see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return fmt.Errorf("unsupported content-type %s", value)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // request body max 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&v); err != nil {
		handelDecodeError(w, err)
		return err
	}

	return v.OK()
}

func handelDecodeError(w http.ResponseWriter, err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		http.Error(w, msg, http.StatusBadRequest)
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request body contains badly-formed JSON")
		http.Error(w, msg, http.StatusBadRequest)
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		http.Error(w, msg, http.StatusBadRequest)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		http.Error(w, msg, http.StatusBadRequest)
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		http.Error(w, msg, http.StatusBadRequest)
	case err.Error() == "http: request body too large":
		msg := "Request body must not be larger than 1MB"
		http.Error(w, msg, http.StatusRequestEntityTooLarge)
	default:
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// NewServer returns a new http server
func NewServer(todoService todo.Service) *Server {
	s := &Server{router: mux.NewRouter(), todoService: todoService}
	s.routes()
	fmt.Println("API running on port 8080")
	return s
}
