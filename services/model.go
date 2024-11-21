package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"https://github.com/lucassilveira96/silveirinha/utils"
)

// GenerateModel generates Go model files for a given model name.
// It creates two files: one in the domain layer and another in the inbound layer.
func GenerateModel(modelName string) error {
	// Convert the name to camelCase for the file and struct
	fileName := utils.ToCamelCase(modelName) // Converts the name to camelCase, e.g., "testeLu"
	structName := strings.Title(fileName)    // Title case for struct (e.g., "TesteLu")

	// Define directories for domain and inbound layers
	domainDir := "internal/app/domain/model"
	inboundDir := "internal/app/transport/inbound"
	mapperDir := "internal/app/transport/mapper"

	// Ensure directories exist
	if err := os.MkdirAll(domainDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating domain directory: %v", err)
	}
	if err := os.MkdirAll(inboundDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating inbound directory: %v", err)
	}
	if err := os.MkdirAll(mapperDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating mapper directory: %v", err)
	}

	// Create the main files in respective directories
	domainFilePath := fmt.Sprintf("%s/%s.go", domainDir, fileName)
	inboundFilePath := fmt.Sprintf("%s/%s.go", inboundDir, fileName)
	mapperFilePath := fmt.Sprintf("%s/%sMapToModel.go", mapperDir, fileName)

	// Write domain and inbound model files
	if err := writeModelFile(domainFilePath, structName, inboundFilePath); err != nil {
		return fmt.Errorf("error writing domain file: %v", err)
	}

	fmt.Printf("Model files generated:\n- %s\n- %s\n", domainFilePath, inboundFilePath)

	// Write the mapper file to map inbound to domain
	if err := writeMapperFile(mapperFilePath, fileName, structName); err != nil {
		return fmt.Errorf("error writing mapper file: %v", err)
	}

	fmt.Printf("Model and Mapper files generated:\n- %s\n- %s\n- %s\n", domainFilePath, inboundFilePath, mapperFilePath)

	err := GenerateRepository(modelName)
	if err != nil {
		return fmt.Errorf("error generating repository: %v", err)
	}

	return nil
}

// writeModelFile creates a model file and writes its struct definition interactively.
// Users specify attributes and their properties.
func writeModelFile(filePath, structName, inboundFilePath string) error {
	// Open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write package declaration and imports
	writer.WriteString("package model\n\n")
	writer.WriteString(`import "time"` + "\n\n")

	// Start defining the struct
	writer.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	writer.WriteString("\tID uint `gorm:\"primaryKey;autoIncrement\" json:\"id\"`\n")

	// Prompt user to add attributes
	for {
		fmt.Print("Add an attribute to the struct? (y/n): ")
		var choice string
		fmt.Scanln(&choice)

		if strings.ToLower(choice) != "y" {
			break
		}

		var attrName, attrType, defaultValue string
		var nullable, hasDefault bool

		// Collect attribute name
		fmt.Print("Attribute name: ")
		fmt.Scanln(&attrName)

		// Collect attribute type
		attrType = selectType()

		// Determine if the attribute is nullable
		fmt.Print("Is it nullable? (y/n): ")
		fmt.Scanln(&choice)
		nullable = strings.ToLower(choice) == "y"
		if nullable {
			attrType = "*" + attrType
		}

		// Determine if the attribute has a default value
		if !nullable {
			fmt.Print("Has a default value? (y/n): ")
			fmt.Scanln(&choice)
			hasDefault = strings.ToLower(choice) == "y"
			if hasDefault {
				fmt.Print("Default value: ")
				fmt.Scanln(&defaultValue)
			}
		}

		// Add attribute definition
		jsonName := utils.ToSnakeCase(attrName)
		gormTag := ""
		if hasDefault {
			gormTag = fmt.Sprintf(`gorm:"default:%s"`, defaultValue)
		} else if !nullable {
			gormTag = `gorm:"not null"`
		}
		writer.WriteString(fmt.Sprintf("\t%s %s `%s json:\"%s\"`\n",
			utils.ToPascalCase(attrName), attrType, gormTag, jsonName))

		// Also add the attribute to the inbound model file
		writeInboundModelFile(inboundFilePath, structName, attrName, attrType, jsonName)
	}

	// Prompt user to add relationships
	for {
		fmt.Print("Would you like to add a relationship? (y/n): ")
		var choice string
		fmt.Scanln(&choice)

		if strings.ToLower(choice) != "y" {
			break
		}

		var relatedModel string
		fmt.Print("Enter the name of the related model: ")
		fmt.Scanln(&relatedModel)

		relationshipName := utils.ToPascalCase(relatedModel)
		relationshipNameSnake := utils.ToSnakeCase(relatedModel)
		relatedField := fmt.Sprintf("%sId", relationshipName)
		relatedFieldSnake := utils.ToSnakeCase(relatedField)

		// Add the foreign key field
		writer.WriteString(fmt.Sprintf("\t%s uint `json:\"%s\"`\n", relatedField, relatedFieldSnake))

		// Add the relationship
		writer.WriteString(fmt.Sprintf("\t%s %s `gorm:\"foreignKey:%s\" json:\"%s\"`\n",
			relationshipName, relationshipName, relatedField, relationshipNameSnake))

		// Also add the relationship to the inbound model file
		writeInboundModelFile(inboundFilePath, structName, relatedField, "uint", relatedFieldSnake)
	}

	// Ask user if they want to include standard date fields
	fmt.Print("Include standard date fields (CreatedAt, UpdatedAt, DeletedAt)? (y/n): ")
	var includeDatesInput string
	fmt.Scanln(&includeDatesInput)
	includeDates := strings.ToLower(includeDatesInput) == "y"

	// Include timestamps if required
	if includeDates {
		writer.WriteString("\tCreatedAt time.Time `gorm:\"autoCreateTime;not null\" json:\"created_at\"`\n")
		writer.WriteString("\tUpdatedAt time.Time `gorm:\"autoUpdateTime;not null\" json:\"updated_at\"`\n")
		writer.WriteString("\tDeletedAt *time.Time `gorm:\"index\" json:\"deleted_at\"`\n")
	}

	// Close the struct definition
	writer.WriteString("}\n\n")

	// Add TableName method for GORM
	writer.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", structName))
	writer.WriteString(fmt.Sprintf("\treturn \"%s\"\n", utils.ToSnakeCase(structName)))
	writer.WriteString("}\n")

	// Flush writer buffer
	writer.Flush()
	return nil
}

