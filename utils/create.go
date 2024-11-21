package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CloneRepository clones the repository into the specified project name
func CloneRepository(projectName string) error {
	cmd := exec.Command("git", "clone", "https://github.com/lucassilveira96/template-go-with-silverinha-file-genarator", projectName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}
	return nil
}

// RemoveGitDirectory removes the .git directory from the project to untrack it
func RemoveGitDirectory(projectName string) error {
	gitDir := filepath.Join(projectName, ".git")

	// Cross-platform removal
	err := os.RemoveAll(gitDir)
	if err != nil {
		return fmt.Errorf("error removing .git directory: %v", err)
	}
	return nil
}

// IsValidProjectName checks if the project name is valid (no spaces)
func IsValidProjectName(name string) bool {
	return strings.TrimSpace(name) != "" && !strings.Contains(name, " ")
}

// ReplacePackagesNames updates the project name in go.mod and source files
func ReplacePackagesNames(projectName string) error {
	// Path of the cloned project
	projectPath := fmt.Sprintf("./%s", projectName)

	// Update go.mod
	goModFile := filepath.Join(projectPath, "go.mod")
	err := ReplaceGoMod(goModFile, projectName)
	if err != nil {
		return fmt.Errorf("error updating go.mod: %v", err)
	}

	// Update package names in source files
	err = UpdateGoFiles(projectPath, projectName)
	if err != nil {
		return fmt.Errorf("error updating package names in .go files: %v", err)
	}

	return nil
}

// UpdateGoFiles updates package names in all Go files
func UpdateGoFiles(path, projectName string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only .go files
		if !info.IsDir() && strings.HasSuffix(filePath, ".go") {
			err = ReplaceTextInFile(filePath, "template-go-with-silverinha-file-genarator", projectName)
			if err != nil {
				return fmt.Errorf("error updating file %s: %v", filePath, err)
			}
		}
		return nil
	})
	return err
}

// ReplaceTextInFile replaces all occurrences of oldText with newText in the given file
func ReplaceTextInFile(filePath, oldText, newText string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	content := string(data)
	newContent := strings.ReplaceAll(content, oldText, newText)

	if newContent != content {
		err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", filePath, err)
		}
	}
	return nil
}

// ReplaceGoMod replaces the module name in go.mod with the new project name
func ReplaceGoMod(goModFile string, projectName string) error {
	data, err := ioutil.ReadFile(goModFile)
	if err != nil {
		return fmt.Errorf("error reading go.mod file: %v", err)
	}

	content := string(data)
	newContent := strings.Replace(content, "module template-go-with-silverinha-file-genarator", fmt.Sprintf("module %s", projectName), 1)

	if newContent == content {
		return fmt.Errorf("module name not found in go.mod")
	}

	err = ioutil.WriteFile(goModFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to go.mod file: %v", err)
	}

	return nil
}
