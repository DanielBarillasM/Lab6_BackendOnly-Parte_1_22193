package main

// Importación de paquetes necesarios
import (
	"database/sql"           // Manejo de base de datos SQL
	"encoding/json"          // Codificación/decodificación JSON
	"fmt"                    // Para imprimir en consola
	"log"                    // Registro de errores
	"net/http"               // Manejo del servidor HTTP
	"strconv"                // Conversión de strings a enteros
	"strings"                // Manipulación de strings
	"time"                   // Control del tiempo (usado para retry)
	_ "github.com/go-sql-driver/mysql" // Driver MySQL (importación anónima)
)

// Estructura que representa un partido (match) en formato JSON
type Match struct {
	ID        int    `json:"id"`         // ID del partido
	HomeTeam  string `json:"homeTeam"`   // Equipo local
	AwayTeam  string `json:"awayTeam"`   // Equipo visitante
	MatchDate string `json:"matchDate"`  // Fecha del partido
}

// Variable global para conexión con la base de datos
var db *sql.DB

// Función principal del programa
func main() {
	initDB()             // Inicializa la base de datos
	defer db.Close()     // Cierra la conexión al finalizar

	// Define las rutas de la API y sus manejadores
	http.HandleFunc("/api/matches", handleMatches)
	http.HandleFunc("/api/matches/", routePatch)

	// Inicia el servidor en puerto 8080
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

// Función para inicializar conexión con MySQL, con reintentos automáticos
func initDB() {
	var err error
	dsn := "root:password@tcp(mysql:3306)/liga_db" // Cadena de conexión

	// Reintenta hasta 10 veces la conexión (esperando 2 segundos entre cada intento)
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			fmt.Println("✅ Conectado exitosamente a MySQL")
			return
		}
		fmt.Println("⏳ Esperando conexión a MySQL...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ No se pudo conectar a MySQL después de varios intentos:", err)
}

// Configura los encabezados CORS para permitir peticiones desde otros orígenes
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Manejador para el endpoint /api/matches
func handleMatches(w http.ResponseWriter, r *http.Request) {
	enableCors(w) // Habilita CORS
	if r.Method == http.MethodOptions {
		return // Responde sin contenido a opciones (preflight)
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Consulta todos los partidos
		rows, err := db.Query("SELECT id, home_team, away_team, match_date FROM matches")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var matches []Match
		for rows.Next() {
			var m Match
			if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			matches = append(matches, m)
		}
		if matches == nil {
			matches = []Match{}
		}
		json.NewEncoder(w).Encode(matches)

	case http.MethodPost:
		// Crea un nuevo partido
		var newMatch Match
		if err := json.NewDecoder(r.Body).Decode(&newMatch); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}
		result, err := db.Exec("INSERT INTO matches(home_team, away_team, match_date) VALUES (?, ?, ?)", newMatch.HomeTeam, newMatch.AwayTeam, newMatch.MatchDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, _ := result.LastInsertId()
		newMatch.ID = int(id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newMatch)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Manejador para endpoints con ID: GET, PUT, DELETE /api/matches/:id
func handleMatchByID(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/matches/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Obtiene un partido por ID
		var m Match
		err := db.QueryRow("SELECT id, home_team, away_team, match_date FROM matches WHERE id = ?", id).
			Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate)
		if err != nil {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(m)

	case http.MethodPut:
		// Actualiza un partido existente
		var updated Match
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("UPDATE matches SET home_team = ?, away_team = ?, match_date = ? WHERE id = ?", updated.HomeTeam, updated.AwayTeam, updated.MatchDate, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updated.ID = id
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		// Elimina un partido
		_, err := db.Exec("DELETE FROM matches WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error al eliminar", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Redirige las rutas PATCH a funciones específicas
func routePatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		// Verifica a qué tipo de PATCH corresponde
		if strings.HasSuffix(r.URL.Path, "/goals") {
			patchIncrement(w, r, "goals")
		} else if strings.HasSuffix(r.URL.Path, "/yellowcards") {
			patchIncrement(w, r, "yellow_cards")
		} else if strings.HasSuffix(r.URL.Path, "/redcards") {
			patchIncrement(w, r, "red_cards")
		} else if strings.HasSuffix(r.URL.Path, "/extratime") {
			patchExtraTime(w, r)
		} else {
			http.Error(w, "Ruta PATCH no válida", http.StatusNotFound)
		}
	} else {
		handleMatchByID(w, r) // Si no es PATCH, delega a GET/PUT/DELETE
	}
}

// Incrementa una columna numérica del partido (goles, tarjetas, etc.)
func patchIncrement(w http.ResponseWriter, r *http.Request, column string) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return
	}
	idStr := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Construye consulta SQL dinámica
	query := fmt.Sprintf("UPDATE matches SET %s = %s + 1 WHERE id = ?", column, column)
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error al actualizar "+column, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": column + " actualizado"})
}

// Establece el tiempo extra del partido
func patchExtraTime(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return
	}
	idStr := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Valor por defecto si no se especifica
	defaultExtra := "5 minutos"

	// Decodifica el JSON recibido
	var data struct {
		ExtraTime string `json:"extraTime"`
	}
	_ = json.NewDecoder(r.Body).Decode(&data)

	extra := data.ExtraTime
	if extra == "" {
		extra = defaultExtra
	}

	// Actualiza la columna extra_time
	_, err = db.Exec("UPDATE matches SET extra_time = ? WHERE id = ?", extra, id)
	if err != nil {
		http.Error(w, "Error al actualizar tiempo extra", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tiempo extra actualizado", "extraTime": extra})
}