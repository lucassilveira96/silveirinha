package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateService generates Go service files for a given model name.
func GenerateService(modelName, structName string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	// Get the current folder name
	currentFolderName := filepath.Base(currentDir)

	// Construct the service directory path
	serviceDir := filepath.Join("internal", "app", "domain", "service", modelName)

	// Ensure the service directory exists
	if err := os.MkdirAll(serviceDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating service directory: %v", err)
	}

	// Generate the Service interface file
	serviceFilePath := filepath.Join(serviceDir, fmt.Sprintf("%sService.go", modelName))
	if err := writeServiceInterfaceFile(serviceFilePath, currentFolderName, modelName, structName); err != nil {
		return fmt.Errorf("error writing service interface file: %v", err)
	}

	// Generate the Service implementation file
	serviceImplFilePath := filepath.Join(serviceDir, fmt.Sprintf("%sServiceImpl.go", modelName))
	if err := writeServiceImplFile(serviceImplFilePath, currentFolderName, modelName, structName); err != nil {
		return fmt.Errorf("error writing service implementation file: %v", err)
	}

	// Edit the services.go file
	servicesFile := filepath.Join("internal", "app", "domain", "services.go")
	if err := editServicesFile(servicesFile, modelName, structName, currentFolderName); err != nil {
		return fmt.Errorf("error editing services.go: %v", err)
	}

	fmt.Printf("Service files generated in: %s\n", serviceDir)
	return nil
}

// writeServiceInterfaceFile creates the service interface file.
func writeServiceInterfaceFile(filePath, currentFolderName, modelName, structName string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the content
	content := fmt.Sprintf(`package %sService

import "%s/internal/app/domain/model"

type %sService interface {
	Create(%s *model.%s) error
	Update(id uint, %s *model.%s) error
	Delete(id uint) error
	FindAll() ([]*model.%s, error)
	FindById(id uint) (*model.%s, error)
}
`, modelName, currentFolderName, structName, modelName, structName, modelName, structName, structName, structName)

	writer.WriteString(content)
	writer.Flush()
	return nil
}

// writeServiceImplFile creates the service implementation file.
func writeServiceImplFile(filePath, currentFolderName, modelName, structName string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the content
	content := fmt.Sprintf(`package %sService

import (
	"%s/internal/app/domain/model"
	%sRepository "%s/internal/app/domain/repository/%s"
)

var _ %sService = (*%sServiceImpl)(nil)

type %sServiceImpl struct {
	repository %sRepository.%sRepository
}

func New%sService(repository %sRepository.%sRepository) *%sServiceImpl {
	return &%sServiceImpl{repository: repository}
}

func (s *%sServiceImpl) Create(%s *model.%s) error {
	return s.repository.Create(%s)
}

func (s *%sServiceImpl) Update(id uint, %s *model.%s) error {
	return s.repository.Update(id, %s)
}

func (s *%sServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *%sServiceImpl) FindAll() ([]*model.%s, error) {
	return s.repository.FindAll()
}

func (s *%sServiceImpl) FindById(id uint) (*model.%s, error) {
	return s.repository.FindById(id)
}
`,
		modelName,
		currentFolderName,
		modelName,
		currentFolderName,
		modelName,
		structName,
		structName,
		structName,
		modelName,
		structName,
		structName,
		modelName,
		structName,
		structName,
		structName,
		structName,
		modelName,
		structName,
		modelName,
		structName,
		modelName,
		structName,
		modelName,
		structName,
		structName,
		structName,
		structName,
		structName)

	writer.WriteString(content)
	writer.Flush()
	return nil
}

// addLineToNewServicesBlock adds a line inside the `services := &Services{}` block in the `NewServices` function.
func addLineToNewServicesBlock(lines []string, newLine string) []string {
	inBlock := false
	var result []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect the start of the `services := &Services{` block
		if strings.Contains(trimmedLine, "services := &Services{") {
			inBlock = true
		}

		// Add the new line inside the block before the closing brace
		if inBlock && trimmedLine == "}" {
			result = append(result, newLine) // Add the new line before the closing brace
			result = append(result, line)    // Add the original closing brace
			inBlock = false                  // Exit the block
			continue
		}

		result = append(result, line)
	}

	return result
}

// editServicesFile updates services.go to include the new service initialization inside `NewServices`.
func editServicesFile(servicesFile, modelName, structName, currentFolderName string) error {
	// Check if services.go exists
	if _, err := os.Stat(servicesFile); os.IsNotExist(err) {
		return fmt.Errorf("services.go not found at %s", servicesFile)
	}

	// Read the content of services.go
	content, err := os.ReadFile(servicesFile)
	if err != nil {
		return fmt.Errorf("error reading services.go: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// Add imports if not already present
	importRepo := fmt.Sprintf("\t%sRepository \"%s/internal/app/domain/repository/%s\"", modelName, currentFolderName, modelName)
	importService := fmt.Sprintf("\t%sService \"%s/internal/app/domain/service/%s\"", modelName, currentFolderName, modelName)
	if !strings.Contains(string(content), importRepo) {
		lines = insertLineAfter(lines, "import (", importRepo)
	}
	if !strings.Contains(string(content), importService) {
		lines = insertLineAfter(lines, "import (", importService)
	}

	// Add the service field in the Services struct if not present
	serviceField := fmt.Sprintf("\t%sService *%sService.%sServiceImpl", modelName, modelName, structName)
	if !strings.Contains(string(content), serviceField) {
		lines = insertLineAfter(lines, "type Services struct {", serviceField)
	}

	// Add initialization in the `services := &Services{}` block
	initLine := fmt.Sprintf("\t\t%sService: %sService.New%sService(%sRepository.New%sRepository(dbs)),",
		modelName, modelName, structName, modelName, structName)
	lines = addLineToNewServicesBlock(lines, initLine)

	// Write the updated content back to services.go
	return os.WriteFile(servicesFile, []byte(strings.Join(lines, "\n")), 0644)
}

// insertLineAfter inserts a line after the first occurrence of a target in lines.
func insertLineAfter(lines []string, target, lineToInsert string) []string {
	for i, line := range lines {
		if strings.Contains(line, target) {
			return append(lines[:i+1], append([]string{lineToInsert}, lines[i+1:]...)...)
		}
	}
	return lines
}
