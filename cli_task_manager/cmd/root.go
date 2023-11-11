/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

var rootCmd = &cobra.Command{
	Use:   db.APP_NAME,
	Short: "CLI task manager",
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}
