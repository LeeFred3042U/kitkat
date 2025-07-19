package main

import (
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/core"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kitkat <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		err := core.InitRepo()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty kitkat repository in .kitkat/")
	
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("usage: kitkat add <file>")
			os.Exit(1)
		}
		err := core.AddFile(os.Args[2])
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("added %s\n", os.Args[2])

	case "ls-files":
		entries, err := storage.LoadIndex()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		for path := range entries {
			fmt.Println(path)
	    }

	default:
		fmt.Println("unknown command:", os.Args[1])
	}
}
