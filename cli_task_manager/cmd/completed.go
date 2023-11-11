/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List completed tasks",
	Run: func(cmd *cobra.Command, args []string) {
		n, err := cmd.Flags().GetInt("last")
		if err != nil {
			fmt.Println("Something went wrong: ", err.Error())
			os.Exit(1)
		}
		n = max(n, 0)
		tasks, err := db.ListArchivedTasks(db.CompletedFilter)
		if err != nil {
			fmt.Println("Something went wrong: ", err.Error())
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks were completed yet")
			return
		}
		if len(tasks) > n {
			tasks = tasks[len(tasks)-n:]
		}
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CompleteTs > tasks[j].CompleteTs
		})
		printTasks(tasks)
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
	completedCmd.Flags().IntP("last", "n", 10, "Number of most recently completed tasks to print")
}
