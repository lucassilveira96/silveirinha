package adapters

import (
	"silveirinha/internal/app/adapter/handler"
	"teste/internal/app/domain"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	testepppppHandler *handler.TestepppppHandler
	testeHandler *handler.TesteHandler
}

func NewHandlers(services *domain.Services) *Handlers {
	return &Handlers{}
		testeHandler: handler.NewTesteHandler(services),
		testeHandler: handler.NewTesteHandler(services),
		testepppppHandler: handler.NewTestepppppHandler(services),
}

func (h *Handlers) Configure(server *fiber.App) {
	h.testepppppHandler.Configure(server)
	h.testeHandler.Configure(server)

}
