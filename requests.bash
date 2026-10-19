#!/usr/bin/env bash

echo "=== Creando un nuevo usuario ==="

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


echo "=== Creando un nuevo usuario sin nombre ==="

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

echo "=== Creando un nueva categoría ==="

curl -i -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Música y Festivales"
  }'

echo -e "\n"

echo "=== Creando un nueva categoría ==="

curl -i -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deportes y Aventura"
  }'

echo -e "\n"

echo "=== Creando un nueva preferencia de usuario ==="

curl -i -X POST http://localhost:8080/users/1/preferences \
  -H "Content-Type: application/json" \
  -d '{
    "category_ids": [2]
  }'

echo -e "\n"


echo "=== GET /events - Listar todos los eventos ==="

curl -i -X GET http://localhost:8080/events

echo -e "\n\n"


echo "=== GET /users/1 - Obtener usuario existente ==="

curl -i -X GET http://localhost:8080/users/1

echo -e "\n\n"


echo "=== GET /users/9999 - Obtener usuario inexistente ==="

curl -i -X GET http://localhost:8080/users/9999

echo -e "\n\n"


echo "=== PUT /users/1 - Actualizar usuario ==="

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


echo "=== PUT /users/1/preferences - Actualizar preferencias del usuario ==="

curl -i -X PUT http://localhost:8080/user-preferences/1 \
  -H "Content-Type: application/json" \
  -d '{
    "category_ids": [1]
  }'

echo -e "\n\n"


echo "=== DELETE /events/1/categories/2 - Eliminar una categoría de un evento ==="

curl -i -X DELETE http://localhost:8080/events/1/categories/2

echo -e "\n\n"


echo "=== DELETE /events/1/categories - Eliminar todas las categorías del evento ==="

curl -i -X DELETE http://localhost:8080/events/1/categories

echo -e "\n\n"