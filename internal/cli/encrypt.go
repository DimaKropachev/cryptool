/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DimaKropachev/cryptool/internal/app"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt files and directories using the selected encryption algorithm.",
	Long: `The encrypt command is used to encrypt individual files or entire directories.
It supports multiple encryption algorithms and modes of operation, allows specifying a password or key, and provides options to tune security and performance.
When encrypting directories, the original directory structure is preserved. The command can either overwrite the source data or create an encrypted copy in a specified output location.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "")
			os.Exit(0)
		}
		inFilePath := filepath.Clean(args[0])

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
		outputPath = filepath.Clean(outputPath)

		// flag "algorithm"
		alg, err := cmd.Flags().GetString("algorithm")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}

		err = app.Encrypt(
			alg,
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
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	encryptCmd.Flags().StringP("output", "o", "", "")
	encryptCmd.Flags().StringP("password", "p", "", "")
	encryptCmd.Flags().StringP("algorithm", "a", "aes256-gcm", "")
}
