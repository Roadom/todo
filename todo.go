package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Task struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type Todo struct {
	Tasks []Task `json:"tasks"`
}

// Load list from file
func LoadTodo(filename string) (Todo, error) {
	var todo Todo

	data, err := os.ReadFile(filename)
	if err != nil {
		return todo, err
	}

	err = json.Unmarshal(data, &todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

// Save JSON (used when modifying data through operations)
func (todo *Todo) Save(filename string) error {
	data, err := json.MarshalIndent(todo, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Todo List Operations
// Insert new task
func (todo *Todo) AddTask(newTask Task) error {
	todo.Tasks = append(todo.Tasks, newTask)

	// Save JSON
	if err := todo.Save("todos.json"); err != nil {
		log.Println("save failed:", err)
	}
	return nil
}

// Toggle Completion
func (todo *Todo) ToggleCompletion(index int) error {
	if !todo.IsValidIndex(index) {
		return fmt.Errorf("invalid index: %d", index)
	}

	todo.Tasks[index].Completed = !todo.Tasks[index].Completed

	// Save JSON
	if err := todo.Save("todos.json"); err != nil {
		log.Println("save failed:", err)
	}
	return nil
}

// Rename Task
func (todo *Todo) RenameTask(index int, newName string) error {
	if !todo.IsValidIndex(index) {
		return fmt.Errorf("invalid index: %d", index)
	}

	newName = strings.TrimSpace(newName)

	// Reject empty name
	if newName == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	// Reject if same name
	if newName == todo.Tasks[index].Name {
		return fmt.Errorf("new name is the same as current name")
	}

	todo.Tasks[index].Name = newName

	// Save JSON
	if err := todo.Save("todos.json"); err != nil {
		log.Println("save failed:", err)
	}

	return nil
}

// Delete task
func (todo *Todo) DeleteTask(index int) error {
	if !todo.IsValidIndex(index) {
		return fmt.Errorf("invalid index: %d", index)
	}

	todo.Tasks = append(todo.Tasks[:index], todo.Tasks[index+1:]...)

	// Save JSON
	if err := todo.Save("todos.json"); err != nil {
		log.Println("save failed:", err)
	}

	return nil
}

// Helper function for checking index
func (todo *Todo) IsValidIndex(index int) bool {
	return index >= 0 && index < len(todo.Tasks)
}
