package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const fileName = "./GoLangToDoApp.json"

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

func main() {
	fmt.Println("==================== Welcome to Manwendra's TO-DO List Application. =======================")
	PrintFlagInstructions()

	list := flag.Bool("list", false, "List all To-Do Items in List")
	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")
	removeAll := flag.Bool("removeAll", false, "Remove all To-Do Items")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	header := flag.String("header", "", "Header")
	desc := flag.String("desc", "", "Description")

	flag.Parse()

	// Load All To Do Items from file
	items := GetAllToDoItems(fileName)
	if items == nil {
		fmt.Println("============================== No To-Do Items in the List =================================")
	}

	// List
	if *list {
		if items != nil {
			PrintToDoItems(items)
		}
	}

	if *add {
		items = AddNewToDoItem(items, *header, *desc)
		fmt.Println("=========================== New To-Do Item added to the List ==============================")
	}

	if *update && *id != 0 {
		items = UpdateToDoItem(items, ToDoItem{*id, *header, *desc})
		fmt.Println("================================== To-Do Item updated ====================================")
	}

	if *remove && *id != 0 {
		items = RemoveToDoItem(items, *id)
		fmt.Println("=================================== To-Do Item removed ====================================")
	}

	if *removeAll {
		items = nil
		fmt.Println("================================ All To-Do Item(s) removed ================================")
	}

	// Save All To Do Items to file
	err := SaveAllToDoItems(items)
	if err != nil {
		fmt.Println("==========> Error saving file:", fileName, err)
	}
}

func PrintFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================")
	fmt.Println("-add -header=<name> -desc <description> to \"Add a new To-Do Item\"")
	fmt.Println("-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"")
	fmt.Println("-remove -id=<itemId> to \"Delete a To-Do Item\"")
	fmt.Println("-removeAll to \"Delete all To-Do Items\"")
	fmt.Println("-list to \"Print all To-Do Items in the List\"")
	fmt.Println("===========================================================================================")
}

func GetAllToDoItems(fileName string) []ToDoItem {
	// Open json file
	byteValue := OpenToDoItemsFile(fileName)
	if byteValue != nil {
		// Parse the json file to ToDoItems
		items := UnMarshalToDoItems(byteValue)
		if items != nil {
			return items
		}
	}
	return nil
}

func SaveAllToDoItems(allItems []ToDoItem) error {
	// Open json file
	data, err := json.MarshalIndent(allItems, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func AddNewToDoItem(currentItems []ToDoItem, header string, desc string) []ToDoItem {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}
	return append(currentItems, ToDoItem{id, header, desc})
}

func UpdateToDoItem(currentItems []ToDoItem, updateItem ToDoItem) []ToDoItem {
	for index, item := range currentItems {
		if item.ItemId == updateItem.ItemId {
			currentItems[index].Header = updateItem.Header
			currentItems[index].Description = updateItem.Description
		}
	}
	return currentItems
}

func RemoveToDoItem(currentItems []ToDoItem, deleteItemId int) []ToDoItem {
	for index, item := range currentItems {
		if item.ItemId == deleteItemId {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
		}
	}
	return currentItems
}

func OpenToDoItemsFile(fileName string) []byte {
	// Open local json file
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("==========> Error opening file:", fileName, err)
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println("==========> Error reading file:", fileName, err)
		} else {
			return byteValue
		}
	}
	return nil
}

func UnMarshalToDoItems(byteValue []byte) []ToDoItem {
	var items []ToDoItem
	err := json.Unmarshal(byteValue, &items)
	if err != nil {
		fmt.Println("==========> Error Unmarshalling file:", err)
	} else {
		return items
	}
	return nil
}

func PrintToDoItems(items []ToDoItem) {
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
}
