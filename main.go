package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

// Server actúa como contenedor de las dependencias necesarias para los handlers
// Decidimos utilizarlo para permitir que luego sea escalable y podamos agregar más dependencias si es necesario.
type Server struct {
	queries *sqlc.Queries
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que sea la ruta exacta "/" (para el 404)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Mostrar el formulario HTML
	http.ServeFile(w, r, "index.html")
}

func main() {
	// Establecer la conexión a la base de datos
	connStr := "user=u password=p dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Verificar que la base de datos responda
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	// Crear la instancia del repositorio generado por sqlc
	queries := sqlc.New(db)

	server := &Server{
		queries: queries,
	}

	// Configurar rutas utilizando las funciones del Server
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.homeHandler)

	// Rutas para category
	mux.HandleFunc("POST /categories", server.CreateCategory)
	mux.HandleFunc("GET /categories", server.ListCategories)
	mux.HandleFunc("GET /categories/{id}", server.GetCategory)
	mux.HandleFunc("PUT /categories/{id}", server.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", server.DeleteCategory)
	//Rutas para user_preferences
	mux.HandleFunc("POST /user-preferences", server.AddUserPreference)
	mux.HandleFunc("GET /user-preferences/{user_id}", server.ListUserPreferences)
	mux.HandleFunc("PUT /user-preferences/{user_id}", server.UpdateUserPreferences)
	mux.HandleFunc("DELETE /user-preferences/{user_id}/{category_id}", server.RemoveUserPreference)
	// Rutas para event
	mux.HandleFunc("POST /events", server.CreateEventHandler)
	mux.HandleFunc("GET /events", server.ListEventsHandler)
	mux.HandleFunc("GET /events/{id}", server.GetEventHandler)
	mux.HandleFunc("PUT /events/{id}", server.UpdateEventHandler)
	mux.HandleFunc("DELETE /events/{id}", server.DeleteEventHandler)
	//Rutas para user
	mux.HandleFunc("POST /users", server.createUser)
	mux.HandleFunc("GET /users", server.listUsers)
	mux.HandleFunc("GET /users/{id}", server.getUser)
	mux.HandleFunc("PUT /users/{id}", server.updateUser)
	mux.HandleFunc("DELETE /users/{id}", server.deleteUser)
	// Rutas para event_category
	mux.HandleFunc("POST /events/{id}/categories", server.AddEventCategoryHandler)
	mux.HandleFunc("GET /events/{id}/categories", server.ListEventCategoriesHandler)
	mux.HandleFunc("DELETE /events/{id}/categories/{category_id}", server.RemoveEventCategoryHandler)
	mux.HandleFunc("DELETE /events/{id}/categories", server.DeleteEventCategoriesHandler)
	mux.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
