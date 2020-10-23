package todo

import (
	"fmt"
	"sort"
)

type service struct {
	r Repository
}

// Service provides todolist features
type Service interface {
	AddList(List) (int, error)
	GetAllLists() ([]List, error)
	GetListByID(id int) (List, error)
	AddItem(item Item, listID int) (int, error)
	DeleteListByID(id int) error
	DeleteItemByID(id int) error
	UpdateItemDone(id int, done bool) error
}

// Repository provides access to todolist repository.
type Repository interface {
	AddList(List) (int, error)
	GetAllLists() ([]List, error)
	GetListByID(id int) (List, error)
	AddItem(item Item, listID int) (int, error)
	DeleteListByID(id int) error
	DeleteItemByID(id int) error
	UpdateItemDone(id int, done bool) error
}

// NewService creates a todo service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddList(l List) (int, error) {
	return s.r.AddList(l)
}

func (s *service) GetAllLists() ([]List, error) {
	allLists, err := s.r.GetAllLists()
	sortListOfListsByName(allLists)
	fmt.Println("all todo lists")
	for _, l := range allLists {
		sortItemsInListByName(&l)
		fmt.Println(l)
	}

	return allLists, err
}
func (s *service) GetListByID(id int) (List, error) {
	l, err := s.r.GetListByID(id)
	sortItemsInListByName(&l)
	return l, err
}
func (s *service) AddItem(item Item, listID int) (int, error) {
	return s.r.AddItem(item, listID)
}
func (s *service) DeleteListByID(id int) error {
	return s.r.DeleteListByID(id)
}
func (s *service) DeleteItemByID(id int) error {
	return s.r.DeleteItemByID(id)
}
func (s *service) UpdateItemDone(id int, done bool) error {
	return s.r.UpdateItemDone(id, done)
}
func sortItemsInListByName(l *List) {
	sort.SliceStable(l.Items, func(i, j int) bool {
		return l.Items[i].Name < l.Items[j].Name
	})
}
func sortListOfListsByName(lists []List) {
	sort.SliceStable(lists, func(i, j int) bool {
		return lists[i].Name < lists[j].Name
	})
}
