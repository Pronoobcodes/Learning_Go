package main

import (
	"fmt"
	"os"
	"encoding/json"
	"errors"
	"flag"
)

var ErrTaskNotFound = errors.New("task not found") 


type TaskRepository interface{
	AddTask(description string)
	GetTasks() []Task
	CompleteTask(id int) error
	DeleteTask(id int) error
}


type Task struct {
	ID          int `json:"id"`
	Description string `json:"description"`
	Completed   bool `json:"completed"`
}


type TaskList struct {
    Tasks []Task `json:"tasks"`

    taskIndex map[int]int
    nextID    int
}


func (t *TaskList) rebuildIndex() {
    t.taskIndex = make(map[int]int)

    t.nextID = 1

    for i, task := range t.Tasks {
        t.taskIndex[task.ID] = i

        if task.ID >= t.nextID {
            t.nextID = task.ID + 1
        }
    }
}


func (t *TaskList) AddTask(description string) {
	if t.taskIndex == nil {
        t.rebuildIndex()
    }

    newTask := Task{
        ID:          t.nextID,
        Description: description,
        Completed:   false,
    }

    t.Tasks = append(t.Tasks, newTask)

    t.taskIndex[newTask.ID] = len(t.Tasks) - 1
    t.nextID++
}


func (t *TaskList) GetTasks() []Task {
	return t.Tasks
}


func (t *TaskList) CompleteTask(id int) error {
	if t.taskIndex == nil {
        t.rebuildIndex()
    }

    index, exists := t.taskIndex[id]
    if !exists {
        return ErrTaskNotFound
    }

    t.Tasks[index].Completed = true

    return nil
}


func (t *TaskList) DeleteTask(id int) error {
	if t.taskIndex == nil {
        t.rebuildIndex()
    }

    index, exists := t.taskIndex[id]
    if !exists {
        return ErrTaskNotFound
    }

    t.Tasks = append(t.Tasks[:index], t.Tasks[index+1:]...)

    t.rebuildIndex()

    return nil
}


func (t *TaskList) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}


func (t *TaskList) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)

    if err != nil {
        if os.IsNotExist(err) {
            t.rebuildIndex()
            return nil
        }

        return err
    }

    if err := json.Unmarshal(data, t); err != nil {
        return err
    }

    t.rebuildIndex()

    return nil
}


func printTasks(repository TaskRepository) {
	for _, task := range repository.GetTasks() {
		fmt.Printf(
			"Task %d: %s - %v\n",
			task.ID,
			task.Description,
			task.Completed,
		)
	}
}


func main() {
	/*
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
	*/

	/*
	store := TaskList{}	

	err := store.LoadFromFile("todo-cli/tasks.json")
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
		return
	}

	fmt.Println("Loaded tasks from tasks.json:")
	fmt.Println(store.GetTasks())
	*/

	/*
	store := &TaskList{}

	store.AddTask("Learn Go")
	store.AddTask("Build Todo CLI")

	printTasks(store)
	*/

	command := flag.String("command", "", "command to execute (add, list, done, delete)")
	description := flag.String("description", "", "task description (required for 'add')")
	id := flag.Int("id", 0, "task ID (required for 'done' and 'delete')")
	flag.Parse()

	store := &TaskList{}
	err := store.LoadFromFile("tasks.json")
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	switch *command {
		case "add":
			if *description == "" {
				fmt.Println("Error: Description is required for adding a task.")
				return
			}
			store.AddTask(*description)
			err := store.SaveToFile("tasks.json")
			if err != nil {
				fmt.Println("Error saving task:", err)
				return
			}
			fmt.Println("Task added successfully.")

		case "list":
			tasks := store.GetTasks()
			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return
			}

			for _, task := range tasks {
				status := " "
				if task.Completed {
					status = "✔️"
				}

				fmt.Printf("[%s  ] %d: %s\n", status, task.ID, task.Description)
			}

		case "done":
			if *id <= 0 {
				fmt.Println("Error: A valid task ID is required.")
				return
			}

			err := store.CompleteTask(*id)

			if errors.Is(err, ErrTaskNotFound) {
				fmt.Printf("Task %d not found.\n", *id)
				return
			}

			err = store.SaveToFile("tasks.json")

			if err != nil {
				fmt.Println("Error saving changes:", err)
				return
			}
			fmt.Printf("Task %d marked as complete.\n", *id)

		case "delete":
			if *id <= 0 {
				fmt.Println("Error: A valid task ID is required.")
				return
			}

			err := store.DeleteTask(*id)
			
			if errors.Is(err, ErrTaskNotFound) {
				fmt.Printf("Task %d not found.\n", *id)
				return
			}

			err = store.SaveToFile("tasks.json")
			
			if err != nil {
				fmt.Println("Error saving changes:", err)
				return
			}
			fmt.Printf("Task %d deleted successfully.\n", *id)

		default:
			fmt.Println("Unknown or missing command. Use -command with add, list, done, or delete.")
	}

	/*

	go run main.go -command list

	
	go run main.go -command delete -id 2

	
	go run main.go -command add -description "Test ID system"

	
	go run main.go -command list
	 
	*/
}
