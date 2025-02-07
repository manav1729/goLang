package todoCon

type Item struct {
	ItemId      int    `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type request struct {
	action string
	item   Item
	id     int
	status string
	desc   string
	resp   chan error
}

type ToDoStore struct {
	filePath string
	items    []Item
	requests chan request
}
