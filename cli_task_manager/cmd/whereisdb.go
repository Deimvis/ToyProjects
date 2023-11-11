/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

var whereisdbCmd = &cobra.Command{
	Use:   "whereisdb",
	Short: "Prints the db file path",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := db.GetDBPath()
		if err != nil {
			fmt.Println("Something went wrong: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(dbPath)
	},
}

func init() {
	rootCmd.AddCommand(whereisdbCmd)
}
