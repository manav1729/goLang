package todo

import (
	"log"
	"os"
	"testing"
)

const tempFile = "test_ToDoData.json"

func TestToDo(t *testing.T) {
	store := &ToDoStore{FilePath: tempFile}
	err := store.saveAllToDoItems()
	if err != nil {
		t.Fatalf("Failed to save to-do items: %v", err)
	}

	testGetAllToDoItems(store, t)
	testAddNewToDoItem(store, t)
	testUpdateToDoItemDesc(store, t)
	testUpdateToDoItemStatus(store, t)
	testDeleteToDoItemStatus(store, t)

	err = os.Remove(tempFile)
	if err != nil {
		log.Fatal(err)
	}
}

func testGetAllToDoItems(store *ToDoStore, t *testing.T) {
	err := store.GetAllToDoItems()
	if err != nil || len(store.Items) != 0 {
		t.Errorf("Failed to Get To-Do Item(s)")
	}
}

func testAddNewToDoItem(store *ToDoStore, t *testing.T) {
	// Test Add New To-Do Item
	err := store.AddNewToDoItem("Test Description")
	if err != nil {
		t.Errorf("Failed to Add New To-Do Item")
	}

	if len(store.Items) != 1 {
		t.Errorf("Failed to Add New To-Do Item")
	}
}

func testUpdateToDoItemDesc(store *ToDoStore, t *testing.T) {
	// Test Update To-Do Item Desc
	err := store.UpdateToDoItem(1, "Updated Description", "")
	if err != nil {
		t.Errorf("Failed to Update To-Do Item")
	}

	if store.Items[0].Description != "Updated Description" {
		t.Errorf("Failed to Update To-Do Item Description")
	}
}

func testUpdateToDoItemStatus(store *ToDoStore, t *testing.T) {
	// Test Update To-Do Status
	err := store.UpdateToDoItem(1, "", "completed")
	if err != nil {
		t.Errorf("Failed to Update To-Do Item")
	}

	if store.Items[0].Status != "completed" {
		t.Errorf("Failed to Update To-Do Item Status")
	}
}

func testDeleteToDoItemStatus(store *ToDoStore, t *testing.T) {
	// Test Delete To-Do Item
	err := store.DeleteToDoItem(1)
	if err != nil {
		t.Errorf("Failed to Delete To-Do Item")
	}

	if len(store.Items) != 0 {
		t.Errorf("Failed to Delete To-Do Item Status")
	}
}
