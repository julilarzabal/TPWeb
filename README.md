# TPWeb
Luego de clonar el repositorio, ejecutar el siguiente comando:  _make run_

Es necesario tener las siguientes herramientas: _docker-compose-v2_ y _templ v0.3.1100_

Antes de ejecutar la aplicación, configurar `.env` en la raíz con el comando: _cp .env.example .env_

### Documentación: 
Creamos 3 tablas principales que son: _User_, _Category_ y _Event_.

La tabla _User_ va a corresponder a todos los tipos de usuarios de la aplicación (usuarios normales, administradores y colaboradores), el tipo será identificado por su atributo _role_.

Al momento de crear un usuario, siempre se comienza como usuario normal, y el rol de administrador o colaborador es asignado por nosotros (superAdministrador). 

La tabla _Category_ será usada tanto por los eventos (categoría/s a la/s cual/es corresponde el evento, por ejemplo taller, fiesta, etc) como por los usuarios (representa las preferencias del usuarios, lo que el usuario quiere ver en la aplicación). 

Las relaciones entre las tablas indicadas requieren 2 tablas adicionales, la tabla _user_preference_ (relación N:M entre User y Category) y _event_category_ (relación N:M entre Event y Category)

Por último, la tabla _Event_ representa a cada evento por sí mismo y cada uno está asociado a un usuario, el cual debe ser "organizador".

En esta entrega en particular abandonamos el enfoque de API + Cliente JS, creado en el tp4, para integrar la vista directamente en el servidor Go.
Para ello, creamos los componentes .templ, refactorizamos el servidor Go, con la creación de un nuevo handler principal encargado de renderizar el componente completo, y también
tuvimos que modificar los siguientes handlers, creados anteriormente, para el manejo de los formularios: CreateUserHandler, CreateEventHandler y CreateCategoryHandler. 
