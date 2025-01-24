package toDoUtil

import (
	"encoding/json"
	"flag"
	"fmt"
	"goLangToDoApp/base"
	"io/ioutil"
	"log/slog"
	"os"
	"slices"
)

const fileName = "./toDoUtil/ToDoAppData.json"

var statuses = []string{"not-started", "started", "completed"}

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func ToDoListApp() {
	base.LogInfo("Welcome to Manwendra's To-Do List Application.", slog.String("method", "ToDoListApp"))

	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")
	removeAll := flag.Bool("removeAll", false, "Delete all To-Do Items")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	desc := flag.String("desc", "", "Description of Item in To-Do List")
	status := flag.String("status", "", "Status of Item in To-Do List")

	flag.Parse()

	// Load All To-Do Items from file
	items, _ := getAllToDoItems(fileName)

	switch {
	case *add:
		// Add a new To-Do Item
		items = AddNewToDoItem(items, *desc)
		base.LogInfo("New To-Do Item added to the List.")
	case *update && *id != 0:
		// Update a To-Do Item
		items = UpdateToDoItem(items, *id, *desc, *status)
	case *remove && *id != 0:
		// Delete a To-Do Item
		items = DeleteToDoItem(items, *id)
	case *removeAll:
		// Delete all To-Do Item(s)
		items = nil
		base.LogInfo("All To-Do Item(s) deleted.")
	default:
		printFlagInstructions()
	}

	// Print All To-Do Item(s)
	printToDoItems(items)

	// Save All To-Do Items to file
	SaveAllToDoItems(items, fileName)
}

// private methods

func printFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================" +
		"\n-add -header=<name> -desc <description> to \"Add a new To-Do Item\"" +
		"\n-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"" +
		"\n-remove -id=<itemId> to \"Delete a To-Do Item\"" +
		"\n-removeAll to \"Delete all To-Do Item(s)\"" +
		"\n===========================================================================================")
}

func printToDoItems(items []ToDoItem) {
	if items != nil && len(items) > 0 {
		base.LogDebug("To-Do Item(s) list.", "To-Do Item(s)", items)
		fmt.Println("================================== Your To-Do Task Items ==================================")
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\nStatus: %s\n", item.ItemId, item.Description, item.Status)
		}
		fmt.Println("===========================================================================================")
	} else {
		base.LogInfo("No To-Do Item(s) in the List.")
	}
}

func AddNewToDoItem(currentItems []ToDoItem, desc string) []ToDoItem {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}
	return append(currentItems, ToDoItem{id, desc, statuses[0]})
}

func UpdateToDoItem(currentItems []ToDoItem, id int, desc string, status string) []ToDoItem {
	if status != "" && !slices.Contains(statuses, status) {
		base.LogError("status of To-Do Item is invalid.", "valid statuses", statuses)
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

func getAllToDoItems(fileName string) ([]ToDoItem, error) {
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
