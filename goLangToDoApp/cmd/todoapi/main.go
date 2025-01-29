package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goLangToDoApp/pkg/base"
	"goLangToDoApp/pkg/todo"
	"net/http"
	"strconv"
)

const fileName = "../data/ToDoData.json"

var ctx context.Context

func main() {
	ctx = base.Init()

	base.LogInfo(ctx, "Welcome to Manwendra's To-Do List Application.", "method", "ToDoListApi")

	// Setup Http Server endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todoapp/create", createFunc)
	mux.HandleFunc("GET /todoapp/get", getFunc)
	mux.HandleFunc("PUT /todoapp/update", updateFunc)
	mux.HandleFunc("DELETE /todoapp/delete", deleteFunc)

	// Wrapping Handlers
	handler := createMiddleware(ctx, mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	base.LogInfo(ctx, "Http Server Listening on port 8080")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		base.LogError(ctx, "Http Server Listening error:", err)
	}

	base.Exit(ctx)
}

func createMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func createFunc(res http.ResponseWriter, req *http.Request) {
	var createReq struct {
		Description string `json:"description"`
	}
	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil || createReq.Description == "" {
		msg := "Invalid request body. Accepted payload: " +
			"\n{\n\"description\" : <Task Description>\n}"
		http.Error(res, msg, http.StatusBadRequest)
		base.LogError(ctx, msg)
		return
	}

	items, _ := todo.GetAllToDoItems(fileName)

	allItems, err := todo.AddNewToDoItem(items, createReq.Description)
	if err != nil {
		msg := fmt.Sprintf("%s/n%s", "Failed to create new To-Do Item.", err)
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	err = todo.SaveAllToDoItems(allItems, fileName)
	if err != nil {
		msg := fmt.Sprintf("%s/n%s", "Failed to save created To-Do Item.", err)
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	base.LogInfo(ctx, "Created new To-Do Item successfully")
	res.WriteHeader(http.StatusCreated)
}

func getFunc(res http.ResponseWriter, _ *http.Request) {
	items, err := todo.GetAllToDoItems(fileName)
	if err != nil {
		msg := fmt.Sprintf("%s/n%s", "Failed to get To-Do Items.", err)
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(items)
	if err != nil {
		base.LogError(ctx, "Failed to encode To-Do Items.")
		return
	}

	base.LogInfo(ctx, "Fetched To-Do Item(s).", "Item(s):", items)
	res.WriteHeader(http.StatusOK)
}

func updateFunc(res http.ResponseWriter, req *http.Request) {
	var updateReq struct {
		ItemId      int    `json:"id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil || (updateReq.ItemId == 0 || (updateReq.Description == "" && updateReq.Status == "")) {
		msg := "Invalid request body. Accepted payload: \n" +
			"{\n" +
			"\"id\" : <Task Id>,\n" +
			"\"description\" : <Task Description>\n," +
			"\"status\" : <Task Status>\n}"
		http.Error(res, msg, http.StatusBadRequest)
		base.LogError(ctx, msg)
		return
	}

	items, _ := todo.GetAllToDoItems(fileName)
	allItems, err := todo.UpdateToDoItem(items, updateReq.ItemId, updateReq.Description, updateReq.Status)
	if err != nil {
		msg := "Failed to update To-Do Item."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	err = todo.SaveAllToDoItems(allItems, fileName)
	if err != nil {
		msg := "Failed to Save To-Do Item(s)."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	base.LogInfo(ctx, "Updated To-Do Item successfully.", "Id", updateReq.ItemId)
	res.WriteHeader(http.StatusOK)
}

func deleteFunc(res http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		msg := "Missing 'id' query parameter."
		http.Error(res, msg, http.StatusBadRequest)
		base.LogError(ctx, msg)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "Invalid 'id' query parameter."
		http.Error(res, msg, http.StatusBadRequest)
		base.LogError(ctx, msg)
		return
	}

	items, _ := todo.GetAllToDoItems(fileName)
	allItems, err := todo.DeleteToDoItem(items, id)
	if err != nil {
		msg := "Failed to Delete To-Do Item."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
		return
	}

	err = todo.SaveAllToDoItems(allItems, fileName)
	if err != nil {
		msg := "Failed to Save To-Do Item(s)."
		http.Error(res, msg, http.StatusInternalServerError)
		base.LogError(ctx, msg)
	}

	base.LogInfo(ctx, "Deleted To-Do Item successfully.", "Id", id)
	res.WriteHeader(http.StatusOK)
}
