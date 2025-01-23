package toDoUtil

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const fileName = "./toDoUtil/ToDoAppData.json"

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

func ToDoApp() {
	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")
	removeAll := flag.Bool("removeAll", false, "Remove all To-Do Items")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	header := flag.String("header", "", "Header")
	desc := flag.String("desc", "", "Description")

	flag.Parse()

	// Load All To-Do Items from file
	items := getAllToDoItems(fileName)

	if *add {
		// Add a new To-Do Item
		items = addNewToDoItem(items, *header, *desc)
		fmt.Println("=========================== New To-Do Item added to the List ==============================")
	} else if *update && *id != 0 {
		// Update a To-Do Item
		items = updateToDoItem(items, *id, *header, *desc)
	} else if *remove && *id != 0 {
		// Delete a To-Do Item
		items = removeToDoItem(items, *id)
	} else if *removeAll {
		// Delete all To-Do Item(s)
		items = nil
		fmt.Println("================================ All To-Do Item(s) deleted ================================")
	} else {
		printFlagInstructions()
	}

	// Print All To-Do Item(s)
	printToDoItems(items)

	// Save All To-Do Items to file
	err := saveAllToDoItems(items, fileName)
	if err != nil {
		fmt.Println("==========> Error saving file:", fileName, err)
	}
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
		fmt.Println("================================== Your To-Do Task Items ==================================")
		//fmt.Println(items)
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\n", item.ItemId, item.Header)
			fmt.Println(item.Description)
		}
		fmt.Println("===========================================================================================")
	} else {
		fmt.Println("============================== No To-Do Items in the List =================================")
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
		fmt.Println("================================== To-Do Item updated ====================================")
	} else {
		fmt.Println("================================ To-Do Item not updated ==================================")
	}
	return currentItems
}

func removeToDoItem(currentItems []ToDoItem, deleteItemId int) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == deleteItemId {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
			success = true
		}
	}
	if success {
		fmt.Println("================================== To-Do Item deleted ====================================")
	} else {
		fmt.Println("================================ To-Do Item not deleted ==================================")
	}
	return currentItems
}

func saveAllToDoItems(allItems []ToDoItem, fileName string) error {
	// Open json file
	data, err := json.MarshalIndent(allItems, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func getAllToDoItems(fileName string) []ToDoItem {
	// Open local json file
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("==========> Error Opening file:", fileName, err)
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println("==========> Error Reading file:", fileName, err)
		} else {
			if byteValue != nil {
				// Parse the json file to ToDoItems
				var items []ToDoItem
				err := json.Unmarshal(byteValue, &items)
				if err != nil {
					fmt.Println("==========> Error Unmarshalling file:", err)
				} else {
					if items != nil && len(items) > 0 {
						return items
					}
				}
			}
		}
	}
	return nil
}
