package cmd

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/spf13/cobra"
)

const daysSixMonths = 183
const outOfRange = 9999
const weeksSixMonths = 26

type column []int

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Shows git contribution graph",
	Long:  `Builds and prints a git contribution graph based on the last 6 months activity`,
	Run: func(cmd *cobra.Command, args []string) {
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal(err)
		}

		commits := ProcessRepos(email)
		PrintStats(commits)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().StringP("email", "e", "", "the email of the git author")
	statsCmd.MarkFlagRequired("email")
}

// ProcessRepos takes email as parameter and returns a commit map for the last sinx months.
// Gets the git file path, parses file to a list of repos, fills a commit map with 0s
// and fills the map with commits while iterating over repos.
func ProcessRepos(email string) map[int]int {
	filepath := GetFilePath()
	repos := ParseFileToSlice(filepath)
	days := daysSixMonths

	commits := make(map[int]int, days)
	for i := days; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = FillCommits(email, path, commits)
	}
	return commits
}

// FillCommits given a valid git repository, commits made my given email will be counted and returned as a map
func FillCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := CalcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := CountDaysSince(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}

		if daysAgo != outOfRange {
			commits[daysAgo]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	repo.Close()
	return commits
}

// CalcOffset calculates a offset based on current day, sets offset to value between 1-7 and returns offset.
func CalcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Saturday:
		offset = 1
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	}
	return offset
}

func CountDaysSince(date time.Time) int {
	now := GetBeginningOfDay(time.Now())
	commitDate := GetBeginningOfDay(date)
	if commitDate.After(now) {
		return 0
	}
	days := int(now.Sub(commitDate).Hours() / 24)

	if days > daysSixMonths {
		return outOfRange
	}
	return days
}

// PrintStats builds git commit graph.
func PrintStats(commits map[int]int) {
	keys := SortIntoSlice(commits)
	cols := BuildCol(keys, commits)
	PrintCells(cols)
}

// SortIntoSlice takes map and sorts it into slice, based on integers.
func SortIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// BuildCol bulds columns based on commits given and puts them on their respective day of the week.
func BuildCol(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := k / 7
		dayinweek := k % 7

		if dayinweek == 0 {
			col = column{}
		}

		col = append(col, commits[k])
		if dayinweek == 6 {
			cols[week] = col
		}
	}
	return cols
}

// PrintCells is used to print out all the different cells.
func PrintCells(cols map[int]column) {
	PrintMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksSixMonths + 1; i >= 0; i-- {
			if i == weeksSixMonths+1 {
				PrintDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == CalcOffset()-1 {
					PrintCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						PrintCell(col[j], false)
						continue
					}
				}
			}
			PrintCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func PrintMonths() {
	week := GetBeginningOfDay(time.Now()).Add(-(daysSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func PrintDayCol(day int) {
	out := "   "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri"
	}
	fmt.Printf(out)
}

// PrintCell colors the different cells with a color based on amount of commits, todays date
// will come up as purple.
func PrintCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// GetBeginningOfDay gets the time a day begins in a specific timezone and returns it.
func GetBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}
