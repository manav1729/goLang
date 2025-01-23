package main

import (
	"flag"
	"fmt"
	"goLangToDoApp/toDoUtil"
)

const fileName = "./toDoUtil/ToDoAppData.json"

func main() {
	fmt.Println("==================== Welcome to Manwendra's To-Do List Application. =======================")

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

	if *add {
		// Add a new To-Do Item
		items = toDoUtil.AddNewToDoItem(items, *header, *desc)
		fmt.Println("=========================== New To-Do Item added to the List ==============================")
	} else if *update && *id != 0 {
		// Update a To-Do Item
		items = toDoUtil.UpdateToDoItem(items, *id, *header, *desc)
	} else if *remove && *id != 0 {
		// Delete a To-Do Item
		items = toDoUtil.RemoveToDoItem(items, *id)
	} else if *removeAll {
		// Delete all To-Do Item(s)
		items = nil
		fmt.Println("================================ All To-Do Item(s) deleted ================================")
	} else {
		toDoUtil.PrintFlagInstructions()
	}

	// Print All To-Do Item(s)
	toDoUtil.PrintToDoItems(items)

	// Save All To-Do Items to file
	err := toDoUtil.SaveAllToDoItems(items, fileName)
	if err != nil {
		fmt.Println("==========> Error saving file:", fileName, err)
	}
}
