/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"nenio/internal/objects"
	"os"
	"path/filepath"

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
	objectsDir := "./.nenio/objects" 
	indexPath := "./.nenio/index.json"
	ignorePath := "./.nignore"

	if err := ensureNenioStructure(objectsDir, indexPath); err != nil {
		return fmt.Errorf("failed to prepare nenio structure: %v", err)
	}

	var allFiles []string
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %v", file, err)
		}

		if info.IsDir() {
			filesInDir, err := GetFilesInDir(file)
			if err != nil {
				return fmt.Errorf("failed to list files in directory %s: %v", file, err)
			}
			allFiles = append(allFiles, filesInDir...)
		} else {
			allFiles = append(allFiles, file)
		}
	}

	ignorePatterns, _ := objects.LoadIgnorePatterns(ignorePath)

	err := objects.AddToIndex(objectsDir, indexPath, ignorePatterns, allFiles)
	if err != nil {
		return err
	}

	return nil
}

//Change


func ensureNenioStructure(objectsDir, indexPath string) error {
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(objectsDir, 0755); err != nil {
			return fmt.Errorf("failed to create objects directory: %v", err)
		}
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		initialIndex := objects.NewIndex()
		if err := objects.SaveIndex(indexPath, initialIndex); err != nil {
			return fmt.Errorf("failed to create initial index: %v", err)
		}
	}

	return nil
}

func GetFilesInDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}