package main

import (
	"goLangToDoApp/base"
	"goLangToDoApp/toDoUtil"
)

func main() {
	base.Init()

	// To-Do List CLI Application
	toDoUtil.ToDoListApp()
}
