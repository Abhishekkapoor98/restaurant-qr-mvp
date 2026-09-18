package main

import (
	"database/sql"
	"log"
	"os"

	"restaurant-qr/delivery"
	"restaurant-qr/repository"
	"restaurant-qr/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Step 1: Database Setup
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		// M4 Mac ke liye 127.0.0.1 aur naya database naam: restaurantdb
		dbConnStr = "postgres://postgres:postgres@127.0.0.1:5432/restaurantdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Step 2: Auto-create tables on startup
	initSQL := `
	CREATE TABLE IF NOT EXISTS menu_items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		category VARCHAR(50) NOT NULL,
		image_url TEXT
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		table_number INTEGER NOT NULL,
		items TEXT NOT NULL,
		total_amount NUMERIC(10, 2) NOT NULL,
		status VARCHAR(20) DEFAULT 'pending'
	);
	`
	if _, err := db.Exec(initSQL); err != nil {
		log.Printf("Warning: Could not auto-create tables: %v", err)
	} else {
		log.Println("Database tables verified/created successfully!")
	}

	// Step 3: Seed Dummy Menu Items (Agar table khali hai)
	seedSQL := `
	INSERT INTO menu_items (name, price, category, image_url)
	SELECT 'Truffle Fries', 199.00, 'Starters', 'https://images.unsplash.com/photo-1576107222453-cefe537a82fc?auto=format&fit=crop&w=300&q=80'
	WHERE NOT EXISTS (SELECT 1 FROM menu_items LIMIT 1);
	
	INSERT INTO menu_items (name, price, category, image_url)
	SELECT 'Farmhouse Pizza', 349.00, 'Main Course', 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=300&q=80'
	WHERE NOT EXISTS (SELECT 1 FROM menu_items WHERE name = 'Farmhouse Pizza');
	`
	db.Exec(seedSQL)

	// Step 4: Initialize Layers (Dependency Injection)
	menuRepo := repository.NewPgMenuRepository(db)
	menuUC := usecase.NewMenuUsecase(menuRepo)

	// Step 5: Start Gin Web Framework
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	delivery.NewMenuHandler(r, menuUC)

	// Step 6: Get Port
	port := os.Getenv("PORT")
	if port == "" {
		// Dhyan do: Isko mene 8081 rakha hai taaki Project 1 (8080) ke sath clash na ho!
		port = "8081"
	}

	// Step 7: Start Server
	log.Printf("Restaurant Server is running at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
