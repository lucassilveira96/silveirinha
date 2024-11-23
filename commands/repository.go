package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
