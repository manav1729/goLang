package main

import (
	"bufio"
	"fmt"
	"goLangToDoApp/pkg/base"
	"goLangToDoApp/pkg/todo"
	"os"
	"strconv"
	"strings"
)

var fileName string

var commands = []string{"list", "add", "update", "delete", "exit"}

func main() {
	fileName = base.DataFile

	fmt.Println("Welcome to Manwendra's To-Do List Application.", "method", "ToDoListRepl")

	// Load All To-Do Items from file
	store, err := todo.NewToDoStore(fileName)
	if err != nil {
		fmt.Println("Failed to get item(s) of To-Do List:", "error", err)
	}

	fmt.Printf("Welcome to the To-Do Read-eval-print! Enter commands (%s, %s, %s, %s, %s).\n",
		commands[0], commands[1], commands[2], commands[3], commands[4])
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == commands[4] {
			fmt.Println("Exiting To-Do Read-eval-print...")
			break
		}

		handleCommand(input, store)
	}
}

func handleCommand(input string, store *todo.ToDoStore) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) < 1 {
		fmt.Printf("Invalid command. Accepted Commands are %s.\n", commands)
		return
	}

	switch parts[0] {
	case commands[0]:
		items := store.GetAllToDoItems()
		if len(items) == 0 {
			fmt.Println("No To-Do items found.")
			return
		}

		for _, item := range items {
			fmt.Printf("%d. %s\nStatus: %s\n", item.ItemId, item.Description, item.Status)
		}
	case commands[1]:
		if len(parts) < 2 {
			fmt.Println("Usage: add <description>")
			return
		}

		err := store.AddNewToDoItem(parts[1])
		if err != nil {
			fmt.Println("Failed to add item to To-Do List:", "error", err)
		} else {
			fmt.Println("Added item to To-Do List.")
		}
	case commands[2]:
		updateParts := strings.SplitN(input, " ", 4)
		if len(updateParts) < 3 {
			fmt.Println("Usage: update <id> <status> <new_description>")
			return
		}
		id, err := strconv.Atoi(updateParts[1])
		if err != nil {
			fmt.Println("Invalid ID.")
			return
		}

		err = store.UpdateToDoItem(id, updateParts[2], updateParts[3])
		if err != nil {
			fmt.Println("Failed to update item to To-Do List:", err)
		} else {
			fmt.Println("To-Do item updated.")
		}
	case commands[3]:
		if len(parts) < 2 {
			fmt.Println("Usage: delete <id>")
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid ID.")
			return
		}

		err = store.DeleteToDoItem(id)
		if err != nil {
			fmt.Println("Failed to delete item from To-Do List:", err)
		} else {
			fmt.Println("To-Do item deleted.")
		}
	default:
		fmt.Printf("Unknown command. "+
			"\nAccepted Commands are %s."+
			"\nUsage: "+
			"\nadd <description>"+
			"\nupdate <id> <status> <new_description>"+
			"\ndelete <id>\n", commands)
	}
}
