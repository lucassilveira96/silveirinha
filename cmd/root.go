package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/lucassilveira96/silveirinha/services"
	"github.com/lucassilveira96/silveirinha/utils"
	"github.com/spf13/cobra"
)

// Define the version string
const version = "1.0.0"

// rootCmd is the main command
var rootCmd = &cobra.Command{
	Use:           "silverinha",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ShowBanner()
	},
	DisableAutoGenTag: true,
	Example: `
# To create a new Go project named 'my-new-project':
silverinha create my-new-project

# To create a new a model named 'modelExample':
silverinha model modelExample

# To see available commands and usage:
silverinha --help
`,
}

// "create" subcommand to create a new project
var createCmd = &cobra.Command{
	Use:           "create [project-name]",
	Aliases:       []string{"-c"},
	Short:         "Create a new Go project",
	Long:          `This command clones the base repository and creates a new project with the given name.`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if projectName == "" {
			fmt.Println("The project name cannot be empty.")
			return
		}

		err := services.CreateProject(projectName)
		if err != nil {
			log.Printf("Error creating project: %v", err)
		} else {
			fmt.Println("Project created successfully!")
		}
	},
	Example: `
# Create a new project with the name 'my-awesome-project':
silverinha create my-awesome-project
`,
}

// "model" subcommand to generate a model
var modelCmd = &cobra.Command{
	Use:           "model [model-name]",
	Aliases:       []string{"-m"},
	Short:         "Generate a new model",
	Long:          `This command generates a new Go model file with the specified name.`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		if modelName == "" {
			fmt.Println("The model name cannot be empty.")
			return
		}

		err := services.GenerateModel(modelName)
		if err != nil {
			log.Printf("Error generating model: %v", err)
		} else {
			fmt.Println("Model generated successfully!")
		}
	},
	Example: `
# Generate a model named 'User':
silverinha model User
`,
}

// Autocomplete subcommand to generate shell completion scripts
var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate shell autocompletion scripts",
	Long:  `This command generates shell autocompletion scripts for Bash, Zsh, or Fish shells.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must specify the shell (bash, zsh, fish).")
			return
		}

		shell := args[0]
		switch shell {
		case "bash":
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				log.Fatalf("Error generating bash completion: %v", err)
			}
		case "zsh":
			if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
				log.Fatalf("Error generating zsh completion: %v", err)
			}
		case "fish":
			if err := rootCmd.GenFishCompletion(os.Stdout, true); err != nil {
				log.Fatalf("Error generating fish completion: %v", err)
			}
		default:
			fmt.Println("Unsupported shell. Supported shells are: bash, zsh, fish.")
		}
	},
}

func init() {
	// Add the completion command to rootCmd
	rootCmd.AddCommand(completionCmd)
}

// Execute executes the root command
func Execute() error {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(modelCmd)
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show the version of Silverinha")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		rootCmd.Usage()
		return err
	}

	return nil
}
