package cli

import (
	"flag"
	"fmt"
	"goLangToDoApp/util"
	"log/slog"
)

func ToDoListCli() {
	util.LogInfo("Welcome to Manwendra's To-Do List Application.", slog.String("method", "ToDoListCli"))

	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	desc := flag.String("desc", "", "Description of Item in To-Do List")
	status := flag.String("status", "", "Status of Item in To-Do List")

	flag.Parse()

	// Load All To-Do Items from file
	items, _ := util.GetAllToDoItems(util.FileName)

	switch {
	case *add:
		// Add a new To-Do Item
		items, _ = util.AddNewToDoItem(items, *desc)
	case *update && *id != 0:
		// Update a To-Do Item
		items, _ = util.UpdateToDoItem(items, *id, *desc, *status)
	case *remove && *id != 0:
		// Delete a To-Do Item
		items, _ = util.DeleteToDoItem(items, *id)
	default:
		printFlagInstructions()
	}

	// Print All To-Do Item(s)
	printToDoItems(items)

	// Save All To-Do Items to file
	util.SaveAllToDoItems(items, util.FileName)
}

// private methods
func printFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================" +
		"\n-add -header=<name> -desc <description> to \"Add a new To-Do Item\"" +
		"\n-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"" +
		"\n-remove -id=<itemId> to \"Delete a To-Do Item\"" +
		"\n===========================================================================================")
}

func printToDoItems(items []util.ToDoItem) {
	if items != nil && len(items) > 0 {
		util.LogDebug("To-Do Item(s) list.", "To-Do Item(s)", items)
		fmt.Println("================================== Your To-Do Task Items ==================================")
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\nStatus: %s\n", item.ItemId, item.Description, item.Status)
		}
		fmt.Println("===========================================================================================")
	} else {
		util.LogInfo("No To-Do Item(s) in the List.")
	}
}
