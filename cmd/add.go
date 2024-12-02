/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"nenio/internal/objects"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files]",
	Short: "Add files to the staging area",
	Long: `The add command stages files to be included in the next commit by generating blobs and updating the index.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addFiles(files []string) error {
	var (
		objectsDir = ".nenio/objects"
		indexPath = ".nenio/index"
	)
}