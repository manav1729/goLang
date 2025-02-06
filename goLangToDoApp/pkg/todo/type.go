package todo

type Item struct {
	ItemId      int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ToDoStore struct {
	filePath string
	commands chan func(*[]Item, *int)
}
