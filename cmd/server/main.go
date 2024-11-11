package main

import (
	"log"
	"net/http"
	"stock/config"
	"stock/internal/api"
	"stock/internal/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	dbConn, err := db.InitDB(cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	handler := api.Hdl_fetch{DB: dbConn}
	http.HandleFunc("/api/fetch_stocks", handler.FetchAndStore_stocks)
	http.HandleFunc("/api/fetch_newQtRpt", handler.FetchAndStore_newQtRpt)
	http.HandleFunc("/api/fetch_theme", handler.FetchAndStore_theme)
	http.HandleFunc("/api/fetch_dividend", handler.FetchAndStore_dividend)
	http.HandleFunc("/api/fetch_history", handler.FetchAndStore_history)
	http.HandleFunc("/api/fetch_balance", handler.FetchAndStore_balance)
	http.HandleFunc("/api/fetch_income", handler.FetchAndStore_income)
	http.HandleFunc("/api/fetch_cashflow", handler.FetchAndStore_cashflow)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
