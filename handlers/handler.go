package handlers

import (
	"net/http"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
	"TPEspecial.com/ProgramacionWeb/TPEspecial/views"
)

// Creas el struct que guardará las dependencias
type Server struct {
	queries *sqlc.Queries
}

// Función constructora (recomendada para inicializar el struct)
func NewServer(q *sqlc.Queries) *Server {
	return &Server{queries: q}
}

// Handler principal para GET /
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todos los registros de la base de datos usando sqlc
	// IndexPage necesita tanto los eventos como las categorías, así que hacemos dos consultas.
	events, err := s.queries.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener datos", http.StatusInternalServerError)
		return
	}

	categories, err := s.queries.ListCategories(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener datos", http.StatusInternalServerError)
		return
	}

	// 2. Llamar al componente de página principal pasándole los datos
	component := views.IndexPage(events, categories)

	// 3. Renderizar el componente completo en el http.ResponseWriter
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
	}
}
