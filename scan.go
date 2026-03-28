package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/joho/godotenv"
)

func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	var path string
	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				// fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "node_modules" || file.Name() == "vendor" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.git-repositories"
	return dotFile
}

func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

func sliceContains(slice []string, element string) bool {
	for _, i := range slice {
		if i == element {
			return true
		}
	}
	return false
}

func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
	fmt.Println("File updated successfully")
}

func addNewSliceElementsToFile(filePath string, newElements []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	fmt.Println("REPOS", existingRepos)
	repos := joinSlices(existingRepos, newElements)
	dumpStringsSliceToFile(repos, filePath)
}

func scan(folder string) {
	fmt.Println("Scanning")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added %d repositories to the file\n", len(repositories))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
