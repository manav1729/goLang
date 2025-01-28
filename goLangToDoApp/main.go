package main

import (
	"goLangToDoApp/base"
	"goLangToDoApp/cli"
)

func main() {
	base.Init()

	// To-Do List CLI Application
	cli.ToDoListApp()

	base.Exit()
}
