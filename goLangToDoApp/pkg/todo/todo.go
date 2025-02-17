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
	}
	err := store.loadAllToDoItems()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store *ToDoStore) AddNewToDoItem(desc string) error {
	id := 1
	itemNos := len(store.items)
	if itemNos > 0 {
		id = store.items[itemNos-1].ItemId + 1
	}

	store.items = append(store.items, Item{id, Statuses[0], desc})
	return store.saveAllToDoItems()
}

func (store *ToDoStore) UpdateToDoItem(id int, status string, desc string) error {
	if status != "" && !slices.Contains(Statuses, status) {
		return errors.New("status of To-Do Item is invalid")
	}

	for index, item := range store.items {
		if item.ItemId == id {
			if status != "" {
				store.items[index].Status = status
			}
			if desc != "" {
				store.items[index].Description = desc
			}
			return store.saveAllToDoItems()
		}
	}
	return errors.New("To-Do Item failed to update")
}

func (store *ToDoStore) DeleteToDoItem(id int) error {
	for index, item := range store.items {
		if item.ItemId == id {
			store.items = append(store.items[:index], store.items[index+1:]...)
			return store.saveAllToDoItems()
		}

	}
	return errors.New("To-Do Item failed to deleted")
}

func (store *ToDoStore) GetAllToDoItems() []Item {
	return store.items
}

func (store *ToDoStore) loadAllToDoItems() error {
	byteValue, err := ioutil.ReadFile(store.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading file %s: %w", store.filePath, err)
	}

	if len(byteValue) > 0 {
		err := json.Unmarshal(byteValue, &store.items)
		if err != nil {
			return fmt.Errorf("error unmarshalling To-Do items: %w", err)
		}
	}
	return nil
}

func (store *ToDoStore) saveAllToDoItems() error {
	// Open json file
	data, err := json.MarshalIndent(store.items, "", "\t")
	if err != nil {
		return fmt.Errorf("%s\n%s", "Error marshalling To-Do Item(s).", err)
	}

	err = ioutil.WriteFile(store.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("%s %s\n%s", "Error saving to file.", store.filePath, err)
	}

	return nil
}
