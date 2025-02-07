package todoCon

import (
	"os"
	"sync"
	"testing"
)

const tempFile = "test_ToDoData.json"

func TestToDoStore_Parallel(t *testing.T) {
	store, err := NewToDoStore(tempFile)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	var wg sync.WaitGroup
	t.Run("Add Items in Parallel", func(t *testing.T) {
		t.Parallel()
		wg.Add(2)

		go func() {
			defer wg.Done()
			if err := store.AddNewToDoItem("Task 1"); err != nil {
				t.Errorf("Failed to add Task 1: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			if err := store.AddNewToDoItem("Task 2"); err != nil {
				t.Errorf("Failed to add Task 2: %v", err)
			}
		}()

		wg.Wait()
	})

	t.Run("Update Items in Parallel", func(t *testing.T) {
		t.Parallel()
		wg.Add(2)

		go func() {
			defer wg.Done()
			if err := store.UpdateToDoItem(1, "started", "Updated Task 1"); err != nil {
				t.Errorf("Failed to update Task 1: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			if err := store.UpdateToDoItem(2, "completed", "Updated Task 2"); err != nil {
				t.Errorf("Failed to update Task 2: %v", err)
			}
		}()

		wg.Wait()
	})

	t.Run("Delete Items in Parallel", func(t *testing.T) {
		t.Parallel()
		wg.Add(2)

		go func() {
			defer wg.Done()
			if err := store.DeleteToDoItem(1); err != nil {
				t.Errorf("Failed to delete Task 1: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			if err := store.DeleteToDoItem(2); err != nil {
				t.Errorf("Failed to delete Task 2: %v", err)
			}
		}()

		wg.Wait()
	})

	t.Cleanup(func() {
		_ = os.Remove(tempFile)
	})
}
