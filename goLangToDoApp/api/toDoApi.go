package api

import (
	"encoding/json"
	"errors"
	"goLangToDoApp/util"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

func ToDoListApi() {
	util.LogInfo("Welcome to Manwendra's To-Do List Application.", slog.String("method", "ToDoListApi"))

	// Setup Http Server endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/todoapp/create", createFunc())
	mux.HandleFunc("/todoapp/get", getFunc())
	mux.HandleFunc("/todoapp/update", updateFunc())
	mux.HandleFunc("/todoapp/delete", deleteFunc())
	mux.HandleFunc("/todoapp/list", listFunc())

	// Serve static files for the /about endpoint
	dir := http.Dir("/webserver/static/")
	mux.Handle("/todoapp/about", http.FileServer(dir))

	// Wrapping Handlers
	handler := util.CreateMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	util.LogInfo("Http Server Listening on port 8080")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		util.LogError("Http Server Listening error:", err)
	}
}

func createFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			msg := "Invalid Method. Accepted methods: POST"
			http.Error(res, msg, http.StatusMethodNotAllowed)
			util.LogError(msg)
			return
		}

		var payload struct {
			Description string `json:"description"`
		}
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil || payload.Description == "" {
			msg := "Invalid request body. Accepted payload: " +
				"\n{\n\"description\" : <Task Description>\n}"
			http.Error(res, msg, http.StatusBadRequest)
			util.LogError(msg)
			return
		}

		items, _ := util.GetAllToDoItems(util.FileName)

		allItems, err := util.AddNewToDoItem(items, payload.Description)
		if err != nil {
			msg := "Failed to create new To-Do Item."
			http.Error(res, msg, http.StatusInternalServerError)
			util.LogError(msg)
			return
		}

		util.SaveAllToDoItems(allItems, util.FileName)
		util.LogInfo("Created new To-Do Item successfully")
		res.WriteHeader(http.StatusCreated)
	}
}

func getFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			msg := "Invalid Method. Accepted methods: GET"
			http.Error(res, msg, http.StatusMethodNotAllowed)
			util.LogError(msg)
			return
		}

		items, _ := util.GetAllToDoItems(util.FileName)
		res.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(res).Encode(items)
		if err != nil {
			util.LogError("Failed to encode To-Do Items.")
			return
		}

		util.LogInfo("Fetched To-Do Items.", "Items", items)
		res.WriteHeader(http.StatusOK)
	}
}

func updateFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			msg := "Invalid Method. Accepted methods: PUT"
			http.Error(res, msg, http.StatusMethodNotAllowed)
			util.LogError(msg)
			return
		}

		var payload struct {
			ItemId      int    `json:"id"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil || (payload.ItemId == 0 || (payload.Description == "" && payload.Status == "")) {
			msg := "Invalid request body. Accepted payload: \n" +
				"{\n" +
				"\"id\" : <Task Id>,\n" +
				"\"description\" : <Task Description>\n," +
				"\"status\" : <Task Status>\n}"
			http.Error(res, msg, http.StatusBadRequest)
			util.LogError(msg)
			return
		}

		items, _ := util.GetAllToDoItems(util.FileName)
		allItems, err := util.UpdateToDoItem(items, payload.ItemId, payload.Description, payload.Status)
		if err != nil {
			msg := "Failed to update To-Do Item."
			http.Error(res, msg, http.StatusInternalServerError)
			util.LogError(msg)
			return
		}

		util.SaveAllToDoItems(allItems, util.FileName)
		util.LogInfo("Updated To-Do Item successfully.", "Id", payload.ItemId)
		res.WriteHeader(http.StatusOK)
	}
}

func deleteFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			msg := "Invalid Method. Accepted methods: DELETE"
			http.Error(res, msg, http.StatusMethodNotAllowed)
			util.LogError(msg)
			return
		}

		idStr := req.URL.Query().Get("id")
		if idStr == "" {
			msg := "Missing 'id' query parameter."
			http.Error(res, msg, http.StatusBadRequest)
			util.LogError(msg)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			msg := "Invalid 'id' query parameter."
			http.Error(res, msg, http.StatusBadRequest)
			util.LogError(msg)
			return
		}

		items, _ := util.GetAllToDoItems(util.FileName)
		allItems, err := util.DeleteToDoItem(items, id)
		if err != nil {
			msg := "Failed to Delete To-Do Item."
			http.Error(res, msg, http.StatusInternalServerError)
			util.LogError(msg)
			return
		}

		util.SaveAllToDoItems(allItems, util.FileName)
		util.LogInfo("Deleted To-Do Item successfully.", "Id", id)
		res.WriteHeader(http.StatusOK)
	}
}

func listFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		tmpl, err := template.ParseFiles("webserver/template/list.html")
		if err != nil {
			msg := "Failed to load template."
			http.Error(res, msg, http.StatusInternalServerError)
			util.LogError(msg)
			return
		}

		items, _ := util.GetAllToDoItems(util.FileName)
		_ = tmpl.Execute(res, items)
	}
}
