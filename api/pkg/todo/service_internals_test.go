package todo

// här kan jag testa internals
import (
	"testing"
)

func TestSortItemsInListByName(t *testing.T) {
	list := &List{Name: "test list"}
	list.Items = append(list.Items, Item{Name: "b"}, Item{Name: "a"})
	sortItemsInListByName(list)
	if list.Items[0].Name != "a" {
		t.Errorf("Expected the first item in the unsroted list to be a but it was %v", list.Items[0].Name)
	}
}

func TestSortListofListsByName(t *testing.T) {
	a := List{Name: "Lista A"}
	b := List{Name: "Lista B"}
	var listOfLists []List
	listOfLists = append(listOfLists, b, a) // adding b first
	sortListOfListsByName(listOfLists)
	if listOfLists[0].Name != "Lista A" {
		t.Errorf("Expected the first list in the list of lists to be Lista A but it was % v", listOfLists[0].Name)
	}
}
