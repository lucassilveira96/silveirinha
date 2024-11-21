package services

import (
	"fmt"

	"github.com/lucassilveira96/silveirinha/utils"
)

// CreateProject creates a new Go project from a template by cloning the repository
func CreateProject(projectName string) error {
	// Check if the project name is valid (e.g., no spaces)
	if !utils.IsValidProjectName(projectName) {
		return fmt.Errorf("invalid project name: %s", projectName)
	}

	// Clone the repository
	err := utils.CloneRepository(projectName)
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	// Remove the .git directory to make the project independent
	err = utils.RemoveGitDirectory(projectName)
	if err != nil {
		return fmt.Errorf("error removing .git directory: %v", err)
	}

	// Replace project name in go.mod and source files
	err = utils.ReplacePackagesNames(projectName)
	if err != nil {
		return fmt.Errorf("error replacing package names: %v", err)
	}

	return nil
}
