package main

import (
	"flag"
	"os"
)

func main() {
	loadEnv()
	folder := os.Getenv("GIT_FOLDER")
	email := os.Getenv("EMAIL")

	addFolderFlag := flag.Bool("add", false, "add a new folder to scan for git repositories")

	flag.Parse()
	if *addFolderFlag != false {
		scan(folder)
		return
	}
	stats(email)
}
