/*
Copyright © 2026 Victor Woydowski Dralle <woydowskidralle@proton.me>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitgraph",
	Short: "A local git contribution graph",
	Long: `GitGraph is a local git contribution graph maker, similar to the one found
on GitHub. GitGraph searches through your given path for .git folders and gives you a graoh based on your last
6 months contributions.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
