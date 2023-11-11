package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

// Returns not completed, not removed tasks and preserves the order
func listAvailableTasks() ([]db.Task, error) {
	tasks, err := db.ListTasks(db.NotCompletedFilter, db.NotRemovedFilter)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreateTs < tasks[j].CreateTs
	})
	return tasks, nil
}

func printTasks(tasks []db.Task) {
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task.Text)
	}
}

func updateAvailableTasks(taskNumbers []int, fn func(db.Task) error) {
	tasks, err := listAvailableTasks()
	if err != nil {
		fmt.Println("Something went wrong: ", err.Error())
		os.Exit(1)
	}
	for _, number := range taskNumbers {
		ind := number - 1
		if !(0 <= ind && ind <= len(tasks)-1) {
			fmt.Println("Got invalid task number:", number)
			continue
		}
		err = fn(tasks[ind])
		if err != nil {
			fmt.Printf("Failed to mark \"%d\" as completed\nError: %s\n", number, err)
		}
	}
}

func parseTaskNumbers(args []string) []int {
	var taskNumbers []int
	for _, arg := range args {
		number, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Failed to parse the argument:", arg)
			continue
		}
		taskNumbers = append(taskNumbers, number)
	}
	return taskNumbers
}
