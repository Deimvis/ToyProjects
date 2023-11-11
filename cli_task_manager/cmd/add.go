package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		taskText := strings.Join(args, " ")
		db.CreateTask(taskText)
		fmt.Printf("New task \"%s\" was created", taskText)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
