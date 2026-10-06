package main

import (
	"fmt"
	// "os"
)


type Task struct {
	ID          int
	Description string
	Completed   bool
}


func main() {
	tasks := []Task{
		{ID: 1, Description: "Buy groceries", Completed: false},
		{ID: 2, Description: "Clean the house", Completed: true},
		{ID: 3, Description: "Finish the project", Completed: false},
	}
	for _, task := range tasks {
		fmt.Printf("Task %d: %s - %v\n", task.ID, task.Description, task.Completed)
	}
}
