package commands

import (
	"fmt"
	"github.com/lucassilveira96/silveirinha/utils"
	"os"
	"path/filepath"
	"strings"
)

// GenerateHandler generates a handler file for a given model in Go.
func GenerateHandler(modelName, structName string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	// Extract current folder name
	currentFolderName := filepath.Base(currentDir)

	// Define the handler directory and file path
	handlerDir := filepath.Join("internal", "app", "adapter", "handler")
	handlerFilePath := filepath.Join(handlerDir, fmt.Sprintf("%sHandler.go", modelName))

	// Ensure the handler directory exists
	if err := os.MkdirAll(handlerDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating handler directory: %v", err)
	}

	// Create and write to the handler file
	if err := writeHandlerFile(handlerFilePath, currentFolderName, modelName, structName); err != nil {
		return fmt.Errorf("error writing handler file: %v", err)
	}

	// Update the `handlers.go` file
	handlersFilePath := filepath.Join("internal", "app", "adapter", "handlers.go")
	if err := updateHandlersFile(handlersFilePath, modelName, structName, currentFolderName); err != nil {
		return fmt.Errorf("error updating handlers.go: %v", err)
	}

	fmt.Printf("Handler successfully created at: %s\n", handlerFilePath)
	return nil
}

// writeHandlerFile generates the content of the handler file, including Swagger documentation.
func writeHandlerFile(filePath, currentFolderName, modelName, structName string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating handler file: %v", err)
	}
	defer file.Close()

	content := fmt.Sprintf(`package handler

import (
	"%s/internal/app/domain"
	"%s/internal/app/domain/model"
	"%s/internal/app/transport/presenter"
	"%s/internal/infra/variables"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type %sHandler struct {
	services *domain.Services
}

func New%sHandler(services *domain.Services) *%sHandler {
	return &%sHandler{
		services: services,
	}
}

func (h *%sHandler) Configure(server *fiber.App) {
	route := variables.PrefixRoute()
	server.Get(route+"/swagger/*", swagger.HandlerDefault)

	// %s Routes
	serviceRoute := route + "/%s"
	server.Get(serviceRoute, h.getAll%ss)
	server.Get(serviceRoute+"/:id", h.get%sById)
	server.Post(serviceRoute, h.create%s)
	server.Put(serviceRoute+"/:id", h.update%s)
	server.Delete(serviceRoute+"/:id", h.delete%s)
}

// @Summary Get all %ss
// @Description Get all %ss from the system
// @Tags %ss
// @Accept json
// @Produce json
// @Success 200 {array} model.%s "Success"
// @Router /api/v1/%s [get]
func (h *%sHandler) getAll%ss(c *fiber.Ctx) error { 
	%ss, err := h.services.%sService.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Data retrieved successfully", %ss))
}

// @Summary Get %s by ID
// @Description Get a %s by ID from the system
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path int true "%s ID"
// @Success 200 {object} model.%s "Success"
// @Router /api/v1/%s/{id} [get]
func (h *%sHandler) get%sById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	%s, err := h.services.%sService.FindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "%s not found"})
	}
	return c.JSON(%s)
}

// @Summary Create a new %s
// @Description Create a new %s in the system
// @Tags %ss
// @Accept json
// @Produce json
// @Param %s body model.%s true "%s Data"
// @Success 201 {object} model.%s "Created"
// @Router /api/v1/%s [post]
func (h *%sHandler) create%s(c *fiber.Ctx) error {
	%s := new(model.%s) 
	if err := c.BodyParser(%s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.services.%sService.Create(%s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(presenter.Success("Success", %s))
}

// @Summary Update an existing %s
// @Description Update a %s by ID in the system
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path int true "%s ID"
// @Param %s body model.%s true "%s Data"
// @Success 200 {object} model.%s "Updated"
// @Router /api/v1/%s/{id} [put]
func (h *%sHandler) update%s(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	%s := new(model.%s)
	if err := c.BodyParser(%s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	%s.ID = uint(id)
	if err := h.services.%sService.Update(%s.ID, %s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Updated successfully", %s))
}

// @Summary Delete a %s
// @Description Delete a %s by ID in the system
// @Tags %ss
// @Accept json
// @Produce json
// @Param id path int true "%s ID"
// @Success 204 "Deleted successfully"
// @Router /api/v1/%s/{id} [delete]
func (h *%sHandler) delete%s(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.services.%sService.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Deleted successfully", nil))
}
`,
		currentFolderName,
		currentFolderName,
		currentFolderName,
		currentFolderName,
		structName,
		structName,
		structName,
		structName,
		structName,
		modelName,
		modelName,
		structName,
		structName,
		structName,
		structName,
		structName,
		structName, //findall
		structName,
		structName,
		structName,
		utils.ToUrlCase(modelName),
		structName,
		structName,
		modelName,
		structName,
		modelName,  //findall
		structName, //findby
		structName,
		structName,
		structName,
		structName,
		utils.ToUrlCase(modelName),
		structName,
		structName,
		modelName,
		structName,
		structName,
		modelName,  //findby
		structName, //create
		structName,
		structName,
		structName,
		structName,
		structName,
		structName,
		utils.ToUrlCase(modelName),
		structName,
		structName,
		modelName,
		structName,
		modelName,
		structName,
		modelName,
		modelName,  //create
		structName, //initial update
		structName,
		structName,
		structName,
		structName,
		structName,
		structName,
		structName,
		utils.ToUrlCase(modelName),
		structName,
		structName,
		modelName,
		structName,
		modelName,
		modelName,
		structName,
		modelName,
		modelName,
		modelName,  //finish update
		structName, //delete
		structName,
		structName,
		structName,
		utils.ToUrlCase(modelName),
		structName,
		structName,
		structName)

	_, err = file.WriteString(content)
	return err
}

