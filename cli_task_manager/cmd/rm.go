/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task (that is not completed yet)",
	Run: func(cmd *cobra.Command, args []string) {
		taskNumbers := parseTaskNumbers(args)
		updateAvailableTasks(taskNumbers, func(t db.Task) error {
			return db.RemoveTask(t.Key)
		})
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
