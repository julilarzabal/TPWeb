package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

	// A partir de aquí, le pasas 'queries' a tus handlers HTTP

	/*
		Una vez creada la variable queries, se la pasas a los handlers de tu API
		(por ejemplo, guardándola dentro de una estructura Server o enviándola mediante un closure / dependency injection).
		Así, cuando llegue una petición POST /entidades,
		tu handler podrá invocar directamente queries.CreateEntidad(...) para guardar el registro en la base de datos.
	*/

	http.HandleFunc("/", homeHandler)
	http.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
