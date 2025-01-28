package main

import (
	"goLangToDoApp/api"
	"goLangToDoApp/cli"
	"goLangToDoApp/util"
)

func main() {
	util.Init()

	// To-Do List CLI Application
	cli.ToDoListCli()

	// To-Do List API Application
	api.ToDoListApi()

	util.Exit()
}
