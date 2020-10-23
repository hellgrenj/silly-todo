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
