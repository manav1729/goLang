package toDoUtil

import (
	"goLangToDoApp/base"
	"testing"
)

const tempFile = "test_ToDoAppData.json"

func TestToDoListApp(t *testing.T) {
	base.Init()

	items := []ToDoItem{{1, "Description 1", Statuses[0]}}
	SaveAllToDoItems(items, tempFile)
	items, _ = GetAllToDoItems(tempFile)

	testAddNewToDoItem(items, t)
	testUpdateToDoItemDesc(items, t)
	testUpdateToDoItemStatus(items, t)
	testDeleteToDoItemStatus(items, t)
}

func testAddNewToDoItem(items []ToDoItem, t *testing.T) {
	// Test Add New To-Do Item
	itemsGot := AddNewToDoItem(items, "Test Description")
	if len(itemsGot) != 2 {
		t.Errorf("Failed to Add New To-Do Item")
	}
}

func testUpdateToDoItemDesc(items []ToDoItem, t *testing.T) {
	// Test Update To-Do Item Desc
	itemsGot := UpdateToDoItem(items, 1, "Updated Description", "")
	if itemsGot[0].Description != "Updated Description" {
		t.Errorf("Failed to Update To-Do Item Description")
	}
}

func testUpdateToDoItemStatus(items []ToDoItem, t *testing.T) {
	// Test Update To-Do Status
	itemsGot := UpdateToDoItem(items, 1, "", "completed")
	if itemsGot[0].Status != "completed" {
		t.Errorf("Failed to Update To-Do Item Status")
	}
}

func testDeleteToDoItemStatus(items []ToDoItem, t *testing.T) {
	// Test Delete To-Do Item
	itemsGot := DeleteToDoItem(items, 1)
	if len(itemsGot) != 0 {
		t.Errorf("Failed to Delete To-Do Item Status")
	}
}
