package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	connURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(requiredEnv("DB_USER"), requiredEnv("DB_PASSWORD")),
		Host:   net.JoinHostPort(requiredEnv("DB_HOST"), requiredEnv("DB_PORT")),
		Path:   requiredEnv("DB_NAME"),
	}
	query := connURL.Query()
	query.Set("sslmode", requiredEnv("DB_SSLMODE"))
	connURL.RawQuery = query.Encode()

	db, err := sql.Open("pgx", connURL.String())
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
	mux.HandleFunc("POST /categories", server.CreateCategoryHandler)
	mux.HandleFunc("GET /categories", server.ListCategoriesHandler)
	mux.HandleFunc("GET /categories/{id}", server.GetCategoryHandler)
	mux.HandleFunc("PUT /categories/{id}", server.UpdateCategoryHandler)
	mux.HandleFunc("DELETE /categories/{id}", server.DeleteCategoryHandler)
	//Rutas para user_preferences
	mux.HandleFunc("POST /user-preferences", server.AddUserPreferenceHandler)
	mux.HandleFunc("GET /user-preferences/{user_id}", server.ListUserPreferencesHandler)
	mux.HandleFunc("PUT /user-preferences/{user_id}", server.UpdateUserPreferencesHandler)
	mux.HandleFunc("DELETE /user-preferences/{user_id}/{category_id}", server.RemoveUserPreferenceHandler)
	// Rutas para event
	mux.HandleFunc("POST /events", server.CreateEventHandler)
	mux.HandleFunc("GET /events", server.ListEventsHandler)
	mux.HandleFunc("GET /events/{id}", server.GetEventHandler)
	mux.HandleFunc("PUT /events/{id}", server.UpdateEventHandler)
	mux.HandleFunc("DELETE /events/{id}", server.DeleteEventHandler)
	//Rutas para user
	mux.HandleFunc("POST /users", server.CreateUserHandler)
	mux.HandleFunc("GET /users", server.ListUsersHandler)
	mux.HandleFunc("GET /users/{id}", server.GetUserHandler)
	mux.HandleFunc("PUT /users/{id}", server.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/{id}", server.DeleteUserHandler)
	// Rutas para event_category
	mux.HandleFunc("POST /events/{id}/categories", server.AddEventCategoryHandler)
	mux.HandleFunc("GET /events/{id}/categories", server.ListEventCategoriesHandler)
	mux.HandleFunc("DELETE /events/{id}/categories/{category_id}", server.RemoveEventCategoryHandler)
	mux.HandleFunc("DELETE /events/{id}/categories", server.DeleteEventCategoriesHandler)
	mux.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}

func requiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}
