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
	store := &ToDoStore{filePath: filePath}
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
	err := store.saveAllToDoItems()
	if err != nil {
		return err
	}

	return nil
}

func (store *ToDoStore) UpdateToDoItem(id int, status string, desc string) error {
	if status != "" && !slices.Contains(Statuses, status) {
		msg := "status of To-Do Item is invalid"
		return errors.New(msg)
	}

	success := false
	for index, item := range store.items {
		if item.ItemId == id {
			if status != "" {
				store.items[index].Status = status
			}
			if desc != "" {
				store.items[index].Description = desc
			}
			success = true
		}
	}
	if !success {
		msg := "To-Do Item failed to update"
		return errors.New(msg)
	} else {
		err := store.saveAllToDoItems()
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *ToDoStore) DeleteToDoItem(id int) error {
	success := false
	for index, item := range store.items {
		if item.ItemId == id {
			store.items = append(store.items[:index], store.items[index+1:]...)
			success = true
		}
	}
	if !success {
		msg := "To-Do Item failed to deleted"
		return errors.New(msg)
	} else {
		err := store.saveAllToDoItems()
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *ToDoStore) GetAllToDoItems() []Item {
	return store.items
}

func (store *ToDoStore) loadAllToDoItems() error {
	// Open local json file
	jsonFile, err1 := os.Open(store.filePath)
	if err1 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error opening file.", store.filePath, err1)
		return errors.New(msg)
	}

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error reading file.", store.filePath, err2)
		return errors.New(msg)
	}
	defer jsonFile.Close()

	if byteValue != nil {
		// Parse the json file to ToDoItems
		err3 := json.Unmarshal(byteValue, &store.items)
		if err3 != nil {
			msg := fmt.Sprintf("%s\n%s", "Error Unmarshalling To-Do Item(s).", err3)
			return errors.New(msg)
		}
	}
	return nil
}

func (store *ToDoStore) saveAllToDoItems() error {
	// Open json file
	data, err1 := json.MarshalIndent(store.items, "", "\t")
	if err1 != nil {
		msg := fmt.Sprintf("%s\n%s", "Error marshalling To-Do Item(s).", err1)
		return errors.New(msg)
	}

	err2 := ioutil.WriteFile(store.filePath, data, 0644)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error saving to file.", store.filePath, err2)
		return errors.New(msg)
	}

	return nil
}
