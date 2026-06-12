package main

import (
	"log"
	"strings"

	"charm.land/huh/v2"
)

func main() {
	// load data from file
	todos, err := LoadTodo("todos.json")
	if err != nil {
		log.Println("could not load file, starting empty:", err)
		todos = Todo{}
	}

	// create selection of tasks (+1 for the add task button)
	for {
		options := make([]huh.Option[int], len(todos.Tasks)+1)
		for index, task := range todos.Tasks {
			label := "[ ] "
			if task.Completed {
				label = "[x] "
			}
			options[index] = huh.NewOption(label+task.Name, index)
		}

		// add "Add Task"
		options[len(todos.Tasks)] = huh.NewOption(" + Add Task", len(todos.Tasks))

		var selectedTask int

		// Form 1: Display Tasks
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[int]().
					Title("Todo List").
					Options(options...).
					Value(&selectedTask),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if selectedTask == len(todos.Tasks) {
			// Form: Add Task
			var newTask string

			addForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("New Task").
						Value(&newTask),
				),
			)

			if err := addForm.Run(); err != nil {
				log.Fatal(err)
			}

			if strings.TrimSpace(newTask) == "" {
				return
			}

			err := todos.AddTask(Task{
				Name:      newTask,
				Completed: false,
			})

			if err != nil {
				log.Println("failed to add task:", err)
				continue
			}

			continue
		} else {
			// Form 2: Action Selection
			var action string

			actionForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Choose Action").
						Options(
							huh.NewOption("Toggle Completion", "toggle"),
							huh.NewOption("Rename Task", "rename"),
							huh.NewOption("Delete Task", "delete"),
							huh.NewOption("Back", "back"),
						).
						Value(&action),
				),
			)

			if err := actionForm.Run(); err != nil {
				log.Fatal(err)
			}

			switch action {
			case "toggle":
				// Toggle completion status
				err := todos.ToggleCompletion(selectedTask)
				if err != nil {
					log.Println("failed to toggle task:", err)
					continue
				}
			case "rename":
				// FORM 3: Rename
				var newName string

				renameForm := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("New Task Name").
							Value(&newName),
					),
				)

				if err := renameForm.Run(); err != nil {
					log.Fatal(err)
				}

				err := todos.RenameTask(selectedTask, newName)
				if err != nil {
					log.Println("failed to rename task:", err)
					continue
				}

			case "delete":
				// FORM 3: Confirm delete
				var confirm bool

				deleteForm := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Delete this task?").
							Value(&confirm),
					),
				)

				if err := deleteForm.Run(); err != nil {
					log.Fatal(err)
				}

				if confirm {
					// delete task
					err := todos.DeleteTask(selectedTask)
					if err != nil {
						log.Println("failed to delete task:", err)
						continue
					}
				}

			case "back":
				return
			}
		}
	}
}
