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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "")
			os.Exit(0)
		}
		inputPath := args[0]

		// flag "password"
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}

		// flag "output"
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}

		err = app.Decrypt(
			inputPath,
			outputPath,
			[]byte(password),
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	decryptCmd.Flags().StringP("output", "o", "", "")
	decryptCmd.Flags().StringP("password", "p", "", "")
}
