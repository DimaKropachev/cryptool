/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/DimaKropachev/cryptool/internal/app"
	"github.com/spf13/cobra"
)

// benchmarkCmd represents the benchmark command
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run an encryption benchmark on a selected file.",
	Long: `The benchmark command compares encryption algorithms based on their speed and memory usage using a specific file.
During execution, the command applies supported encryption algorithms sequentially and measures execution time as well as peak and average memory consumption.
The results are presented in an easy-to-analyze format and can be used to choose the most suitable encryption algorithm for specific performance and resource requirements.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "you must specify the path of the input file/dir")
			os.Exit(0)
		}
		inFilePath := args[0]

		err := app.Benchmark(inFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benchmarkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchmarkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
