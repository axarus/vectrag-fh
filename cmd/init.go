/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/axarus/vectrag/internal/application"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new VectraG project",
	Long: `Initialize a new VectraG project with the required directory structure
and configuration files.

This command will create:
  - vectrag.config.yaml (project configuration)
  - models/ (directory for model definitions)
  - config/(configuration files)
  - .vectrag/ (internal state directory)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}

		// Prompt for project name
		projectPrompt := promptui.Prompt{
			Label:   "Project name",
			Default: "my-app",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("project name cannot be empty")
				}

				if strings.Contains(input, " ") {
					return fmt.Errorf("project name cannot contain spaces")
				}
				return nil
			},
		}

		projectName, err := projectPrompt.Run()
		if err != nil {
			return err
		}

		// Prompt for port
		portPrompt := promptui.Prompt{
			Label:   "Port",
			Default: "3000",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("port cannot be empty")
				}
				port, err := strconv.Atoi(input)
				if err != nil {
					return fmt.Errorf("port must be a number")
				}
				if port < 1 || port > 65535 {
					return fmt.Errorf("port must be between 1 and 65535")
				}
				return nil
			},
		}
		port, err := portPrompt.Run()
		if err != nil {
			return err
		}

		// Prompt for database
		dbPrompt := promptui.Select{
			Label:     "Database",
			Items:     []string{"PostgreSQL", "MySQL", "SQLite"},
			CursorPos: 2,
		}

		_, database, err := dbPrompt.Run()
		if err != nil {
			return err
		}

		// Create init service and initialize project
		initService := application.NewInitService()
		config := application.InitConfig{
			ProjectName: projectName,
			Port:        port,
			Database:    database,
		}

		fmt.Printf("\nInitializing VectraG project at: %s\n", path)
		fmt.Printf("Project name: %s\n", projectName)
		fmt.Printf("Port: %s\n", port)
		fmt.Printf("Database: %s\n", database)

		if err := initService.InitializeProject(path, config); err != nil {
			return fmt.Errorf("failed to initialize project: %w", err)
		}

		projectPath := filepath.Join(path, projectName)

		fmt.Println("✅ Project initialized successfully!")
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  1. Navigate to the project: cd %s\n", projectPath)
		fmt.Printf("  2. Start development: vectrag develop\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
