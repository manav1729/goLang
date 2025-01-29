package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"goLangToDoApp/pkg/base"
	"goLangToDoApp/pkg/todo"
	"html/template"
	"net/http"
)

var (
	//go:embed static
	static embed.FS
)

var ctx context.Context
var fileName string

func main() {
	ctx = base.Init()
	fileName = base.DataFile

	base.LogInfo(ctx, "Welcome to Manwendra's To-Do List Application.", "method", "ToDoListWeb")

	// Setup Http Server endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo/list", listFunc)

	// Serve static files for the /about endpoint
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	base.LogInfo(ctx, "Http Server Listening on port 8081")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		msg := fmt.Sprintf("%s/n %s", "Http Server Listening error.", err)
		base.LogError(ctx, msg)
	}

	base.Exit(ctx)
}

func listFunc(res http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("dynamic/list.html")
	if err != nil {
		msg := "Failed to load template."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	items, err := todo.GetAllToDoItems(fileName)
	if err != nil {
		msg := "Failed to load To-Do Items."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}
	_ = tmpl.Execute(res, items)
}
