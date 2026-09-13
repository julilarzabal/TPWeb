# TPWeb
Luego de clonar el repositorio, ejecutar el siguiente comando:  _make test_

### Documentación: 
Creamos 3 tablas principales que son: _User_, _Category_ y _Event_.

La tabla _User_ va a corresponder a todos los tipos de usuarios de la aplicación (usuarios normales, administradores y colaboradores), el tipo será identificado por su atributo _role_.

Al momento de crear un usuario, siempre se comienza como usuario normal, y el rol de administrador o colaborador es asignado por nosotros (superAdministrador). 

La tabla _Category_ será usada tanto por los eventos (categoría/s a la/s cual/es corresponde el evento, por ejemplo taller, fiesta, etc) como por los usuarios (representa las preferencias del usuarios, lo que el usuario quiere ver en la aplicación). 

Las relaciones entre las tablas indicadas requieren 2 tablas adicionales, la tabla _user_preference_ (relación N:M entre User y Category) y _event_category_ (relación N:M entre Event y Category)

Por último, la tabla _Event_ representa a cada evento por sí mismo y cada uno está asociado a un usuario, el cual (en un futuro) deberá ser colaborador.
