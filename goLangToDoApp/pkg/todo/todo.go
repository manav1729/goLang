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

func AddNewToDoItem(currentItems []Item, desc string) ([]Item, error) {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}

	currentItems = append(currentItems, Item{id, desc, Statuses[0]})
	return currentItems, nil
}

func UpdateToDoItem(currentItems []Item, id int, desc string, status string) ([]Item, error) {
	if status != "" && !slices.Contains(Statuses, status) {
		msg := "status of To-Do Item is invalid"
		return currentItems, errors.New(msg)
	}

	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			if desc != "" {
				currentItems[index].Description = desc
			}
			if status != "" {
				currentItems[index].Status = status
			}
			success = true
		}
	}
	if !success {
		msg := "To-Do Item failed to update"
		return currentItems, errors.New(msg)
	}
	return currentItems, nil
}

func DeleteToDoItem(currentItems []Item, id int) ([]Item, error) {
	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
			success = true
		}
	}
	if !success {
		msg := "To-Do Item failed to deleted"
		return currentItems, errors.New(msg)
	}
	return currentItems, nil
}

func SaveAllToDoItems(allItems []Item, fileName string) error {
	// Open json file
	data, err1 := json.MarshalIndent(allItems, "", "\t")
	if err1 != nil {
		msg := fmt.Sprintf("%s\n%s", "Error marshalling To-Do Item(s).", err1)
		return errors.New(msg)
	}

	err2 := ioutil.WriteFile(fileName, data, 0644)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error saving to file.", fileName, err2)
		return errors.New(msg)
	}

	return nil
}

func GetAllToDoItems(fileName string) ([]Item, error) {
	// Open local json file
	jsonFile, err1 := os.Open(fileName)
	if err1 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error opening file.", fileName, err1)
		return nil, errors.New(msg)
	}

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		msg := fmt.Sprintf("%s %s\n%s", "Error reading file.", fileName, err2)
		return nil, errors.New(msg)
	}

	if byteValue != nil {
		// Parse the json file to ToDoItems
		var items []Item
		err3 := json.Unmarshal(byteValue, &items)
		if err3 != nil {
			msg := fmt.Sprintf("%s\n%s", "Error Unmarshalling To-Do Item(s).", err3)
			return nil, errors.New(msg)
		}

		if items != nil && len(items) > 0 {
			return items, nil
		}
	}
	return nil, nil
}
