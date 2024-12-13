/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Usage: nenio init",
	Long: `Initializes an empty local repository in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := InitializeRepo(); err != nil {
			fmt.Printf("Failed to initialize repository: %v", err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty nenio repository in ./.nenio")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitializeRepo() error {
// 	/.nenio/
// ├── config          # Файл конфігурації репозиторію
// |-- index 		   # Файл, який зберігає метадані про файли, додані до стейджинг-зони (індексу).
// ├── HEAD            # Поточна гілка або коміт
// ├── objects/        # Директорія для збереження блобів, дерев, комітів
// └── refs/
//     ├── heads/      # Гілки
//     └── tags/       # Теги

var (
	basePath = ".nenio"
	subDirs = []string{
		"objects",
		"refs/heads",
		"refs/tags",
	}
	files = map[string]string{
		"HEAD": "ref: refs/heads/main\n", // гілка по дефолту main
		"config": "[core]\n\trepositoryFormatVersion =0\n",
		"index.json": "",
	}
)

if _, err := os.Stat(basePath); !os.IsNotExist(err) {
	return fmt.Errorf(".nenio already exists")
}

for _, dir := range subDirs {
	path := filepath.Join(basePath, dir)
	if err := os.MkdirAll(path, 0755); err != nil  {
		return err
	}
}

for file, content := range files {
	filePath := filepath.Join(basePath, file)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}
}


return nil
}
