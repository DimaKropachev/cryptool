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
	Short: "Decrypt previously encrypted files and directories.",
	Long: `The decrypt command decrypts files and directories that were previously encrypted using this utility.
It supports the same algorithms and parameters that were used during encryption and requires the correct password or key.
When decrypting directories, the original file structure is restored. The command allows selecting an output directory and controlling behavior in case of file conflicts.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "")
			os.Exit(0)
		}
		inFilePath := args[0]

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
			inFilePath,
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