// addLineToNewHandlersBlock adds a line inside the `return &Handlers{}` block in the `NewHandlers` function.
// Handles both empty and non-empty blocks.
func addLineToNewHandlersBlock(lines []string, newLine string) []string {
	inBlock := false
	var result []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect the start of the `return &Handlers{` block
		if strings.Contains(trimmedLine, "return &Handlers{") {
			inBlock = true
			result = append(result, line) // Add the line containing "return &Handlers{"

			// Check if the block is immediately closed (e.g., `&Handlers{}`)
			if strings.HasSuffix(trimmedLine, "}") {
				// Split the line and insert the new line inside
				result[len(result)-1] = strings.Replace(line, "}", fmt.Sprintf("\n%s\n}", newLine), 1)
				inBlock = false
			}
			continue
		}

		// Add the new line before the closing brace if in the block
		if inBlock && trimmedLine == "}" {
			result = append(result, newLine) // Add the new line inside the block
			result = append(result, line)    // Add the closing brace
			inBlock = false                  // Exit the block
			continue
		}

		// Ensure lines outside the block are added normally
		result = append(result, line)
	}

	return result
}

// updateHandlersFile updates the handlers.go file to include the new handler.
func updateHandlersFile(handlersFilePath, modelName, structName, currentFolderName string) error {
	// Read the content of handlers.go
	content, err := os.ReadFile(handlersFilePath)
	if err != nil {
		return fmt.Errorf("error reading handlers.go: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// Add import statement for the handler
	importLine := fmt.Sprintf("\t\"%s/internal/app/adapter/handler\"", currentFolderName)
	if !strings.Contains(string(content), importLine) {
		lines = insertLineHandlerAfter(lines, "import (", importLine)
	}

	// Add the handler field in the Handlers struct
	handlerField := fmt.Sprintf("\t%sHandler *handler.%sHandler", modelName, structName)
	if !strings.Contains(string(content), handlerField) {
		lines = insertLineHandlerAfter(lines, "type Handlers struct {", handlerField)
	}

	// Add initialization of the handler in NewHandlers
	handlerInit := fmt.Sprintf("\t\t%sHandler: handler.New%sHandler(services),", modelName, structName)
	lines = addLineToNewHandlersBlock(lines, handlerInit)

	// Add Configure call in Handlers.Configure
	configureCall := fmt.Sprintf("\th.%sHandler.Configure(server)", modelName)
	lines = insertLineHandlerAfter(lines, "func (h *Handlers) Configure(server *fiber.App) {", configureCall)

	// Write the updated content back to handlers.go
	return os.WriteFile(handlersFilePath, []byte(strings.Join(lines, "\n")), 0644)
}

// insertLineHandlerAfter inserts a line after the first occurrence of a target in lines.
func insertLineHandlerAfter(lines []string, target, lineToInsert string) []string {
	for i, line := range lines {
		if strings.Contains(line, target) {
			return append(lines[:i+1], append([]string{lineToInsert}, lines[i+1:]...)...)
		}
	}
	return lines
}
