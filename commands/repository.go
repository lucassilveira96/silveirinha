package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateRepository generates Go repository files for a given model name.
func GenerateRepository(modelName, structName string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	// Get the current folder name
	currentFolderName := filepath.Base(currentDir)

	// Construct the repository directory path
	repositoryDir := filepath.Join("internal", "app", "domain", "repository", modelName)

	// Ensure the repository directory exists
	if err := os.MkdirAll(repositoryDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating repository directory: %v", err)
	}

	// Generate the Repository interface file
	repositoryFilePath := filepath.Join(repositoryDir, fmt.Sprintf("%sRepository.go", modelName))
	if err := writeRepositoryInterfaceFile(repositoryFilePath, currentFolderName, modelName, structName); err != nil {
		return fmt.Errorf("error writing repository interface file: %v", err)
	}

	// Generate the Repository implementation file
	repositoryImplFilePath := filepath.Join(repositoryDir, fmt.Sprintf("%sRepositoryImpl.go", modelName))
	if err := writeRepositoryImplFile(repositoryImplFilePath, currentFolderName, modelName, structName); err != nil {
		return fmt.Errorf("error writing repository implementation file: %v", err)
	}

	fmt.Printf("Repository files generated in: %s\n", repositoryDir)

	if err := addModelToMigrations(structName, currentFolderName); err != nil {
		return fmt.Errorf("error writing migrations file: %v", err)
	}
	return nil
}

// writeRepositoryInterfaceFile creates the repository interface file.
func writeRepositoryInterfaceFile(filePath, currentFolderName, modelName, structName string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the content
	content := fmt.Sprintf(`package %sRepository

import "%s/internal/app/domain/model"

type %sRepository interface {
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

// writeRepositoryImplFile creates the repository implementation file.
func writeRepositoryImplFile(filePath, currentFolderName, modelName, structName string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the content
	content := fmt.Sprintf(`package %sRepository

import (
	"%s/internal/app/domain/model"
	"%s/internal/infra/database"
	"time"
)

var _ %sRepository = (*%sRepositoryImpl)(nil)

type %sRepositoryImpl struct {
	db *database.Databases
}

func New%sRepository(db *database.Databases) *%sRepositoryImpl {
	return &%sRepositoryImpl{db: db}
}

func (r *%sRepositoryImpl) Create(%s *model.%s) error {
	return r.db.Write.Create(%s).Error
}

func (r *%sRepositoryImpl) Update(id uint, %s *model.%s) error {
	existing := &model.%s{}
	if err := r.db.Write.First(existing, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(existing).Updates(%s).Error
}

func (r *%sRepositoryImpl) Delete(id uint) error {
	%s := &model.%s{}
	if err := r.db.Write.First(%s, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(%s).Update("deleted_at", time.Now()).Error
}

func (r *%sRepositoryImpl) FindAll() ([]*model.%s, error) {
	var %ss []*model.%s
	err := r.db.Read.
		Where("deleted_at IS NULL").
		Find(&%ss).Error
	return %ss, err
}

func (r *%sRepositoryImpl) FindById(id uint) (*model.%s, error) {
	var %s model.%s
	err := r.db.Read.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&%s).Error
	return &%s, err
}
`,
		modelName,
		currentFolderName,
		currentFolderName,
		structName,
		structName,
		structName,
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
		structName,
		modelName,
		structName,
		modelName,
		structName,
		modelName,
		modelName,
		structName,
		structName,
		modelName,
		structName,
		modelName,
		modelName,
		structName,
		structName,
		modelName,
		structName,
		modelName,
		modelName)

	writer.WriteString(content)
	writer.Flush()
	return nil
}

func addModelToMigrations(structName, currentFolderName string) error {
	// Path to the databases.go file
	databasesFilePath := filepath.Join("internal", "infra", "database", "databases.go")

	// Read the existing content of the file
	content, err := os.ReadFile(databasesFilePath)
	if err != nil {
		return fmt.Errorf("error reading databases.go: %v", err)
	}

	// Convert to string for modification
	contentStr := string(content)

	// Check if the model is already added to migrations
	if containsModel(contentStr, structName) {
		fmt.Printf("Model %s already exists in migrations.\n", structName)
		return nil
	}

	// Construct the model import path dynamically
	modelImport := fmt.Sprintf(`"%s/internal/app/domain/model"`, currentFolderName)
	if !containsImport(contentStr, modelImport) {
		contentStr = addImport(contentStr, modelImport)
	}

	// Add the model to the runMigrations function
	migrationEntry := fmt.Sprintf("&model.%s{},", structName)
	contentStr = addMigrationEntry(contentStr, migrationEntry)

	// Write the updated content back to the file
	err = os.WriteFile(databasesFilePath, []byte(contentStr), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing to databases.go: %v", err)
	}

	fmt.Printf("Model %s added to migrations.\n", structName)
	return nil
}

func containsModel(content, structName string) bool {
	return filepath.Base(content) == structName
}

func containsImport(content, importStr string) bool {
	return filepath.Base(content) == importStr
}

func addImport(content, importStr string) string {
	importSection := "import ("
	index := strings.Index(content, importSection) // Use strings.Index to find the position of "import ("
	if index < 0 {
		return content // Return the content unchanged if "import (" is not found
	}
	return content[:index+len(importSection)] + "\n\t" + importStr + content[index+len(importSection):]
}

func addMigrationEntry(content, entry string) string {
	migrationFunc := "func (d *Databases) runMigrations(db *gorm.DB) {"
	index := strings.Index(content, migrationFunc) // Locate the runMigrations function
	if index < 0 {
		return content // If runMigrations function not found, return content as is
	}

	// Find the AutoMigrate call within the runMigrations function
	autoMigrateIndex := strings.Index(content[index:], "db.AutoMigrate(")
	if autoMigrateIndex < 0 {
		return content // If AutoMigrate is not found, return content as is
	}

	// Calculate the start of the db.AutoMigrate call
	start := index + autoMigrateIndex + len("db.AutoMigrate(")

	// Find the closing parenthesis for the AutoMigrate call
	end := strings.Index(content[start:], ")")
	if end < 0 {
		return content // If the closing parenthesis is not found, return content as is
	}

	// Extract the current content inside db.AutoMigrate
	migrateContent := strings.TrimSpace(content[start : start+end])

	// Clean up trailing commas in the existing AutoMigrate entries
	if strings.HasSuffix(migrateContent, ",") {
		migrateContent = strings.TrimSuffix(migrateContent, ",")
	}

	// Add the new entry
	if migrateContent == "" {
		// If AutoMigrate is empty, add the entry
		migrateContent = "\n\t\t" + entry
	} else {
		// If AutoMigrate has existing entries, append the new entry
		migrateContent += ",\n\t\t" + entry
	}

	// Replace the content inside db.AutoMigrate
	return content[:start] + migrateContent + content[start+end:]
}
