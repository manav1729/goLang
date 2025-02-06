package todo

type Item struct {
	ItemId      int    `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ToDoStore struct {
	filePath string
	items    []Item
}
