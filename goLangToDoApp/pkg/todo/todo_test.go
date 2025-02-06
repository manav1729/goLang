package todo

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

const tempFile = "test_ToDoData.json"

func TestToDo(t *testing.T) {
	testParallelCreateRead(t)
}

func testParallelCreateRead(t *testing.T) {
	store, err := NewToDoStore(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize ToDoStore: %s", err)
	}

	var wg sync.WaitGroup
	numWorkers := 10

	// Run concurrent writes
	t.Run("Concurrent Writes", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				store.AddNewToDoItem("Task " + strconv.Itoa(i))
			}(i)
		}
		wg.Wait()
	})

	// Run concurrent reads
	t.Run("Concurrent Reads", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = store.GetAllToDoItems()
			}()
		}
		wg.Wait()
	})

	var items []Item
	err = loadAllToDoItems(tempFile, &items)
	if err != nil {
		t.Fatalf("Failed to load all ToDoItems: %s", err)
	}
	if len(items) != numWorkers {
		t.Fatalf("Expected %d items, got %d", numWorkers, len(items))
	}

	// Clean up test file
	os.Remove(tempFile)
}
