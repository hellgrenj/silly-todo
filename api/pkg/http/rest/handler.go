package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hellgrenj/silly-todo/pkg/domainerrors"

	"github.com/gorilla/mux"
	"github.com/hellgrenj/silly-todo/pkg/todo"
)

type operationResult struct {
	Result string `json:"result"`
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have reached %s!", r.URL.Path[1:])
}

func (s *Server) addItemToList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["listID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid list id"})
		return
	}
	var item todo.Item
	if err := s.decode(w, r, &item); err != nil {
		fmt.Fprintf(w, "Something went wrong %v", err)
		return
	}
	fmt.Printf("Received item %v for list id %v\n", item, listID)
	id, err := s.todoService.AddItem(item, listID)
	if err != nil {
		handleError(w, err)
		return
	}
	result := fmt.Sprintf("item created with id %v", id)
	json.NewEncoder(w).Encode(&operationResult{Result: result})

}

func (s *Server) getAllTodoLists(w http.ResponseWriter, r *http.Request) {
	list, err := s.todoService.GetAllLists()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(list)

}
func (s *Server) getTodoListByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["listID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid list id"})
		return
	}
	list, err := s.todoService.GetListByID(listID)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(list)

}
func (s *Server) createTodoList(w http.ResponseWriter, r *http.Request) {
	var l todo.List
	if err := s.decode(w, r, &l); err != nil {
		fmt.Fprintf(w, "Failed to parse request payload %v", err)
		return
	}
	fmt.Printf("Received item %v\n", l)
	id, err := s.todoService.AddList(l)
	if err != nil {
		handleError(w, err)
		return
	}
	result := fmt.Sprintf("todolist created with id %v", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&operationResult{Result: result})

}
func (s *Server) deleteListByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["listID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid list id"})
		return
	}
	err = s.todoService.DeleteListByID(listID)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&operationResult{Result: fmt.Sprintf("Successfully deleted list with id %v", listID)})
}
func (s *Server) deleteItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid item id"})
		return
	}
	err = s.todoService.DeleteItemByID(itemID)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&operationResult{Result: fmt.Sprintf("Successfully deleted item with id %v", itemID)})
}
func (s *Server) setItemDone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid item id"})
		return
	}
	done, err := strconv.ParseBool(vars["done"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "done needs to be a boolean (second arg)"})
		return
	}
	err = s.todoService.UpdateItemDone(itemID, done)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&operationResult{Result: fmt.Sprintf("successfully set item.done to %v", done)})
}

func handleError(w http.ResponseWriter, err error) {
	if notFoundError, ok := err.(domainerrors.NotFoundError); ok {
		fmt.Println(notFoundError.Error())
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(&operationResult{Result: notFoundError.ExternalError()})
		return
	}
	if conflictError, ok := err.(domainerrors.ConflictError); ok {
		fmt.Println(conflictError.Error())
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(&operationResult{Result: conflictError.ExternalError()})
		return
	}

	w.WriteHeader(500)
	log.Fatal("Something went wrong!")
}