// writeInboundModelFile creates the inbound model file (without ID, date fields, GORM tags, and TableName function)
func writeInboundModelFile(filePath, structName, attrName, attrType, jsonName string) error {
	// Create the file in the inbound directory
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the package declaration for inbound model
	writer.WriteString("package inbound\n\n")

	// Start defining the struct (same as the model, without ID and date fields)
	writer.WriteString(fmt.Sprintf("type %s struct {\n", utils.ToPascalCase(structName)))

	// Add the field to the inbound struct with JSON tags
	writer.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", utils.ToPascalCase(attrName), attrType, jsonName))

	// Close the struct definition
	writer.WriteString("}\n")

	// Flush writer buffer
	writer.Flush()
	return nil
}

// writeMapperFile generates the mapper file to map from inbound model to domain model
func writeMapperFile(filePath, fileName, structName string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	// Get the name of the current directory (the project folder)
	projectFolder := filepath.Base(currentDir)
	fmt.Println("debugger")
	fmt.Println(projectFolder)

	// Create the mapper file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating mapper file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write package declaration and imports
	writer.WriteString("package mapper\n\n")
	writer.WriteString("import (\n")
	writer.WriteString(fmt.Sprintf("\t\"reflect\"\n"))
	writer.WriteString(fmt.Sprintf("\t\"%s/internal/app/domain/model\"\n", projectFolder))
	writer.WriteString(fmt.Sprintf("\t\"%s/internal/app/transport/inbound\"\n", projectFolder))
	writer.WriteString(")\n\n")

	// Define the mapping function
	writer.WriteString(fmt.Sprintf("func %sMapToModel(inbound inbound.%s) model.%s {\n", structName, structName, structName))
	writer.WriteString(fmt.Sprintf("\tvar modelObj model.%s\n", structName))
	writer.WriteString("\tinboundValue := reflect.ValueOf(inbound)\n")
	writer.WriteString("\tmodelValue := reflect.ValueOf(&modelObj).Elem()\n\n")

	// Loop through each field in the inbound struct and map it to the model
	writer.WriteString("\t// Loop through each field in the inbound struct\n")
	writer.WriteString("\tfor i := 0; i < inboundValue.NumField(); i++ {\n")
	writer.WriteString("\t\tfieldName := inboundValue.Type().Field(i).Name\n")
	writer.WriteString("\t\tmodelField := modelValue.FieldByName(fieldName)\n\n")

	writer.WriteString("\t\t// If the field exists in the model, copy the value\n")
	writer.WriteString("\t\tif modelField.IsValid() && modelField.CanSet() {\n")
	writer.WriteString("\t\t\tmodelField.Set(inboundValue.Field(i))\n")
	writer.WriteString("\t\t}\n")
	writer.WriteString("\t}\n\n")

	// Return the mapped model
	writer.WriteString("\treturn modelObj\n")
	writer.WriteString("}\n")

	// Flush the writer buffer
	writer.Flush()

	return nil
}

// ShowGoTypes lists the supported Go types for attributes.
// It displays a menu for user selection during attribute definition.
func ShowGoTypes() {
	fmt.Println("Choose a type for the attribute:")
	fmt.Println("1) int")
	fmt.Println("2) uint")
	fmt.Println("3) int8")
	fmt.Println("4) uint8")
	fmt.Println("5) int16")
	fmt.Println("6) uint16")
	fmt.Println("7) int32")
	fmt.Println("8) uint32")
	fmt.Println("9) int64")
	fmt.Println("10) uint64")
	fmt.Println("11) string")
	fmt.Println("12) float32")
	fmt.Println("13) float64")
	fmt.Println("14) complex64")
	fmt.Println("15) complex128")
	fmt.Println("16) bool")
	fmt.Println("17) byte")
	fmt.Println("18) rune")
	fmt.Println("19) time.Time")
	fmt.Println("20) []byte")
}

// selectType allows users to select a type from a predefined list.
func selectType() string {
	ShowGoTypes()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter the number corresponding to the type: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > 20 {
			fmt.Println("Invalid choice. Please select a valid number.")
			continue
		}

		switch choice {
		case 1:
			return "int"
		case 2:
			return "uint"
		case 3:
			return "int8"
		case 4:
			return "uint8"
		case 5:
			return "int16"
		case 6:
			return "uint16"
		case 7:
			return "int32"
		case 8:
			return "uint32"
		case 9:
			return "int64"
		case 10:
			return "uint64"
		case 11:
			return "string"
		case 12:
			return "float32"
		case 13:
			return "float64"
		case 14:
			return "complex64"
		case 15:
			return "complex128"
		case 16:
			return "bool"
		case 17:
			return "byte"
		case 18:
			return "rune"
		case 19:
			return "time.Time"
		case 20:
			return "[]byte"
		}
	}
}
