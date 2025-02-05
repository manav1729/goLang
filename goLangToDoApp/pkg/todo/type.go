package todo

type Item struct {
	ItemId      int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ToDoStore struct {
	FilePath string
	Items    []Item
}
