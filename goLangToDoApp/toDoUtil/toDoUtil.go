package toDoUtil

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
)

const fileName = "./toDoUtil/ToDoAppData.json"

var logger *slog.Logger

func initLogger() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

func ToDoListApp() {
	initLogger()
	logger.Info("Welcome to Manwendra's To-Do List Application.", slog.String("method", "ToDoListApp"))

	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")
	removeAll := flag.Bool("removeAll", false, "Delete all To-Do Items")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	header := flag.String("header", "", "Header")
	desc := flag.String("desc", "", "Description")

	flag.Parse()

	// Load All To-Do Items from file
	items, _ := getAllToDoItems(fileName)

	if *add {
		// Add a new To-Do Item
		items = addNewToDoItem(items, *header, *desc)
		logger.Info("New To-Do Item added to the List.")
	} else if *update && *id != 0 {
		// Update a To-Do Item
		items = updateToDoItem(items, *id, *header, *desc)
	} else if *remove && *id != 0 {
		// Delete a To-Do Item
		items = deleteToDoItem(items, *id)
	} else if *removeAll {
		// Delete all To-Do Item(s)
		items = nil
		logger.Info("All To-Do Item(s) deleted.")
	} else {
		printFlagInstructions()
	}

	// Print All To-Do Item(s)
	printToDoItems(items)

	// Save All To-Do Items to file
	saveAllToDoItems(items, fileName)
}

// private methods

func printFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================")
	fmt.Println("-add -header=<name> -desc <description> to \"Add a new To-Do Item\"")
	fmt.Println("-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"")
	fmt.Println("-remove -id=<itemId> to \"Delete a To-Do Item\"")
	fmt.Println("-removeAll to \"Delete all To-Do Item(s)\"")
	fmt.Println("===========================================================================================")
}

func printToDoItems(items []ToDoItem) {
	if items != nil && len(items) > 0 {
		logger.Debug("To-Do Item(s) list.", "To-Do Item(s)", items)
		fmt.Println("================================== Your To-Do Task Items ==================================")
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\n", item.ItemId, item.Header)
			fmt.Println(item.Description)
		}
		fmt.Println("===========================================================================================")
	} else {
		logger.Info("No To-Do Item(s) in the List.")
	}
}

func addNewToDoItem(currentItems []ToDoItem, header string, desc string) []ToDoItem {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}
	return append(currentItems, ToDoItem{id, header, desc})
}

func updateToDoItem(currentItems []ToDoItem, id int, header string, desc string) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			currentItems[index].Header = header
			currentItems[index].Description = desc
			success = true
		}
	}
	if success {
		logger.Info("To-Do Item updated.", "Id:", id)
	} else {
		logger.Warn("To-Do Item not updated.", "Id:", id)
	}
	return currentItems
}

func deleteToDoItem(currentItems []ToDoItem, id int) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
			success = true
		}
	}
	if success {
		logger.Info("To-Do Item deleted.", "Id:", id)
	} else {
		logger.Warn("To-Do Item not deleted.", "Id:", id)
	}
	return currentItems
}

func saveAllToDoItems(allItems []ToDoItem, fileName string) {
	// Open json file
	data, err1 := json.MarshalIndent(allItems, "", "\t")
	if err1 != nil {
		logger.Error("Error marshalling To-Do Item(s).", err1)
		return
	}

	err2 := ioutil.WriteFile(fileName, data, 0644)
	if err2 != nil {
		logger.Error("Error saving to file.", "fileName", fileName, err2)
		return
	}

	logger.Info("To-Do Items Saved to file.", "fileName", fileName)
}

func getAllToDoItems(fileName string) ([]ToDoItem, error) {
	// Open local json file
	jsonFile, err1 := os.Open(fileName)
	if err1 != nil {
		logger.Error("Error Opening file.", "fileName", fileName, err1)
		return nil, err1
	}

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		logger.Error("Error Reading file.", "fileName", fileName, err2)
		return nil, err2
	}

	if byteValue != nil {
		// Parse the json file to ToDoItems
		var items []ToDoItem
		err3 := json.Unmarshal(byteValue, &items)
		if err3 != nil {
			logger.Error("Error Unmarshalling To-Do Item(s).", err3)
			return nil, err3
		}

		if items != nil && len(items) > 0 {
			return items, nil
		}
	}
	return nil, nil
}
