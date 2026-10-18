package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
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

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Estructura para recibir el cuerpo JSON de la petición
	/*
		Es una estructura anónima que se usa como un DTO (Data Transfer Object).
		Su única función es definir qué campos JSON esperas que el cliente
		(por ejemplo, Postman, un frontend o curl) te envíe en el cuerpo (body) de la petición HTTP.
	*/
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		UserID      int32  `json:"user_id"`
	}

	// Carga de datos en req
	/*
		r.Body: Contiene el flujo de bytes (stream) con el JSON que envió el usuario.
		json.NewDecoder(...).Decode(&req): Lee ese JSON, analiza sus claves y rellena los campos de la variable req.
	*/
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Validación (ejemplo: validar que el título no esté vacío)
	if req.Title == "" {
		http.Error(w, "El título es obligatorio", http.StatusBadRequest)
		return
	}

	// 3. Ejecutar la consulta generada por sqlc
	/*
		s.queries.CreateEvent: Llama a la función generada por sqlc para crear un nuevo evento.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.CreateEventParams{...}: Estructura con los parámetros para la consulta.
	*/
	event, err := s.queries.CreateEvent(r.Context(), db.CreateEventParams{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
	})

	/*
		http.Error(): Es una función auxiliar de Go para devolver respuestas de error en texto plano. Recibe 3 parámetros:
		http.Error(w, "Mensaje de error", CÓDIGO_HTTP)
		(Ideal para responder errores)
	*/
	if err != nil {
		http.Error(w, "Error al crear el evento", http.StatusInternalServerError)
		return
	}

	// 4. Devolver la respuesta como JSON con estado 201 Created
	// w.WriteHeader(...): (Ideal para respuestas Exitosas)
	// a. Decirle al navegador/cliente que el contenido devuelto será JSON
	w.Header().Set("Content-Type", "application/json")

	// b. A AQUÍ SE INDICA EL CÓDIGO DE ESTADO (201 Created)
	w.WriteHeader(http.StatusCreated) // O directamente: w.WriteHeader(201)

	// c. Escribir el cuerpo del JSON que devolverás al cliente
	json.NewEncoder(w).Encode(event)
}

