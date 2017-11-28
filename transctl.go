package main

import (
	"dtom/pkg/cmd"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need input swagger.json path")
		os.Exit(0)
	}
	fmt.Println(os.Args[1])
	sa := cmd.SwaggerAction{}
	sa.ShowConf(os.Args[1])
}
