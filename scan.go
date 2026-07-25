package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

// Searches folder and sub folders for .git, appends them to slice.
// vendor and node modules folders are ignored.
func scanFolders(folders []string, folder string) []string {
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

	// Uses files slice and searches for .git, appends them to folders.
	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}

			// ignores vendor and node_modules folders
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanFolders(folders, path)
		}
	}
	return folders
}

// Starts the scanFolders function.
func recursivelyScanFolder(folder string) []string {
	return scanFolders(make([]string, 0), folder)
}

// Returns file path of the repo list.
// Will create a new file if its not found.
func getFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	file := usr.HomeDir + "/.gitlocalstats"

	return file
}

// Adds a slice to a given file.
func addSliceToFile(filepath string, newRepos []string) {
	existingRepos := parseFileToSlice(filepath)
	repos := joinSlices(newRepos, existingRepos)
	dumpSlicesToFile(repos, filepath)
}

// Parses contents of a file to a string slice.
func parseFileToSlice(filepath string) []string {
	f := openFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}

// Writes slice to given file.
func dumpSlicesToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	w := os.WriteFile(filepath, []byte(content), 0755)
	if w != nil {
		panic(w)
	}
}

// Opens file at filepath or creates it if it doesnt exist.
func openFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if os.IsExist(err) {
			_, err := os.Create(filepath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

// Joins slices and uses sliceContains to see if duplicate values
// are present, if not, new values are appended.
func joinSlices(new []string, old []string) []string {
	for _, i := range new {
		if !sliceContains(old, i) {
			old = append(old, i)
		}
	}
	return old
}

// Checks to see if slice contains value, returns a true or false
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
