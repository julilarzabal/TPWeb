package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// DTO para recibir el JSON en POST /users
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      string `json:"age"`      // Formato esperado: "YYYY-MM-DD"
	Industry string `json:"industry"` // Opcional
	Role     string `json:"role"`
}

// DTO para recibir el JSON en PUT /users/{id}
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      string `json:"age"`
	Industry string `json:"industry"`
	Role     string `json:"role"`
}

//sirve porque definimos estructuras intermedias (DTOs) que reciban los datos en formato de texto plano (strings)
//  No exponemos db.User directamente porque sql.NullString / sql.NullTime se serializan como {"String":"x","Valid":true}

// Handler para crear un nuevo usuario
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el JSON recibido en la petición
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "JSON de entrada inválido")
		return
	}

	// 2. Validación de regla de negocio
	if req.Name == "" || req.Email == "" || req.Role == "" {
		s.respondWithError(w, http.StatusBadRequest, "Los campos 'name', 'email' y 'role' son obligatorios")
		return
	}

	// 3. Conversión de datos (String -> sql.NullTime / sql.NullString)
	var ageNull sql.NullTime
	if req.Age != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Age)
		if err == nil {
			ageNull = sql.NullTime{Time: parsedDate, Valid: true}
		}
	}

	industryNull := sql.NullString{
		String: req.Industry,
		Valid:  req.Industry != "",
	}

	// 4. Guardar en la base de datos mediante sqlc
	user, err := s.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Age:      ageNull,
		Industry: industryNull,
		Role:     req.Role,
	})
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Error al guardar el usuario en la base de datos")
		return
	}

	// 5. Responder con estado 201 Created y el usuario recién creado
	s.respondWithJSON(w, http.StatusCreated, user)
}

// handler para listar los usuarios
func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.queries.ListUsers(r.Context())
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Error al consultar la lista de usuarios")
		return
	}

	// Asegurar que si no hay registros retorne un arreglo vacío [] y no null
	if users == nil {
		users = []sqlc.User{}
	}

	s.respondWithJSON(w, http.StatusOK, users)
}

// handler para obtener un usuario por ID
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "El ID de usuario debe ser un número entero válido")
		return
	}

	// 2. Buscar en la base de datos
	user, err := s.queries.GetUser(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			s.respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		s.respondWithError(w, http.StatusInternalServerError, "Error en el servidor al buscar el usuario")
		return
	}

	s.respondWithJSON(w, http.StatusOK, user)
}

// handler para actualizar un usuario
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "El ID de usuario debe ser un número entero válido")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "JSON de entrada inválido")
		return
	}

	if req.Name == "" || req.Email == "" || req.Role == "" {
		s.respondWithError(w, http.StatusBadRequest, "Los campos 'name', 'email' y 'role' son obligatorios")
		return
	}

	var ageNull sql.NullTime
	if req.Age != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Age)
		if err == nil {
			ageNull = sql.NullTime{Time: parsedDate, Valid: true}
		}
	}

	industryNull := sql.NullString{
		String: req.Industry,
		Valid:  req.Industry != "",
	}

	err = s.queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		ID:       int32(id),
		Name:     req.Name,
		Email:    req.Email,
		Age:      ageNull,
		Industry: industryNull,
		Role:     req.Role,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		s.respondWithError(w, http.StatusInternalServerError, "Error al actualizar el usuario")
		return
	}

	// Devolver el objeto actualizado
	updatedUser, _ := s.queries.GetUser(r.Context(), int32(id))
	s.respondWithJSON(w, http.StatusOK, updatedUser)
}

// handler para eliminar un usuario por id
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "El ID de usuario debe ser un número entero válido")
		return
	}

	err = s.queries.DeleteUser(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			s.respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		s.respondWithError(w, http.StatusInternalServerError, "Error al eliminar el usuario")
		return
	}

	// 204 No Content indica que la acción fue exitosa y no hay contenido para retornar
	w.WriteHeader(http.StatusNoContent)
}

// Métodos auxiliares de Server para responder en formato JSON
func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, map[string]string{"error": message})
}

// REVISAR EL TIPO DE RESPUESTA QUE DEVUELVE, EL FORMATO ES DISTINTO EN LO DE EVENT Y ESTO (aca uso respondWWithJSON, en lo de event SE USA json.NewEncoder(w).Encode(event))
