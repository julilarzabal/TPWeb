package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	sqlc "TPEspecial.com/ProgramacionWeb/TPEspecial/db/sqlc"
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

// Handler para crear un nuevo usuario
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el JSON recibido en la petición
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON de entrada inválido", http.StatusBadRequest)
		return
	}

	// 2. Validación de regla de negocio
	if req.Name == "" || req.Email == "" || req.Role == "" {
		http.Error(w, "Los campos 'name', 'email' y 'role' son obligatorios", http.StatusBadRequest)
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
		http.Error(w, "Error al guardar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}

	// 5. Responder con estado 201 Created y el usuario recién creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Handler para listar los usuarios
func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Error al consultar la lista de usuarios", http.StatusInternalServerError)
		return
	}

	// Asegurar que si no hay registros retorne un arreglo vacío [] y no null
	if users == nil {
		users = []sqlc.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Handler para obtener un usuario por ID
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID de usuario debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	// 2. Buscar en la base de datos
	user, err := s.queries.GetUser(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en el servidor al buscar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handler para actualizar un usuario
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID de usuario debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON de entrada inválido", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Role == "" {
		http.Error(w, "Los campos 'name', 'email' y 'role' son obligatorios", http.StatusBadRequest)
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
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	updatedUser, _ := s.queries.GetUser(r.Context(), int32(id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// Handler para eliminar un usuario por id
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El ID de usuario debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	err = s.queries.DeleteUser(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	// 204 No Content indica éxito sin cuerpo
	w.WriteHeader(http.StatusNoContent)
}
