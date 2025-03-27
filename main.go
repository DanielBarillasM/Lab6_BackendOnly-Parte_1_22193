package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

var db *sql.DB

func main() {
	initDB()

	defer db.Close()

	http.HandleFunc("/api/matches", handleMatches)
	http.HandleFunc("/api/matches/", handleMatchByID)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

// 🔁 Reintento para conectar a MySQL
func initDB() {
	var err error
	dsn := "root:password@tcp(mysql:3306)/liga_db"

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

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleMatches(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
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

		// ✅ Asegurarse de que no sea nil
		if matches == nil {
			matches = []Match{}
		}

		// 🐛 Imprimir en consola para depurar
		fmt.Println("📦 Matches desde la base de datos:", matches)

		json.NewEncoder(w).Encode(matches)

	case http.MethodPost:
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
		var m Match
		err := db.QueryRow("SELECT id, home_team, away_team, match_date FROM matches WHERE id = ?", id).Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate)
		if err != nil {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(m)

	case http.MethodPut:
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