/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"nenio/internal/objects"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files]",
	Short: "Add files to the staging area",
	Long: `The add command stages files to be included in the next commit by generating blobs and updating the index.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No files specified")
			os.Exit(1)
		}

		if err := addFiles(args); err != nil {
			fmt.Printf("failed to add files: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Files successfully added.")
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

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}

		hash, err := objects.CreateBlob(objectsDir, content) 
		if err != nil{
			return fmt.Errorf("failed to create blob for file %s: %w", file, err)
		}

		if err := objects.UpdateIndex(indexPath, hash, file); err != nil {
			return fmt.Errorf("failed to update index for file %s: %w", file, err)
		}
	}

	return nil
}