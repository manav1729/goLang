package main

import (
	"context"
	"flag"
	"fmt"
	"goLangToDoApp/pkg/base"
	"goLangToDoApp/pkg/todo"
	"log/slog"
)

var fileName string

func main() {
	ctx := base.Init()
	fileName = base.DataFile

	slog.InfoContext(ctx, "Welcome to Manwendra's To-Do List Application.", "method", "ToDoListCli")

	add := flag.Bool("add", false, "Add a new To-Do Item to List")
	update := flag.Bool("update", false, "Update a To-Do Item")
	remove := flag.Bool("remove", false, "Delete a To-Do Item")

	id := flag.Int("id", 0, "ID of Item in To-Do List")
	desc := flag.String("desc", "", "Description of Item in To-Do List")
	status := flag.String("status", "", "Status of Item in To-Do List")

	flag.Parse()

	// Load All To-Do Items from file
	store, err := todo.NewToDoStore(fileName)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get item(s) of To-Do List:", "error", err)
	}

	switch {
	case *add:
		// Add a new To-Do Item
		err = store.AddNewToDoItem(*desc)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to add item to To-Do List:", "error", err)
		}
	case *update && *id != 0:
		// Update a To-Do Item
		err = store.UpdateToDoItem(*id, *desc, *status)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to update item to To-Do List:", "error", err)
		}
	case *remove && *id != 0:
		// Delete a To-Do Item
		err = store.DeleteToDoItem(*id)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to remove item from To-Do List:", "error", err)
		}
	default:
		printFlagInstructions()
	}

	// Print All To-Do Item(s)
	printToDoItems(ctx, store.GetAllToDoItems())

	base.Exit(ctx)
}

// private methods
func printFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================" +
		"\n-add -header=<name> -desc <description> to \"Add a new To-Do Item\"" +
		"\n-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"" +
		"\n-remove -id=<itemId> to \"Delete a To-Do Item\"" +
		"\n===========================================================================================")
}

func printToDoItems(ctx context.Context, items []todo.Item) {
	if items != nil && len(items) > 0 {
		slog.DebugContext(ctx, "To-Do Item(s) list.", "To-Do Item(s)", items)
		fmt.Println("================================== Your To-Do Task Items ==================================")
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\nStatus: %s\n", item.ItemId, item.Description, item.Status)
		}
		fmt.Println("===========================================================================================")
	} else {
		slog.InfoContext(ctx, "No To-Do Item(s) in the List.")
	}
}
