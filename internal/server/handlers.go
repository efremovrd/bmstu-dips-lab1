package server

import (
	h "bmstu-dips-lab1/internal/person/delivery/http"
	"bmstu-dips-lab1/internal/person/repo"
	"bmstu-dips-lab1/internal/person/usecase"
)

func (s *Server) MapHandlers() error {
	pRepo := repo.NewPersonRepo(s.db)
	pUC := usecase.NewPersonUseCase(pRepo)
	pH := h.NewPersonHandlers(pUC)

	api := s.router.Group("/api")

	v1 := api.Group("/v1")

	persons := v1.Group("/persons")
	h.MapPersonRoutes(persons, pH)

	return nil
}
