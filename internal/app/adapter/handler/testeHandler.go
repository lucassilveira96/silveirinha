package handler

import (
	"silveirinha/internal/app/domain"
	"silveirinha/internal/app/domain/model"
	"silveirinha/internal/app/transport/presenter"
	"silveirinha/internal/infra/variables"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type TesteHandler struct {
	services *domain.Services
}

func NewTesteHandler(services *domain.Services) *TesteHandler {
	return &TesteHandler{
		services: services,
	}
}

func (h *TesteHandler) Configure(server *fiber.App) {
	route := variables.PrefixRoute()
	server.Get(route+"/swagger/*", swagger.HandlerDefault)

	// teste Routes
	serviceRoute := route + "/teste"
	server.Get(serviceRoute, h.getAllTestes)
	server.Get(serviceRoute+"/:id", h.getTesteById)
	server.Post(serviceRoute, h.createTeste)
	server.Put(serviceRoute+"/:id", h.updateTeste)
	server.Delete(serviceRoute+"/:id", h.deleteTeste)
}

var validate = validator.New()

// @Summary Get all Testes
// @Description Get all Testes from the system
// @Tags Testes
// @Accept json
// @Produce json
// @Success 200 {array} model.Teste "Success"
// @Router /api/v1/teste [get]
func (h *TesteHandler) getAllTestes(c *fiber.Ctx) error { 
	testes, err := h.services.TesteService.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Data retrieved successfully", testes))
}

// @Summary Get Teste by ID
// @Description Get a Teste by ID from the system
// @Tags Testes
// @Accept json
// @Produce json
// @Param id path int true "Teste ID"
// @Success 200 {object} model.Teste "Success"
// @Router /api/v1/teste/{id} [get]
func (h *TesteHandler) getTesteById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	teste, err := h.services.TesteService.FindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Teste not found"})
	}
	return c.JSON(teste)
}

// @Summary Create a new Teste
// @Description Create a new Teste in the system
// @Tags Testes
// @Accept json
// @Produce json
// @Param Teste body model.Teste true "Teste Data"
// @Success 201 {object} model.Teste "Created"
// @Router /api/v1/teste [post]
func (h *TesteHandler) createTeste(c *fiber.Ctx) error {
	teste := new(model.Teste) 
	if err := c.BodyParser(teste); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.services.TesteService.Create(teste); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(presenter.Success("Success", teste))
}

// @Summary Update an existing Teste
// @Description Update a Teste by ID in the system
// @Tags Testes
// @Accept json
// @Produce json
// @Param id path int true "Teste ID"
// @Param Teste body model.Teste true "Teste Data"
// @Success 200 {object} model.Teste "Updated"
// @Router /api/v1/teste/{id} [put]
func (h *TesteHandler) updateTeste(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	teste := new(model.Teste)
	if err := c.BodyParser(teste); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	teste.ID = uint(id)
	if err := h.services.TesteService.Update(teste.ID, teste); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Updated successfully", teste))
}

// @Summary Delete a Teste
// @Description Delete a Teste by ID in the system
// @Tags Testes
// @Accept json
// @Produce json
// @Param id path int true "Teste ID"
// @Success 204 "Deleted successfully"
// @Router /api/v1/teste/{id} [delete]
func (h *TesteHandler) deleteTeste(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.services.TesteService.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Deleted successfully", nil))
}
