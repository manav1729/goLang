package main

import (
	"embed"
	"errors"
	"fmt"
	"goLangToDoApp/pkg/base"
	"goLangToDoApp/pkg/todo"
	"html/template"
	"log/slog"
	"net/http"
)

var (
	//go:embed static
	static embed.FS
)

var fileName string
var store *todo.ToDoStore

func main() {
	ctx := base.Init()
	fileName = base.DataFile

	var err error
	store, err = todo.NewToDoStore(fileName)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load to-do data", "error", err)
		return
	}

	slog.InfoContext(ctx, "Welcome to Manwendra's To-Do List Application.", "method", "ToDoListWeb")

	// Setup Http Server endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo/list", listFunc)

	// Serve static files for the /about endpoint
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	slog.InfoContext(ctx, "Http Server Listening on port 8081")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		msg := fmt.Sprintf("%s/n %s", "Http Server Listening error.", err)
		slog.ErrorContext(ctx, msg)
	}

	base.Exit(ctx)
}

func listFunc(res http.ResponseWriter, _ *http.Request) {
	ctx := base.Init()
	tmpl, err := template.ParseFiles("dynamic/list.html")
	if err != nil {
		msg := "Failed to load template."
		http.Error(res, msg, http.StatusInternalServerError)
		slog.ErrorContext(ctx, msg)
		return
	}

	items := store.GetAllToDoItems()
	_ = tmpl.Execute(res, items)
}
