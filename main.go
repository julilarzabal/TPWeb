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
	"TPEspecial.com/ProgramacionWeb/TPEspecial/views"
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

// Handler principal para GET /
func handleIndex(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Obtener todos los registros de la base de datos usando sqlc
		// IndexPage necesita tanto los eventos como las categorías, así que hacemos dos consultas.
		events, err := queries.ListEvents(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener datos", http.StatusInternalServerError)
			return
		}

		categories, err := queries.ListCategories(r.Context())
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

	// Maneja la ruta principal con el handler que obtiene datos de la base de datos
	http.HandleFunc("/", handleIndex(queries))

	// Servir los archivos estáticos (CSS, imágenes)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func requiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}
