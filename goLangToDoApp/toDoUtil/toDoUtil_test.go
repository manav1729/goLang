package toDoUtil

import (
	"os"
	"testing"
)

func TestToDoUtil(t *testing.T) {
	tempFile, err1 := os.CreateTemp("", "test_ToDoAppData.json")
	if err1 != nil {
		t.Error("Error creating temp file")
	}

	err2 := tempFile.Close()
	if err2 != nil {
		return
	}
}
