package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"https://github.com/lucassilveira96/silveirinha/utils"
)

// GenerateRepository generates Go repository files for a given model name.
func GenerateRepository(modelName string) error {
	// Convert the model name to camelCase for the file and struct
	fileName := utils.ToCamelCase(modelName)
	structName := strings.Title(fileName)

	// Get the dynamic project path for the internal directory
	repositoryDir, err := getRepositoryDirPath(modelName) // Pass the model name to the directory function
	if err != nil {
		return fmt.Errorf("error getting repository directory path: %v", err)
	}

	// Ensure the repository directory exists
	if err := os.MkdirAll(repositoryDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating repository directory: %v", err)
	}

	// Generate the Repository interface file
	repositoryFilePath := fmt.Sprintf("%s/%sRepository.go", repositoryDir, fileName)
	if err := writeRepositoryInterfaceFile(repositoryFilePath, modelName, structName, repositoryDir); err != nil {
		return fmt.Errorf("error writing repository interface file: %v", err)
	}

	// Generate the RepositoryImpl file
	repositoryImplFilePath := fmt.Sprintf("%s/%sRepositoryImpl.go", repositoryDir, fileName)
	if err := writeRepositoryImplFile(repositoryImplFilePath, modelName, structName, repositoryDir); err != nil {
		return fmt.Errorf("error writing repository implementation file: %v", err)
	}

	fmt.Printf("Repository files generated in: %s\n", repositoryDir)

	return nil
}

// getRepositoryDirPath retrieves the dynamic path for the repository directory based on the model name.
func getRepositoryDirPath(modelName string) (string, error) {
	// Get the current working directory (the directory where the program is running)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %v", err)
	}

	// Construct the path to the repository directory using the model name
	repositoryDir := filepath.Join(currentDir, "internal", "app", "domain", "repository", modelName)

	return repositoryDir, nil
}

// writeRepositoryInterfaceFile creates the repository interface file.
func writeRepositoryInterfaceFile(filePath, modelName, structName, repositoryDir string) error {
	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the package declaration and imports with the correct module path
	writer.WriteString(fmt.Sprintf("package %sRepository\n\n", repositoryDir))
	writer.WriteString(fmt.Sprintf(`import (
		"%s/internal/app/domain/model"
	)`+"\n\n", repositoryDir))

	// Write the interface definition
	writer.WriteString(fmt.Sprintf("type %sRepository interface {\n", structName))
	writer.WriteString(fmt.Sprintf("\tCreate(%s *model.%s) error\n", modelName, structName))
	writer.WriteString(fmt.Sprintf("\tUpdate(id uint, %s *model.%s) error\n", modelName, structName))
	writer.WriteString(fmt.Sprintf("\tDelete(id uint) error\n", modelName))
	writer.WriteString(fmt.Sprintf("\tFindAll() ([]*model.%s, error)\n", modelName))
	writer.WriteString(fmt.Sprintf("\tFindById(id uint) (*model.%s, error)\n", modelName))
	writer.WriteString("}\n")

	// Flush the writer buffer
	writer.Flush()
	return nil
}

// writeRepositoryImplFile creates the repository implementation file.
func writeRepositoryImplFile(filePath, modelName, structName, repositoryDir string) error {
	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the package declaration and imports with the correct module path
	writer.WriteString(fmt.Sprintf("package %sRepository\n\n", modelName))
	writer.WriteString(fmt.Sprintf(`import (
		"%s/internal/app/domain/model"
		"%s/internal/infra/database"
		"time"
	)`+"\n\n", repositoryDir, repositoryDir))

	// Write the repository implementation struct and constructor
	repositoryName := strings.Title(modelName)
	writer.WriteString(fmt.Sprintf("var _ %sRepository = (*%sRepositoryImpl)(nil)\n\n", repositoryName, repositoryName))
	writer.WriteString(fmt.Sprintf("type %sRepositoryImpl struct {\n", repositoryName))
	writer.WriteString("\tdb *database.Databases\n")
	writer.WriteString("}\n\n")
	writer.WriteString(fmt.Sprintf("func New%sRepository(db *database.Databases) *%sRepositoryImpl {\n", structName, repositoryName))
	writer.WriteString(fmt.Sprintf("\treturn &%sRepositoryImpl{db: db}\n", repositoryName))
	writer.WriteString("}\n\n")

	// Write the functions for the repository implementation
	writeRepositoryImplFunction(writer, "Create", modelName, structName, "r.db.Write.Create")
	writeRepositoryImplFunction(writer, "Update", modelName, structName, "r.db.Write.Model(existing).Create")
	writeRepositoryImplFunction(writer, "Delete", modelName, structName, "r.db.Write.Model(%s).Update(\"deleted_at\", time.Now())")
	writeRepositoryImplFunction(writer, "FindAll", modelName, structName, "r.db.Read.Where(\"deleted_at IS NULL\").Find(&%s)")
	writeRepositoryImplFunction(writer, "FindById", modelName, structName, "r.db.Read.Where(\"id = ? AND deleted_at IS NULL\", id).First(&%s)")

	// Flush the writer buffer
	writer.Flush()
	return nil
}

// writeRepositoryImplFunction generates a function for the repository implementation.
func writeRepositoryImplFunction(writer *bufio.Writer, functionName, modelName, structName, dbCall string) {
	writer.WriteString(fmt.Sprintf("func (r *%sRepositoryImpl) %s(%s *model.%s) error {\n", functionName, modelName, structName))
	writer.WriteString(fmt.Sprintf("\t%s(%s).Error\n", dbCall, modelName))
	writer.WriteString("\treturn err\n")
	writer.WriteString("}\n\n")
}
