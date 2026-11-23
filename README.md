# TPWeb
Luego de clonar el repositorio, ejecutar el siguiente comando:  _make run_

Es necesario tener las siguientes herramientas: _Go 1.25_, _docker-compose-v2_, _sqlc_ y _templ v0.3.1020_ (el Makefile lo ejecuta con `go run`, por lo que no requiere instalación aparte)

Antes de ejecutar la aplicación, configurar `.env` en la raíz con el comando: _cp .env.example .env_

### Documentación: 

#### 1. Introducción

**Tandil App** es una aplicación web que funciona como agenda cultural de la ciudad de Tandil: permite publicar eventos (talleres, música, ferias, paseos), organizarlos por categorías y registrar usuarios, algunos de los cuales actúan como organizadores de eventos.

La aplicación se renderiza íntegramente en el servidor: un servidor Go genera el HTML a partir de componentes **templ** con datos obtenidos de PostgreSQL mediante código generado por **sqlc**. En el navegador, **HTMX** envía las acciones del usuario por AJAX y reemplaza solo los fragmentos de la página que cambian, por lo que crear, editar o eliminar elementos no recarga la página completa.

Las diferencias respecto del TP5 se detallan en la [sección 10](#10-actualizaciones-del-tp5-al-tp6).

#### 2. Objetivos

**Objetivo general:** ofrecer una agenda cultural renderizada en el servidor que se comporte como una aplicación de página única (SPA), actualizando fragmentos HTML con HTMX en lugar de recargar la página completa y sin escribir JavaScript propio.

**Objetivos específicos:**
1. Integrar la librería HTMX en el layout de la aplicación.
2. Crear entidades (eventos y categorías) sin recargar la página, devolviendo desde el servidor solo el fragmento de la lista actualizada.
3. Eliminar entidades directamente desde la lista, quitando el elemento del DOM a partir de una respuesta vacía del servidor.
4. Mantener la separación de responsabilidades entre presentación, lógica y datos.

#### 3. Modelo de datos

Creamos 3 tablas principales que son: _User_, _Category_ y _Event_.

La tabla _User_ va a corresponder a todos los tipos de usuarios de la aplicación, el tipo será identificado por su atributo _role_. Actualmente la aplicación maneja dos roles: _user_ (usuario común) y _creator_ (organizador de eventos).

Al momento de registrarse desde el formulario _Sign Up_, el usuario elige su tipo: "Usuario común" u "Organizador de eventos". Solo los usuarios con rol _creator_ aparecen como opción de organizador al crear un evento. 

La tabla _Category_ será usada tanto por los eventos (categoría/s a la/s cual/es corresponde el evento, por ejemplo taller, fiesta, etc) como por los usuarios (representa las preferencias del usuarios, lo que el usuario quiere ver en la aplicación). 

Las relaciones entre las tablas indicadas requieren 2 tablas adicionales, la tabla _user_preference_ (relación N:M entre User y Category) y _event_category_ (relación N:M entre Event y Category)

Por último, la tabla _Event_ representa a cada evento por sí mismo y cada uno está asociado a un usuario, el cual debe ser "organizador".

| Tabla | Campos principales |
|---|---|
| `user` | `id`, `name`, `email` (único), `age` (fecha, opcional), `industry` (opcional), `role` |
| `category` | `id`, `name` |
| `event` | `id`, `title`, `description` (opcional), `location`, `age_range`, `creator_id`, `fecha` |
| `user_preference` | `user_id`, `category_id` |
| `event_category` | `event_id`, `category_id` |

```
user 1 ──── N event                (event.creator_id → user.id)
user N ──── M category             (vía user_preference)
event N ──── M category            (vía event_category)
```

Todas las claves foráneas se definen con **`ON DELETE CASCADE`**: al eliminar un evento o una categoría se eliminan automáticamente sus filas asociadas en las tablas intermedias, por lo que el borrado desde la interfaz no falla por restricciones de integridad.

#### 4. Arquitectura del sistema

La aplicación sigue una arquitectura **three‑tier**:

```
┌──────────────────────────┐
│   Navegador + HTMX       │  ← Presentación (cliente)
└────────────┬─────────────┘
             │ HTTP (página completa o fragmentos HTML)
┌────────────▼─────────────┐
│  Servidor Go (net/http)  │
│  ├─ handlers/            │  ← Lógica: valida, coordina, decide qué renderizar
│  └─ views/ (templ)       │  ← Presentación (servidor): componentes HTML
└────────────┬─────────────┘
             │ db/sqlc (código generado)
┌────────────▼─────────────┐
│      PostgreSQL          │  ← Datos
└──────────────────────────┘
```

- **Presentación:** componentes templ (`views/`) que el servidor convierte en HTML. HTMX, en el navegador, decide qué parte de la página reemplazar con cada respuesta.
- **Lógica:** los handlers (`handlers/`) leen el formulario, validan reglas de negocio, invocan la capa de datos y eligen qué componente devolver.
- **Datos:** consultas SQL (`db/queries/`) compiladas por sqlc a funciones Go (`db/sqlc/`).

**Inyección de dependencias:** `main.go` crea la conexión a la base, construye `sqlc.Queries` y lo inyecta en el `Server` de handlers mediante el constructor `handlers.NewServer(queries)`. Ningún handler abre conexiones por su cuenta.

#### 5. Estructura del proyecto

```
.
├── main.go                  # Configuración, conexión a la BD y registro de rutas
├── handlers/                # Lógica HTTP
│   ├── handler.go           # Server, NewServer, HandleIndex, renderFragment
│   ├── event_handlers.go
│   ├── category_handlers.go
│   ├── user_handlers.go
│   ├── event_category_handlers.go
│   └── user_preferences_handlers.go
├── views/                   # Componentes templ
│   ├── layout.templ         # Estructura HTML común (head, header, footer)
│   ├── IndexPage.templ      # Página principal
│   ├── entity_form.templ    # Formularios (usuario, evento, categoría, edición)
│   └── entity_list.templ    # Listas de eventos y categorías
├── db/
│   ├── schema/schema.sql    # Definición de tablas
│   ├── queries/queries.sql  # Consultas para sqlc
│   └── sqlc/                # Código generado por sqlc (no versionado)
├── static/                  # CSS (Pico + estilos propios)
├── tests/db_test.go         # Pruebas de integración contra la BD
├── requests.bash            # Peticiones curl de prueba (endpoints JSON)
├── docker-compose.yml       # Contenedor de PostgreSQL
├── sqlc.yaml                # Configuración de sqlc
├── Makefile                 # Automatiza generación, compilación y ejecución
└── .env.example             # Plantilla de configuración
```

Los archivos `*_templ.go` y `db/sqlc/*.go` se generan automáticamente (`templ generate` y `sqlc generate`) y no se versionan.

#### 6. Capa de acceso a datos (sqlc)

Las consultas se escriben en SQL plano en `db/queries/queries.sql` y sqlc genera funciones Go con tipos seguros. La configuración (`sqlc.yaml`) mapea las columnas opcionales a punteros (`*string`, `*time.Time`) para distinguir "sin valor" de "cadena vacía".

| Entidad | Consultas disponibles |
|---|---|
| Usuario | `CreateUser`, `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser` |
| Categoría | `CreateCategory`, `GetCategory`, `ListCategories`, `UpdateCategory`, `DeleteCategory` |
| Evento | `CreateEvent`, `GetEvent`, `ListEvents`, `UpdateEvent`, `DeleteEvent` |
| Preferencias | `ListUserPreferences`, `UpdateUserPreferences`, `AddUserPreference`, `GetPreferencesByUser`, `RemoveUserPreference` |
| Evento–Categoría | `AddEventCategory`, `GetCategoriesByEvent`, `RemoveEventCategory`, `ListEventCategories`, `DeleteEventCategories` |

#### 7. Servidor HTTP: rutas y handlers

| Método y ruta | Handler | Respuesta |
|---|---|---|
| `GET /` | `HandleIndex` | Página completa |
| `POST /users` | `CreateUserHandler` | Redirección 303 a `/` |
| `POST /categories` | `CreateCategoryHandler` | Fragmento `#category-list` |
| `POST /events` | `CreateEventHandler` | Fragmento `#event-list` |
| `DELETE /categories/{id}` | `DeleteCategoryHandler` | 200 OK, cuerpo vacío |
| `DELETE /events/{id}` | `DeleteEventHandler` | 200 OK, cuerpo vacío |
| `GET /events/{id}/edit` | `ShowEditEventFormHandler` | Fragmento: formulario de edición |
| `PUT /events/{id}` | `UpdateEventFormHandler` | Fragmento `#event-list` |
| `GET /static/` | `http.FileServer` | Archivos estáticos |

**Estructura común de un handler de formulario:**
1. `r.ParseForm()` para leer los campos enviados.
2. Validación de reglas de negocio (campos obligatorios, ID positivo, formato de fecha `YYYY-MM-DD`). Ante un error se responde **400 Bad Request**.
3. Llamada a la función de sqlc correspondiente. Ante un error de base de datos se responde **500 Internal Server Error**.
4. Respuesta: el fragmento HTML actualizado (excepto el registro de usuarios, que redirige a `/`).

**Función auxiliar `renderFragment`** (`handlers/handler.go`): centraliza la respuesta de fragmentos para HTMX. Establece `Content-Type: text/html; charset=utf-8` y renderiza el componente. Si el renderizado falla, el error se registra en el log en lugar de llamar a `http.Error`, porque parte de la respuesta podría haberse enviado ya.

#### 8. Capa de presentación (templ)

| Componente | Archivo | Descripción |
|---|---|---|
| `Layout` | `layout.templ` | Estructura HTML común: `<head>` (CSS y script de HTMX), cabecera con navegación y panel de registro, pie de página. |
| `IndexPage` | `IndexPage.templ` | Página principal: portada, `EventList` y `CategoryList`. |
| `EventList` / `CategoryList` | `entity_list.templ` | Sección completa: título, botón del panel de creación y la lista. |
| `EventItems` / `CategoryItems` | `entity_list.templ` | **Fragmento** que HTMX reemplaza. Siempre renderiza el contenedor `#event-list` / `#category-list`, incluso con la lista vacía. Incluye los botones de borrar (y de modificar en eventos). |
| `SignUpForm` / `SignUpPanel` | `entity_form.templ` | Registro de usuarios (envío tradicional). |
| `EventForm` / `EventPanel` | `entity_form.templ` | Creación de eventos con HTMX. |
| `EventEditForm` | `entity_form.templ` | Edición de un evento, precargado con sus datos. |
| `eventFormFields` | `entity_form.templ` | Campos compartidos por la creación y la edición de eventos. |
| `CategoryForm` / `CategoryPanel` | `entity_form.templ` | Creación de categorías con HTMX. |

#### 9. Interacción dinámica con HTMX

**Creación sin recarga:**

```html
<form id="form-categoria"
      hx-post="/categories"
      hx-target="#category-list"
      hx-swap="outerHTML"
      hx-on::after-request="if (event.detail.successful) this.closest('details').open = false">
```

| Atributo | Función |
|---|---|
| `hx-post` | Envía el formulario por AJAX a la ruta de creación de la entidad. |
| `hx-target` | Indica qué elemento se actualiza con la respuesta: el contenedor de la lista. |
| `hx-swap="outerHTML"` | Reemplaza el contenedor completo; como la respuesta trae su propio `<div id="...">`, no quedan elementos anidados ni ids duplicados. |
| `hx-on::after-request` | Si la petición fue exitosa, cierra el panel desplegable del formulario. |

La validación HTML5 (`required`) se mantiene: HTMX no envía el formulario si falta un campo obligatorio.

**Borrado sin recarga:**

```html
<button type="button" class="boton-eliminar"
        hx-delete="/events/{id}"
        hx-target="closest article"
        hx-swap="outerHTML"
        hx-confirm="¿Eliminar este evento?">Delete</button>
```

1. `hx-confirm` muestra un diálogo de confirmación.
2. `hx-delete` envía `DELETE /events/{id}`.
3. El servidor elimina el registro y responde **200 OK con cuerpo vacío**.
4. HTMX reemplaza el elemento objetivo (`closest article` en eventos, `closest li` en categorías) por la respuesta vacía, y el elemento desaparece de la vista.

**Edición de eventos (adicional):**
1. El botón **Modify** hace `hx-get="/events/{id}/edit"` con `hx-target="closest article"`: la tarjeta del evento se reemplaza por un formulario precargado.
2. El formulario envía `hx-put="/events/{id}"` con `hx-target="#event-list"`: el servidor actualiza el evento y devuelve la lista completa.

#### 10. Actualizaciones del TP5 al TP6

| Aspecto | TP5 | TP6 |
|---|---|---|
| Envío de formularios | `<form method="post" action="...">` tradicional | `hx-post` / `hx-put` (AJAX) |
| Respuesta del servidor al crear | Redirección 303 (patrón PRG) y recarga completa | Fragmento HTML con la lista actualizada |
| Actualización de la vista | Se vuelve a renderizar toda la página | Solo se reemplaza el fragmento afectado |
| Borrado desde la interfaz | No disponible | Botón por elemento con confirmación, sin recarga |
| Edición desde la interfaz | No disponible | Formulario en línea para eventos |
| Respuesta de los `DELETE` | 204 No Content | 200 OK con cuerpo vacío |
| Componentes de lista | La lista estaba dentro de la sección junto al formulario | La lista es un componente propio con un `id` fijo |
| JavaScript propio | Ninguno | Ninguno (solo atributos HTMX) |

#### 11. Conclusiones

HTMX permite que una aplicación renderizada íntegramente en el servidor ofrezca una interfaz que se siente como una SPA, sin escribir JavaScript propio. El servidor sigue siendo la única fuente de verdad: genera el HTML y la base de datos determina qué se muestra, mientras que HTMX solo decide **dónde** colocar cada fragmento.

La clave está en la organización de los componentes: separar las listas en fragmentos con identificadores fijos fue lo que hizo posible actualizar partes de la página de forma precisa. Las validaciones se mantienen en el servidor, el acceso a datos y las vistas son código tipado generado por sqlc y templ, y templ brinda protección contra XSS.

#### 12. Referencias

- Documentación de HTMX: <https://htmx.org/docs/>
- Ejemplos de HTMX (Out‑of‑Band Swaps): <https://htmx.org/examples/update-other-content/#oob>
- Documentación de templ: <https://templ.guide/>
- Documentación de sqlc: <https://docs.sqlc.dev/>
- Paquete `net/http` de Go: <https://pkg.go.dev/net/http>
- Material de la cátedra: *Programación Web — Diseño*, Dr. A. Zunino y Dr. A. Teyseyre.
- Enunciado: *Trabajo Práctico 6 — Interfaces Dinámicas con HTMX*.
