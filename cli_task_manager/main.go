package main

import (
	"fmt"
	"os"

	"github.com/Deimvis/toyprojects/cli_task_manager/cmd"
	db "github.com/Deimvis/toyprojects/cli_task_manager/src"
)

func main() {
	check(db.InitDB())
	check(cmd.Execute())
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		db.StopDB()
		os.Exit(1)
	}
}
