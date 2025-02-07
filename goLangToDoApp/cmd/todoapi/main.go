package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goLangToDoApp/pkg/base"
	todo "goLangToDoApp/pkg/todoCon"
	"log/slog"
	"net/http"
	"strconv"
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

	slog.InfoContext(ctx, "Welcome to Manwendra's To-Do List Application.", "method", "ToDoListApi")

	// Setup Http Server endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todo/create", createFunc)
	mux.HandleFunc("GET /todo/get", getFunc)
	mux.HandleFunc("PUT /todo/update", updateFunc)
	mux.HandleFunc("DELETE /todo/delete", deleteFunc)

	// Wrapping Handlers
	handler := createMiddleware(ctx, mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	slog.InfoContext(ctx, "Http Server Listening on port 8080")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, "Http Server Listening error:", err)
	}

	base.Exit(ctx)
}

func createMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func createFunc(res http.ResponseWriter, req *http.Request) {
	ctx := base.Init()
	var createReq struct {
		Description string `json:"description"`
	}
	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil || createReq.Description == "" {
		msg := "Invalid request body. Accepted payload: " +
			"\n{\n\"description\" : <Task Description>\n}"
		http.Error(res, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
		return
	}

	err = store.AddNewToDoItem(createReq.Description)
	if err != nil {
		msg := fmt.Sprintf("%s/n%s", "Failed to create new To-Do Item.", err)
		http.Error(res, msg, http.StatusInternalServerError)
		slog.ErrorContext(ctx, msg)
		return
	}

	slog.InfoContext(ctx, "Created new To-Do Item successfully")
	res.WriteHeader(http.StatusCreated)
}

func getFunc(res http.ResponseWriter, _ *http.Request) {
	ctx := base.Init()
	items, err := store.GetAllToDoItems()
	if err != nil {
		msg := fmt.Sprintf("%s/n%s", "Failed to get all To-Do Items.", err)
		http.Error(res, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(items)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to encode To-Do Items.")
		return
	}

	slog.InfoContext(ctx, "Fetched To-Do Item(s).")
	slog.DebugContext(ctx, "Item(s):", items)
	res.WriteHeader(http.StatusOK)
}

func updateFunc(res http.ResponseWriter, req *http.Request) {
	ctx := base.Init()
	var updateReq struct {
		ItemId      int    `json:"id"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}
	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil || (updateReq.ItemId == 0 || (updateReq.Description == "" && updateReq.Status == "")) {
		msg := "Invalid request body. Accepted payload: \n" +
			"{\n" +
			"\"id\" : <Task Id>,\n" +
			"\"status\" : <Task Status>,\n" +
			"\"description\" : <Task Description>\n}"

		http.Error(res, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
		return
	}

	err = store.UpdateToDoItem(updateReq.ItemId, updateReq.Status, updateReq.Description)
	if err != nil {
		msg := "Failed to update To-Do Item."
		http.Error(res, msg, http.StatusInternalServerError)
		slog.ErrorContext(ctx, msg)
		return
	}

	slog.InfoContext(ctx, "Updated To-Do Item successfully.", "Id", updateReq.ItemId)
	res.WriteHeader(http.StatusOK)
}

func deleteFunc(res http.ResponseWriter, req *http.Request) {
	ctx := base.Init()
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		msg := "Missing 'id' query parameter."
		http.Error(res, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "Invalid 'id' query parameter."
		http.Error(res, msg, http.StatusBadRequest)
		slog.ErrorContext(ctx, msg)
		return
	}

	err = store.DeleteToDoItem(id)
	if err != nil {
		msg := "Failed to Delete To-Do Item."
		http.Error(res, msg, http.StatusInternalServerError)
		slog.ErrorContext(ctx, msg)
		return
	}

	slog.InfoContext(ctx, "Deleted To-Do Item successfully.", "Id", id)
	res.WriteHeader(http.StatusOK)
}
