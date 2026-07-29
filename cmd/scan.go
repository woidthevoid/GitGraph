package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a folder for git repositories",
	Long: `Recursively walks the given folder and searches for git repos,
saves them so stats can process them later`,
	Run: func(cmd *cobra.Command, args []string) {
		folder, err := cmd.Flags().GetString("folder")
		if err != nil {
			log.Fatal(err)
		}

		folders := RecursivelyScanFolder(folder)
		filepath := GetFilePath()
		AddSliceToFile(filepath, folders)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringP("folder", "f", "", "folder to scan")
	scanCmd.MarkFlagRequired("folder")
}

// Searches folder and sub folders for .git, appends them to slice.
// vendor and node modules folders are ignored.
func ScanFolders(folders []string, folder string) []string {
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
			folders = ScanFolders(folders, path)
		}
	}
	return folders
}

// Starts the scanFolders function.
func RecursivelyScanFolder(folder string) []string {
	return ScanFolders(make([]string, 0), folder)
}

// Returns file path of the repo list.
// Will create a new file if its not found.
func GetFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	file := usr.HomeDir + "/.gitlocalstats"

	return file
}

// Adds a slice to a given file.
func AddSliceToFile(filepath string, newRepos []string) {
	existingRepos := ParseFileToSlice(filepath)
	repos := JoinSlices(newRepos, existingRepos)
	DumpSlicesToFile(repos, filepath)
}

// Parses contents of a file to a string slice.
func ParseFileToSlice(filepath string) []string {
	f := OpenFile(filepath)
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
func DumpSlicesToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	w := os.WriteFile(filepath, []byte(content), 0755)
	if w != nil {
		panic(w)
	}
}

// Opens file at filepath or creates it if it doesnt exist.
func OpenFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

// Joins slices and uses SliceContains to see if duplicate values
// are present, if not, new values are appended.
func JoinSlices(new []string, old []string) []string {
	for _, i := range new {
		if !SliceContains(old, i) {
			old = append(old, i)
		}
	}
	return old
}

// Checks to see if slice contains value, returns a true or false
func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
