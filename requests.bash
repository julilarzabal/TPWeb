#!/usr/bin/env bash

echo "=== 1. POST /users - Creando un nuevo usuario ==="

# -i : permite que la terminal te muestre las cabeceras de respuesta que devuelve tu servidor, para verificar fácilmente el código de estado

curl -i -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juliana Pedone",
    "email": "juliana.pedone@example.com",
    "age": "2004-09-17",
    "industry": "Tecnología",
    "role": "Software Engineer"
  }'

echo -e "\n"


echo "=== 2. POST /users - Creando un nuevo usuario sin nombre ==="

curl -i -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "email": "juliana.larzabal@example.com",
    "age": "2005-04-27",
    "industry": "Tecnología",
    "role": "Software Engineer"
  }'

echo -e "\n"

echo "=== 3. POST /categories - Creando un nueva categoría ==="

curl -i -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Música y Festivales"
  }'

echo -e "\n"

echo "=== 4. POST /categories - Creando un nueva categoría ==="

curl -i -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deportes y Aventura"
  }'

echo -e "\n"

echo "=== 5. POST /user-preferences - Creando un nueva preferencia de usuario ==="

curl -i -X POST http://localhost:8080/user-preferences \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "category_id": 2
  }'

echo -e "\n"


echo "=== 6. POST /events - Creando un nuevo evento ==="

curl -i -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Conferencia de Go y PostgreSQL",
    "description": "Charla sobre arquitectura limpia, APIs REST y sqlc",
    "location": "Centro de Convenciones",
    "age_range": "18-65",
    "user_id": 1
  }'

echo -e "\n"

echo "=== 7. POST /events/1/categories - Agregando una categoría a un evento ==="

curl -i -X POST http://localhost:8080/events/6/categories \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 2
  }'

echo -e "\n"


echo "=== 8. GET /events - Listar todos los eventos ==="

curl -i -X GET http://localhost:8080/events

echo -e "\n\n"


echo "=== 9. GET /users/1 - Obtener usuario existente ==="

curl -i -X GET http://localhost:8080/users/1

echo -e "\n\n"


echo "=== 10. GET /users/9999 - Obtener usuario inexistente ==="

curl -i -X GET http://localhost:8080/users/9999

echo -e "\n\n"


echo "=== 11. PUT /users/1 - Actualizar usuario ==="

curl -i -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juliana Pedone",
    "email": "juliana.pedone@example.com",
    "age": "2004-09-17",
    "industry": "Tecnología",
    "role": "IT Engineer"
  }'

echo -e "\n\n"


echo "=== 12. PUT /users/1/preferences - Actualizar preferencias del usuario ==="

curl -i -X PUT http://localhost:8080/user-preferences/1 \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1
  }'

echo -e "\n\n"


echo "=== 13. DELETE /events/1/categories/1 - Eliminar una categoría de un evento ==="

curl -i -X DELETE http://localhost:8080/events/1/categories/2

echo -e "\n\n"


echo "=== 14. DELETE /events/1/categories - Eliminar todas las categorías del evento ==="

curl -i -X DELETE http://localhost:8080/events/1/categories

echo -e "\n\n"