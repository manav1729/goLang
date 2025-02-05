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
	store := &ToDoStore{FilePath: filePath}
	err := store.loadAllToDoItems()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store *ToDoStore) AddNewToDoItem(desc string) error {
	id := 1
	itemNos := len(store.Items)
	if itemNos > 0 {
		id = store.Items[itemNos-1].ItemId + 1
	}

	store.Items = append(store.Items, Item{id, desc, Statuses[0]})
	err := store.saveAllToDoItems()
	if err != nil {
		return err
	}

	return nil
}

func (store *ToDoStore) UpdateToDoItem(id int, desc string, status string) error {
	if status != "" && !slices.Contains(Statuses, status) {
		msg := "status of To-Do Item is invalid"
		return errors.New(msg)
	}

	success := false
	for index, item := range store.Items {
		if item.ItemId == id {
			if desc != "" {
				store.Items[index].Description = desc
			}
			if status != "" {
				store.Items[index].Status = status
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
	for index, item := range store.Items {
		if item.ItemId == id {
			store.Items = append(store.Items[:index], store.Items[index+1:]...)
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
	return store.Items
}

func (store *ToDoStore) loadAllToDoItems() error {
	// Open local json file
	jsonFile, err1 := os.Open(store.FilePath)
	if err1 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error opening file.", store.FilePath, err1)
		return errors.New(msg)
	}

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error reading file.", store.FilePath, err2)
		return errors.New(msg)
	}
	defer jsonFile.Close()

	if byteValue != nil {
		// Parse the json file to ToDoItems
		err3 := json.Unmarshal(byteValue, &store.Items)
		if err3 != nil {
			msg := fmt.Sprintf("%s\n%s", "Error Unmarshalling To-Do Item(s).", err3)
			return errors.New(msg)
		}
	}
	return nil
}

func (store *ToDoStore) saveAllToDoItems() error {
	// Open json file
	data, err1 := json.MarshalIndent(store.Items, "", "\t")
	if err1 != nil {
		msg := fmt.Sprintf("%s\n%s", "Error marshalling To-Do Item(s).", err1)
		return errors.New(msg)
	}

	err2 := ioutil.WriteFile(store.FilePath, data, 0644)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error saving to file.", store.FilePath, err2)
		return errors.New(msg)
	}

	return nil
}
