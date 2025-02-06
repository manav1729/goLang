package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"slices"
)

var Statuses = []string{"not-started", "started", "completed"}

// NewToDoStore initializes a new ToDoStore
func NewToDoStore(filePath string) (*ToDoStore, error) {
	store := &ToDoStore{
		filePath: filePath,
		commands: make(chan func(*[]Item, *int)),
	}

	var items []Item
	err := loadAllToDoItems(filePath, &items)
	if err != nil {
		return nil, err
	}

	// Determine the next ID (ensure IDs are continuous)
	nextID := 1
	if len(items) > 0 {
		nextID = items[len(items)-1].ItemId + 1
	}

	// Start the actor goroutine
	go store.actorLoop(items, nextID)

	return store, nil
}

func (store *ToDoStore) actorLoop(items []Item, nextID int) {
	for command := range store.commands {
		command(&items, &nextID)
		err := saveAllToDoItems(store.filePath, items)
		if err != nil {
			return
		}
	}
}

func (store *ToDoStore) AddNewToDoItem(desc string) error {
	store.commands <- func(items *[]Item, nextID *int) {
		*items = append(*items, Item{ItemId: *nextID, Description: desc, Status: Statuses[0]})
		*nextID++
	}
	return nil
}

func (store *ToDoStore) UpdateToDoItem(id int, desc string, status string) error {
	if status != "" && !slices.Contains(Statuses, status) {
		msg := "status of To-Do Item is invalid"
		return errors.New(msg)
	}

	errChan := make(chan error)
	store.commands <- func(items *[]Item, _ *int) {
		for i := range *items {
			if (*items)[i].ItemId == id {
				if desc == "" {
					desc = (*items)[i].Description
				}

				if status == "" {
					status = (*items)[i].Status
				}

				(*items)[i] = Item{ItemId: id, Description: desc, Status: status}
				errChan <- nil
				return
			}
		}
		errChan <- errors.New("To-Do Item not found")
	}
	return <-errChan
}

func (store *ToDoStore) DeleteToDoItem(id int) error {
	errChan := make(chan error)
	store.commands <- func(items *[]Item, _ *int) {
		for i := range *items {
			if (*items)[i].ItemId == id {
				*items = append((*items)[:i], (*items)[i+1:]...)
				errChan <- nil
				return
			}
		}
		errChan <- errors.New("To-Do Item not found")
	}
	return <-errChan
}

func (store *ToDoStore) GetAllToDoItems() []Item {
	response := make(chan []Item)
	store.commands <- func(items *[]Item, _ *int) {
		result := make([]Item, len(*items))
		copy(result, *items) // Return a copy to avoid race conditions
		response <- result
	}
	return <-response
}

func (store *ToDoStore) SaveAll() error {
	// Open json file
	data, err1 := json.MarshalIndent(store.GetAllToDoItems(), "", "\t")
	if err1 != nil {
		msg := fmt.Sprintf("%s %s", "Error marshalling To-Do Item(s).", err1)
		return errors.New(msg)
	}

	err2 := ioutil.WriteFile(store.filePath, data, 0644)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s %s", "Error saving to file.", store.filePath, err2)
		return errors.New(msg)
	}
	return nil
}

func loadAllToDoItems(filePath string, items *[]Item) error {
	byteValue, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		if os.IsNotExist(err1) {
			return nil
		}
		return errors.New(fmt.Sprintf("%s %s %s", "Error reading file.", filePath, err1))
	}

	// Parse the json file to ToDoItems
	err2 := json.Unmarshal(byteValue, items)
	if err2 != nil {
		return errors.New(fmt.Sprintf("%s %s", "Error Unmarshalling To-Do Item(s).", err2))
	}
	return nil
}

func saveAllToDoItems(filePath string, items []Item) error {
	// Open json file
	data, err1 := json.MarshalIndent(items, "", "\t")
	if err1 != nil {
		msg := fmt.Sprintf("%s %s", "Error marshalling To-Do Item(s).", err1)
		return errors.New(msg)
	}

	err2 := ioutil.WriteFile(filePath, data, 0644)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s %s", "Error saving to file.", filePath, err2)
		return errors.New(msg)
	}

	return nil
}
