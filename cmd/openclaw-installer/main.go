package main

import (
	"fmt"
	"os"

	"github.com/liuyingwen/openclaw-installer/internal/cli"
)

func main() {
	app := cli.NewApp(nil)
	if err := app.Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
