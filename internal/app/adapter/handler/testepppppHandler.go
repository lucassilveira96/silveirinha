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

type TestepppppHandler struct {
	services *domain.Services
}

func NewTestepppppHandler(services *domain.Services) *TestepppppHandler {
	return &TestepppppHandler{
		services: services,
	}
}

func (h *TestepppppHandler) Configure(server *fiber.App) {
	route := variables.PrefixRoute()
	server.Get(route+"/swagger/*", swagger.HandlerDefault)

	// testeppppp Routes
	serviceRoute := route + "/testeppppp"
	server.Get(serviceRoute, h.getAllTesteppppps)
	server.Get(serviceRoute+"/:id", h.getTestepppppById)
	server.Post(serviceRoute, h.createTesteppppp)
	server.Put(serviceRoute+"/:id", h.updateTesteppppp)
	server.Delete(serviceRoute+"/:id", h.deleteTesteppppp)
}

var validate = validator.New()

// @Summary Get all Testeppppps
// @Description Get all Testeppppps from the system
// @Tags Testeppppps
// @Accept json
// @Produce json
// @Success 200 {array} model.Testeppppp "Success"
// @Router /api/v1/testeppppp [get]
func (h *TestepppppHandler) getAllTesteppppps(c *fiber.Ctx) error { 
	testeppppps, err := h.services.TestepppppService.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Data retrieved successfully", testeppppps))
}

// @Summary Get Testeppppp by ID
// @Description Get a Testeppppp by ID from the system
// @Tags Testeppppps
// @Accept json
// @Produce json
// @Param id path int true "Testeppppp ID"
// @Success 200 {object} model.Testeppppp "Success"
// @Router /api/v1/testeppppp/{id} [get]
func (h *TestepppppHandler) getTestepppppById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	testeppppp, err := h.services.TestepppppService.FindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Testeppppp not found"})
	}
	return c.JSON(testeppppp)
}

// @Summary Create a new Testeppppp
// @Description Create a new Testeppppp in the system
// @Tags Testeppppps
// @Accept json
// @Produce json
// @Param Testeppppp body model.Testeppppp true "Testeppppp Data"
// @Success 201 {object} model.Testeppppp "Created"
// @Router /api/v1/testeppppp [post]
func (h *TestepppppHandler) createTesteppppp(c *fiber.Ctx) error {
	testeppppp := new(model.Testeppppp) 
	if err := c.BodyParser(testeppppp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.services.TestepppppService.Create(testeppppp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(presenter.Success("Success", testeppppp))
}

// @Summary Update an existing Testeppppp
// @Description Update a Testeppppp by ID in the system
// @Tags Testeppppps
// @Accept json
// @Produce json
// @Param id path int true "Testeppppp ID"
// @Param Testeppppp body model.Testeppppp true "Testeppppp Data"
// @Success 200 {object} model.Testeppppp "Updated"
// @Router /api/v1/testeppppp/{id} [put]
func (h *TestepppppHandler) updateTesteppppp(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	testeppppp := new(model.Testeppppp)
	if err := c.BodyParser(testeppppp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	testeppppp.ID = uint(id)
	if err := h.services.TestepppppService.Update(testeppppp.ID, testeppppp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Updated successfully", testeppppp))
}

// @Summary Delete a Testeppppp
// @Description Delete a Testeppppp by ID in the system
// @Tags Testeppppps
// @Accept json
// @Produce json
// @Param id path int true "Testeppppp ID"
// @Success 204 "Deleted successfully"
// @Router /api/v1/testeppppp/{id} [delete]
func (h *TestepppppHandler) deleteTesteppppp(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.services.TestepppppService.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Success("Deleted successfully", nil))
}