func (s *Server) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener los registros mediante sqlc
	/*
		s.queries.ListEvents: Llama a la función generada por sqlc para obtener la lista de eventos.
		r.Context(): Pasa el contexto de la solicitud HTTP.
	*/
	events, err := s.queries.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de eventos", http.StatusInternalServerError)
		return
	}

	// Si no hay datos, devolver un array vacío en lugar de null
	if events == nil {
		events = []db.Event{}
	}

	// 2. Responder en formato JSON
	w.Header().Set("Content-Type", "application/json")
	// No realiza w.WriteHeader(...) porque no debo devolver un código de error, sino un 200 OK por defecto.
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer el ID de la URL (usando net/http de Go)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Buscar en la base de datos con sqlc
	/*
		s.queries.GetEvent: Llama a la función generada por sqlc para obtener un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(id): ID del evento a buscar.
	*/
	event, err := s.queries.GetEvent(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no existe, devolver estado 404
			http.Error(w, "Evento no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	// 3. Devolver el evento hallado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Leer JSON
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		AgeRange    string `json:"age_range"`
		UserID      int32  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 3. Validar datos
	if req.Title == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	// 4. Ejecutar la actualización en sqlc
	/*
		s.queries.UpdateEvent: Llama a la función generada por sqlc para actualizar un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		db.UpdateEventParams{...}: Estructura con los parámetros para la consulta de actualización.
	*/
	err = s.queries.UpdateEvent(r.Context(), db.UpdateEventParams{
		ID:          int32(id),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		AgeRange:    req.AgeRange,
		UserID:      req.UserID,
	})
	if err != nil {
		http.Error(w, "Error al actualizar el evento", http.StatusInternalServerError)
		return
	}

	// No es necesario devolver el evento actualizado, solo un estado 200 OK
	// Por eso no se realiza json.NewEncoder(w).Encode(...)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener ID
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Convertir el ID a entero
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// 2. Eliminar usando sqlc
	/*
		s.queries.DeleteEvent: Llama a la función generada por sqlc para eliminar un evento específico.
		r.Context(): Pasa el contexto de la solicitud HTTP.
		int32(id): ID del evento a eliminar.
	*/
	err = s.queries.DeleteEvent(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el evento", http.StatusInternalServerError)
		return
	}

	// 3. Responder con estado 204 No Content
	w.WriteHeader(http.StatusNoContent)
}


//HANDLERS DE CATEGORY

// POST /categories: Crea una categoría nueva
func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Definimos una estructura temporal para guardar el JSON que envía el usuario
	var input struct {
		Name string `json:"name"`
	}

	// Leemos el cuerpo (Body) de la petición HTTP y lo guardamos en 'input'
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Cumplimos la consigna: validamos que el nombre no esté vacío
	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Llamamos al "cocinero" (sqlc) para que ejecute el INSERT en PostgreSQL
	category, err := s.queries.CreateCategory(r.Context(), input.Name)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// Devolvemos el status 201 (Created) y el objeto recién guardado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GET /categories: Trae todas las categorías
func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) {
	// Ejecuta el SELECT de todas las categorías
	categories, err := s.queries.ListCategories(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las categorías", http.StatusInternalServerError)
		return
	}

	// Si la tabla está vacía, evitamos devolver 'null' y devolvemos un array vacío '[]'
	if categories == nil {
		categories = []sqlc.Category{}
	}

	// Empaquetamos la lista en JSON y la enviamos al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GET /categories/{id}: Trae una categoría específica
func (s *Server) GetCategory(w http.ResponseWriter, r *http.Request) {
	// Extraemos el número que el usuario puso en la URL (ej: /categories/5)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32) // Convertimos ese texto '5' al número entero 5
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Buscamos el registro en la base de datos por su ID
	category, err := s.queries.GetCategory(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			// Cumplimos la consigna: Si no existe, devolvemos error 404
			http.Error(w, "Categoría no encontrada", http.StatusNotFound) 
			return
		}
		http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// PUT /categories/{id}: Modifica una categoría
func (s *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Obtenemos y convertimos el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Leemos los nuevos datos enviados en el JSON
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Validamos que no nos manden un nombre vacío
	if input.Name == "" {
		http.Error(w, "El nombre de la categoría no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Ejecutamos el UPDATE usando la estructura de parámetros que armó sqlc
	err = s.queries.UpdateCategory(r.Context(), sqlc.UpdateCategoryParams{
		ID:   int32(id),
		Name: input.Name,
	})
	
	if err != nil {
		http.Error(w, "Error al actualizar la categoría", http.StatusInternalServerError)
		return
	}

	// Respondemos que todo salió bien (Status 200)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Categoría actualizada con éxito"}`))
}

// DELETE /categories/{id}: Elimina una categoría
func (s *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Obtenemos y convertimos el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Ejecutamos el DELETE en la base de datos
	err = s.queries.DeleteCategory(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar la categoría", http.StatusInternalServerError)
		return
	}

	// Cumplimos la consigna: Devolver estado 204 (No Content)
	w.WriteHeader(http.StatusNoContent) 
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
	//CATEGORIES
    mux.HandleFunc("POST /categories", server.CreateCategory)
	mux.HandleFunc("GET /categories", server.ListCategories)
	mux.HandleFunc("GET /categories/{id}", server.GetCategory)
	mux.HandleFunc("PUT /categories/{id}", server.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", server.DeleteCategory)
	
	mux.HandleFunc("/", server.homeHandler)
	mux.HandleFunc("/events", server.CreateEventHandler)
	mux.HandleFunc("/events", server.ListEventsHandler)
	mux.HandleFunc("/events/{id}", server.GetEventHandler)
	mux.HandleFunc("/events/{id}", server.UpdateEventHandler)
	mux.HandleFunc("/events/{id}", server.DeleteEventHandler)
	mux.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("Servidor iniciado en http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
