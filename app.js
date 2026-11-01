// ---------------------------------------------------------------
// Configuración
// ---------------------------------------------------------------
// Dirección donde corre tu API en Go. Cambiala si usás otro puerto.
const API_URL = "http://localhost:8080";

// NOTA: los nombres de los campos JSON (title, description, location,
// age_range, creator_id, fecha, name, email, age, industry, role,
// category_id) coinciden con db/schema/schema.sql.

// Acá guardamos los datos que trae la API para poder dibujarlos
let eventos = [];
let categorias = [];
let usuarios = [];
let busqueda = "";

// ---------------------------------------------------------------
// Funciones de ayuda
// ---------------------------------------------------------------

// Hace un fetch a la API y devuelve el JSON (o null si no hay cuerpo)
async function pedir(ruta, opciones) {
  const respuesta = await fetch(API_URL + ruta, opciones);
  if (!respuesta.ok) {
    throw new Error("La API respondió con error " + respuesta.status);
  }
  const texto = await respuesta.text();
  return texto ? JSON.parse(texto) : null;
}

// Manda datos en formato JSON (para POST)
function enviarJSON(ruta, datos) {
  return pedir(ruta, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(datos),
  });
}

// Evita que un texto escrito por alguien se interprete como HTML
function escapar(texto) {
  const div = document.createElement("div");
  div.textContent = texto ?? "";
  return div.innerHTML;
}

function mostrarMensaje(texto) {
  document.getElementById("mensaje").textContent = texto;
}

// ---------------------------------------------------------------
// Cargar datos desde la API (GET)
// ---------------------------------------------------------------
async function cargarTodo() {
  try {
    // Pedimos las tres listas al mismo tiempo.
    // Si la API devuelve null (lista vacía), usamos [].
    const [e, c, u] = await Promise.all([
      pedir("/events"),
      pedir("/categories"),
      pedir("/users"),
    ]);
    eventos = e || [];
    categorias = c || [];
    usuarios = u || [];

    // A cada evento le pedimos sus categorías
    await Promise.all(
      eventos.map(async (evento) => {
        try {
          evento.categorias = (await pedir(`/events/${evento.id}/categories`)) || [];
        } catch {
          evento.categorias = [];
        }
      })
    );

    dibujarEventos();
    dibujarCategorias();
    dibujarUsuarios();
    mostrarMensaje("");
  } catch (error) {
    mostrarMensaje("No se pudo conectar con la API. ¿Está encendida? " + error.message);
  }
}

// ---------------------------------------------------------------
// Dibujar en pantalla
// ---------------------------------------------------------------
function dibujarEventos() {
  const contenedor = document.getElementById("lista-eventos");
  const termino = busqueda.trim().toLowerCase();
  const visibles = termino
    ? eventos.filter((evento) =>
        [evento.title, evento.description, evento.location, evento.age_range]
          .join(" ")
          .toLowerCase()
          .includes(termino)
      )
    : eventos;

  if (visibles.length === 0) {
    const texto = eventos.length === 0
      ? "Todavía no hay eventos. Publicá el primero con el botón Create Event."
      : "No hay eventos que coincidan con tu búsqueda.";
    contenedor.innerHTML = `<p class="vacio">${texto}</p>`;
    return;
  }

  let html = "";
  for (const evento of visibles) {
    const fecha = new Date(evento.fecha).toLocaleDateString("es-AR", {
      dateStyle: "full",
      timeZone: "UTC",
    });
    const organizador = usuarios.find((u) => u.id === evento.creator_id);
    const etiquetas = (evento.categorias || [])
      .map((cat) => `<span class="etiqueta">${escapar(cat.name)}</span>`)
      .join("");

    html += `
      <article class="evento">
        <div class="etiquetas">${etiquetas}</div>
        <h3>${escapar(evento.title)}</h3>
        <p class="fecha">${escapar(fecha)}</p>
        <p class="rango">Rango de edad: ${escapar(evento.age_range)}</p>
        <p>${escapar(evento.description)}</p>
        <p class="lugar">${escapar(evento.location)}${organizador ? " · Organiza " + escapar(organizador.name) : ""}</p>
        <button class="borrar" data-tipo="events" data-id="${evento.id}">Eliminar</button>
      </article>`;
  }
  contenedor.innerHTML = html;
}

function dibujarCategorias() {
  // Lista de categorías con botón de borrar
  const lista = document.getElementById("lista-categorias");
  lista.innerHTML = categorias
    .map(
      (cat) => `
      <li>
        <span>${escapar(cat.name)}</span>
        <button class="borrar" data-tipo="categories" data-id="${cat.id}">Eliminar</button>
      </li>`
    )
    .join("");

  // Casillas para elegir categorías al crear un evento
  const casillas = document.getElementById("evento-categorias");
  casillas.innerHTML = categorias
    .map(
      (cat) => `
      <label class="casilla">
        <input type="checkbox" value="${cat.id}"> ${escapar(cat.name)}
      </label>`
    )
    .join("");
}

