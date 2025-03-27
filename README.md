# ⚽ La Liga Tracker - Backend

Este proyecto implementa el backend de la aplicación web **La Liga Tracker**, desarrollado en **Go** con conexión a base de datos **MySQL**, desplegado con **Docker Compose**. Permite gestionar partidos de fútbol de La Liga mediante una API RESTful.

---

## 📦 Funcionalidad RESTful

| Método  | Endpoint                  | Descripción                          |
|---------|---------------------------|--------------------------------------|
| GET     | `/api/matches`            | Obtener todos los partidos           |
| GET     | `/api/matches/:id`        | Obtener un partido por su ID         |
| POST    | `/api/matches`            | Crear un nuevo partido               |
| PUT     | `/api/matches/:id`        | Actualizar un partido existente      |
| DELETE  | `/api/matches/:id`        | Eliminar un partido por su ID        |

---

## 🐳 Uso con Docker y Docker Compose

### 📁 Estructura del proyecto

```bash
├── main.go               # Código fuente del backend en Go
├── Dockerfile            # Imagen del backend (puerto 8080)
├── docker-compose.yml    # Orquestador de servicios (MySQL + Backend)
├── init.sql              # Script para inicializar tabla en MySQL
├── LaLigaTracker.html    # Frontend proporcionado
├── screenshot.png        # Captura del frontend funcionando
├── go.mod / go.sum       # Dependencias del proyecto
└── README.md             # Este documento
```

### ▶️ Cómo ejecutar el proyecto

```bash
# 1. Construir y levantar los servicios
docker compose up --build

# 2. Backend disponible en:
http://localhost:8080

# 3. Abrir el frontend (HTML local)
http://localhost:5500/LaLigaTracker.html
```

---

## 🛠️ Detalles técnicos

- Lenguaje: **Go**
- DB: **MySQL 8.0**
- Driver usado: `github.com/go-sql-driver/mysql`
- Persistencia completa en base de datos (`init.sql` crea la tabla `matches`)
- El backend se reconecta automáticamente a la base de datos si aún se está iniciando.

---

## 📸 Captura del Frontend funcionando

![alt text](image.png)

---

## 👨‍💻 Autor

- **Nombre:** Pablo Daniel Barillas Moreno 
- **Carné:** 22193

--- 
