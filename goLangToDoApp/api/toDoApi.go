package api

import (
	"encoding/json"
	"goLangToDoApp/base"
	"io/ioutil"
	"os"
	"slices"
)

const FileName = "./data/ToDoAppData.json"

var Statuses = []string{"not-started", "started", "completed"}

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func AddNewToDoItem(currentItems []ToDoItem, desc string) []ToDoItem {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}

	currentItems = append(currentItems, ToDoItem{id, desc, Statuses[0]})
	base.LogInfo("New To-Do Item added to the List.")

	return currentItems
}

func UpdateToDoItem(currentItems []ToDoItem, id int, desc string, status string) []ToDoItem {
	if status != "" && !slices.Contains(Statuses, status) {
		base.LogError("status of To-Do Item is invalid.", "valid statuses", Statuses)
		return currentItems
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
	if success {
		base.LogInfo("To-Do Item updated.", "Id:", id)
	} else {
		base.LogWarn("To-Do Item not updated.", "Id:", id)
	}
	return currentItems
}

func DeleteToDoItem(currentItems []ToDoItem, id int) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
			success = true
		}
	}
	if success {
		base.LogInfo("To-Do Item deleted.", "Id:", id)
	} else {
		base.LogWarn("To-Do Item not deleted.", "Id:", id)
	}
	return currentItems
}

func SaveAllToDoItems(allItems []ToDoItem, fileName string) {
	// Open json file
	data, err1 := json.MarshalIndent(allItems, "", "\t")
	if err1 != nil {
		base.LogError("Error marshalling To-Do Item(s).", err1)
		return
	}

	err2 := ioutil.WriteFile(fileName, data, 0644)
	if err2 != nil {
		base.LogError("Error saving to file.", "fileName", fileName, err2)
		return
	}

	base.LogInfo("To-Do Items Saved to file.", "fileName", fileName)
}

func GetAllToDoItems(fileName string) ([]ToDoItem, error) {
	// Open local json file
	jsonFile, err1 := os.Open(fileName)
	if err1 != nil {
		base.LogError("Error Opening file.", "fileName", fileName, err1)
		return nil, err1
	}

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		base.LogError("Error Reading file.", "fileName", fileName, err2)
		return nil, err2
	}

	if byteValue != nil {
		// Parse the json file to ToDoItems
		var items []ToDoItem
		err3 := json.Unmarshal(byteValue, &items)
		if err3 != nil {
			base.LogError("Error Unmarshalling To-Do Item(s).", err3)
			return nil, err3
		}

		if items != nil && len(items) > 0 {
			return items, nil
		}
	}
	return nil, nil
}
