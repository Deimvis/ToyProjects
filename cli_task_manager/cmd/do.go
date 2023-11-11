package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

var doCmd = &cobra.Command{
	Use:   "do [task numbers]",
	Short: "Mark a task completed",
	Long: fmt.Sprintf(`Specify numbers of completed tasks.
For example:
%s do 1 2 3`, db.APP_NAME),
	Run: func(cmd *cobra.Command, args []string) {
		taskNumbers := parseTaskNumbers(args)
		updateAvailableTasks(taskNumbers, func(t db.Task) error {
			return db.CompleteTask(t.Key)
		})
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
