package rest

func (s *Server) routes() {
	s.router.HandleFunc("/api/", s.handleAPI).Methods("GET")
	s.router.HandleFunc("/todolist/", s.createTodoList).Methods("POST")
	s.router.HandleFunc("/todolist/", s.getAllTodoLists).Methods("GET")
	s.router.HandleFunc("/todolist/{listID}", s.getTodoListByID).Methods("GET")
	s.router.HandleFunc("/todolist/{listID}", s.deleteListByID).Methods("DELETE")
	s.router.HandleFunc("/todolist/{listID}/item", s.addItemToList).Methods("POST")
	s.router.HandleFunc("/todolist/item/{itemID}", s.deleteItemByID).Methods("DELETE")
	s.router.HandleFunc("/todolist/item/{itemID}/{done}", s.setItemDone).Methods("PATCH")
}
