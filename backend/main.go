package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Все секретные переменные читаются из окружения БЕЗ дефолтов.
	// Если их нет — приложение падает и пишет в лог, что не так.
	dbUser := mustEnv("DB_USER")
	dbPassword := mustEnv("DB_PASSWORD")
	dbName := mustEnv("DB_NAME")
	dbHost := env("DB_HOST", "db")
	dbPort := env("DB_PORT", "5432")

	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ждём, пока БД поднимется
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			log.Println("connected to db")
			break
		}
		log.Printf("waiting for db... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	staticDir := env("STATIC_DIR", "/app/static")
	fs := http.FileServer(http.Dir(staticDir))

	// API
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		items := []Item{}
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			items = append(items, it)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	// Статика. Здесь ИСПРАВЛЕН баг с бесконечным редиректом:
	// на "/" отдаём index.html напрямую через ServeFile.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	addr := ":" + env("HTTP_PORT", "8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

// env — читает переменную с дефолтом (для несекретных значений: хост, порт, путь)
func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// mustEnv — читает переменную и падает, если её нет (для секретов)
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env %s is not set", key)
	}
	return v
}