function dibujarUsuarios() {
  const lista = document.getElementById("lista-usuarios");
  lista.innerHTML = usuarios
    .map(
      (u) => `
      <li>
        <span>${escapar(u.name)} <small>${u.role === "creator" ? "(organizador)" : ""}</small>
          <small>${escapar(u.industry)}${u.age ? " · " + escapar(u.age.slice(0, 10)) : ""}</small></span>
        <button class="borrar" data-tipo="users" data-id="${u.id}">Eliminar</button>
      </li>`
    )
    .join("");

  // Opciones del selector de organizador (solo los usuarios organizadores)
  const selector = document.getElementById("evento-creador");
  const organizadores = usuarios.filter((u) => u.role === "creator");
  selector.innerHTML = organizadores.length
    ? organizadores.map((u) => `<option value="${u.id}">${escapar(u.name)}</option>`).join("")
    : "<option value=''>Primero creá un organizador</option>";
}

// ---------------------------------------------------------------
// Crear (POST)
// ---------------------------------------------------------------
document.getElementById("form-evento").addEventListener("submit", async (evento) => {
  evento.preventDefault();
  try {
    const nuevo = await enviarJSON("/events", {
      title: document.getElementById("evento-titulo").value,
      description: document.getElementById("evento-descripcion").value,
      fecha: document.getElementById("evento-fecha").value,
      location: document.getElementById("evento-lugar").value,
      age_range: document.getElementById("evento-rango-edad").value,
      creator_id: Number(document.getElementById("evento-creador").value),
    });

    // Asociamos las categorías tildadas al evento recién creado
    const tildadas = document.querySelectorAll("#evento-categorias input:checked");
    for (const casilla of tildadas) {
      await enviarJSON(`/events/${nuevo.id}/categories`, {
        category_id: Number(casilla.value),
      });
    }

    evento.target.reset();
    evento.target.closest("dialog")?.close();
    await cargarTodo(); // volvemos a pedir la lista completa
  } catch (error) {
    mostrarMensaje("No se pudo crear el evento: " + error.message);
  }
});

document.getElementById("form-categoria").addEventListener("submit", async (evento) => {
  evento.preventDefault();
  try {
    await enviarJSON("/categories", {
      name: document.getElementById("categoria-nombre").value,
    });
    evento.target.reset();
    evento.target.closest("dialog")?.close();
    await cargarTodo();
  } catch (error) {
    mostrarMensaje("No se pudo crear la categoría: " + error.message);
  }
});

document.getElementById("form-usuario").addEventListener("submit", async (evento) => {
  evento.preventDefault();
  try {
    await enviarJSON("/users", {
      name: document.getElementById("usuario-nombre").value,
      email: document.getElementById("usuario-email").value,
      age: document.getElementById("usuario-edad").value,
      industry: document.getElementById("usuario-industria").value,
      role: document.getElementById("usuario-rol").value,
    });
    evento.target.reset();
    evento.target.closest("dialog")?.close();
    await cargarTodo();
  } catch (error) {
    mostrarMensaje("No se pudo crear el usuario: " + error.message);
  }
});

// ---------------------------------------------------------------
// Búsqueda, carrusel y diálogos
// ---------------------------------------------------------------
document.getElementById("buscador").addEventListener("input", (evento) => {
  busqueda = evento.target.value;
  dibujarEventos();
});

const pistaEventos = document.getElementById("lista-eventos");
document.getElementById("carrusel-prev").addEventListener("click", () => {
  pistaEventos.scrollBy({ left: -340, behavior: "smooth" });
});
document.getElementById("carrusel-next").addEventListener("click", () => {
  pistaEventos.scrollBy({ left: 340, behavior: "smooth" });
});

document.querySelectorAll("[data-abrir-dialogo]").forEach((boton) => {
  boton.addEventListener("click", () => {
    document.getElementById(boton.dataset.abrirDialogo).showModal();
  });
});

document.querySelectorAll("[data-cerrar-dialogo]").forEach((boton) => {
  boton.addEventListener("click", () => {
    boton.closest("dialog").close();
  });
});

document.querySelectorAll("dialog").forEach((dialogo) => {
  dialogo.addEventListener("click", (evento) => {
    if (evento.target === dialogo) dialogo.close();
  });
});

// ---------------------------------------------------------------
// Eliminar (DELETE) con delegación de eventos
// ---------------------------------------------------------------
// Los botones se crean dinámicamente, así que escuchamos los clics
// en todo el documento y vemos si el clic fue sobre un botón "borrar".
document.addEventListener("click", async (clic) => {
  const boton = clic.target.closest(".borrar");
  if (!boton) return;

  const tipo = boton.dataset.tipo; // "events", "categories" o "users"
  const id = boton.dataset.id;

  if (!confirm("¿Seguro que querés eliminarlo?")) return;

  try {
    await pedir(`/${tipo}/${id}`, { method: "DELETE" });
    await cargarTodo(); // refrescamos la vista
  } catch (error) {
    mostrarMensaje("No se pudo eliminar: " + error.message);
  }
});

// ---------------------------------------------------------------
// Al cargar la página
// ---------------------------------------------------------------
cargarTodo();
