package toDoUtil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ToDoItem struct {
	ItemId      int    `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

func PrintFlagInstructions() {
	fmt.Println("======================== Use following flags for various operations =======================")
	fmt.Println("-add -header=<name> -desc <description> to \"Add a new To-Do Item\"")
	fmt.Println("-update -id=<itemId> -header=<name> -desc <description> to \"Update a To-Do Item\"")
	fmt.Println("-remove -id=<itemId> to \"Delete a To-Do Item\"")
	fmt.Println("-removeAll to \"Delete all To-Do Item(s)\"")
	fmt.Println("===========================================================================================")
}

func PrintToDoItems(items []ToDoItem) {
	if items != nil && len(items) > 0 {
		fmt.Println("================================== Your To-Do Task Items ==================================")
		//fmt.Println(items)
		for index, item := range items {
			if index != 0 {
				fmt.Println("-------------------------------------------------------------------------------------------")
			}
			fmt.Printf("%d. %s\n", item.ItemId, item.Header)
			fmt.Println(item.Description)
		}
		fmt.Println("===========================================================================================")
	} else {
		fmt.Println("============================== No To-Do Items in the List =================================")
	}
}

func AddNewToDoItem(currentItems []ToDoItem, header string, desc string) []ToDoItem {
	id := 1
	itemNos := len(currentItems)
	if itemNos > 0 {
		id = currentItems[itemNos-1].ItemId + 1
	}
	return append(currentItems, ToDoItem{id, header, desc})
}

func UpdateToDoItem(currentItems []ToDoItem, id int, header string, desc string) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == id {
			currentItems[index].Header = header
			currentItems[index].Description = desc
			success = true
		}
	}
	if success {
		fmt.Println("================================== To-Do Item updated ====================================")
	} else {
		fmt.Println("================================ To-Do Item not updated ==================================")
	}
	return currentItems
}

func RemoveToDoItem(currentItems []ToDoItem, deleteItemId int) []ToDoItem {
	success := false
	for index, item := range currentItems {
		if item.ItemId == deleteItemId {
			currentItems = append(currentItems[:index], currentItems[index+1:]...)
			success = true
		}
	}
	if success {
		fmt.Println("================================== To-Do Item deleted ====================================")
	} else {
		fmt.Println("================================ To-Do Item not deleted ==================================")
	}
	return currentItems
}

func SaveAllToDoItems(allItems []ToDoItem, fileName string) error {
	// Open json file
	data, err := json.MarshalIndent(allItems, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func GetAllToDoItems(fileName string) []ToDoItem {
	// Open local json file
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("==========> Error Opening file:", fileName, err)
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println("==========> Error Reading file:", fileName, err)
		} else {
			if byteValue != nil {
				// Parse the json file to ToDoItems
				var items []ToDoItem
				err := json.Unmarshal(byteValue, &items)
				if err != nil {
					fmt.Println("==========> Error Unmarshalling file:", err)
				} else {
					if items != nil && len(items) > 0 {
						return items
					}
				}
			}
		}
	}
	return nil
}
