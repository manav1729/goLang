package main

import (
	"flag"
	"fmt"
	"goLangToDoApp/toDoUtil"
)

const fileName = "./toDoUtil/ToDoAppData.json"

func main() {
	fmt.Println("==================== Welcome to Manwendra's TO-DO List Application. =======================")
	toDoUtil.PrintFlagInstructions()

	list := flag.Bool("list", false, "List all To-Do Items in List")
	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")
	removeAll := flag.Bool("removeAll", false, "Remove all To-Do Items")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	header := flag.String("header", "", "Header")
	desc := flag.String("desc", "", "Description")

	flag.Parse()

	// Load All To-Do Items from file
	items := toDoUtil.GetAllToDoItems(fileName)
	if items == nil {
		fmt.Println("============================== No To-Do Items in the List =================================")
	}

	// List all the To-Do Items
	if *list {
		if items != nil {
			toDoUtil.PrintToDoItems(items)
		}
	}

	// Add a new To-Do Item
	if *add {
		items = toDoUtil.AddNewToDoItem(items, *header, *desc)
		fmt.Println("=========================== New To-Do Item added to the List ==============================")
	}

	// Update a To-Do Item
	if *update && *id != 0 {
		items = toDoUtil.UpdateToDoItem(items, *id, *header, *desc)
		fmt.Println("================================== To-Do Item updated ====================================")
	}

	// Delete a To-Do Item
	if *remove && *id != 0 {
		items = toDoUtil.RemoveToDoItem(items, *id)
		fmt.Println("=================================== To-Do Item deleted ====================================")
	}

	// Delete all To-Do Item(s)
	if *removeAll {
		items = nil
		fmt.Println("================================ All To-Do Item(s) deleted ================================")
	}

	// Save All To-Do Items to file
	err := toDoUtil.SaveAllToDoItems(items, fileName)
	if err != nil {
		fmt.Println("==========> Error saving file:", fileName, err)
	}
}
