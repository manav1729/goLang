package todo

import (
	"testing"
)

const tempFile = "test_ToDoData.json"

func TestToDo(t *testing.T) {
	itemsGot := []Item{{1, "Description 1", Statuses[0]}}
	err := SaveAllToDoItems(itemsGot, tempFile)
	if err != nil {
		t.Errorf("Failed to Get To-Do Item(s)")
	}

	testGetAllToDoItems(t)
	testAddNewToDoItem(itemsGot, t)
	testUpdateToDoItemDesc(itemsGot, t)
	testUpdateToDoItemStatus(itemsGot, t)
	testDeleteToDoItemStatus(itemsGot, t)
}

func testGetAllToDoItems(t *testing.T) []Item {
	itemsGot, err := GetAllToDoItems(tempFile)
	if err != nil || len(itemsGot) != 1 {
		t.Errorf("Failed to Get To-Do Item(s)")
	}
	return itemsGot
}

func testAddNewToDoItem(items []Item, t *testing.T) {
	// Test Add New To-Do Item
	itemsGot, _ := AddNewToDoItem(items, "Test Description")
	if len(itemsGot) != 2 {
		t.Errorf("Failed to Add New To-Do Item")
	}
}

func testUpdateToDoItemDesc(items []Item, t *testing.T) {
	// Test Update To-Do Item Desc
	itemsGot, _ := UpdateToDoItem(items, 1, "Updated Description", "")
	if itemsGot[0].Description != "Updated Description" {
		t.Errorf("Failed to Update To-Do Item Description")
	}
}

func testUpdateToDoItemStatus(items []Item, t *testing.T) {
	// Test Update To-Do Status
	itemsGot, _ := UpdateToDoItem(items, 1, "", "completed")
	if itemsGot[0].Status != "completed" {
		t.Errorf("Failed to Update To-Do Item Status")
	}
}

func testDeleteToDoItemStatus(items []Item, t *testing.T) {
	// Test Delete To-Do Item
	itemsGot, _ := DeleteToDoItem(items, 1)
	if len(itemsGot) != 0 {
		t.Errorf("Failed to Delete To-Do Item Status")
	}
}
