package main

import (
	"fmt"
	"os"
	"encoding/json"

)


type Task struct {
	ID          int `json:"id"`
	Description string `json:"description"`
	Completed   bool `json:"completed"`
}


type TaskList struct {
	Tasks []Task `json:"tasks"`
}


func (t *TaskList) AddTask(description string) {
	newID := len(t.Tasks) + 1
	newTask := Task{
		ID:          newID,
		Description: description,
		Completed:   false,
	}
	t.Tasks = append(t.Tasks, newTask)
}


func (t *TaskList) GetTasks() []Task {
	return t.Tasks
}


func (t *TaskList) CompleteTask(id int) error {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}


func (t *TaskList) DeleteTask(id int) error {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)	
}


func (t *TaskList) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}


func main() {
	list := TaskList{
		Tasks: []Task{
			{ID: 1, Description: "Buy groceries", Completed: false},
			{ID: 2, Description: "Clean the house", Completed: true},
			{ID: 3, Description: "Finish the project", Completed: false},
		},
	}

	list.AddTask("Read a book")
	list.AddTask("Go for a walk") 

	err := list.CompleteTask(1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(list.GetTasks())

	err = list.SaveToFile("tasks.json")

	if err != nil {
		fmt.Println("Error saving tasks to file:", err)
	} else {
		fmt.Println("Tasks saved to tasks.json")
	}

	for _, task := range list.Tasks {
		fmt.Printf("Task %d: %s - %v\n", task.ID, task.Description, task.Completed)
	}

	 
	fmt.Println(list.GetTasks())

	for _, task := range list.Tasks {
		fmt.Printf("Task %d: %s - %v\n", task.ID, task.Description, task.Completed)
	}
}
