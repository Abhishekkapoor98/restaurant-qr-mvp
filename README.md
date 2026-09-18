# 🍽️ Restaurant QR Ordering & Live Kitchen Dashboard

A real-time, zero-friction table ordering MVP built for speed and seamless restaurant operations. Customers scan a table-specific QR code to access a digital menu and place orders, which instantly trigger a live alert on the kitchen's display without a single page refresh.

## 🚀 Key Features

* **Zero-Friction Ordering:** Customers access the menu via URL parameters (e.g., `/?table=5`), which automatically locks their table number—no manual entry required.
* **Real-Time Kitchen Sync (WebSockets):** Orders are broadcasted instantly to the kitchen dashboard via `gorilla/websocket`.
* **Database Auto-Migration & Seeding:** The Go server automatically verifies, creates necessary PostgreSQL tables, and injects dummy menu items on startup.
* **Clean Architecture:** Strictly separated layers (Domain, Repository, Usecase, Delivery) for high maintainability and scalable business logic.
* **Mobile-First UI:** Built with Tailwind CSS, ensuring the customer menu looks perfect on phones while the kitchen dashboard utilizes wide screens effectively.

## 🛠️ Tech Stack

* **Backend:** Go (Golang), Gin Web Framework
* **Real-Time Engine:** Gorilla WebSockets
* **Database:** PostgreSQL (`lib/pq`)
* **Frontend:** Server-Side Rendered HTML, Vanilla JavaScript, Tailwind CSS (CDN)
* **Architecture:** Clean Architecture Pattern

## 📂 Project Structure

├── cmd/
│   └── main.go              # Application entry point, DB setup, Auto-migration
├── delivery/
│   └── http_handler.go      # Gin HTTP routes, WebSocket upgrader, and logic
├── domain/
│   └── menu.go              # Core structs (MenuItem, Order) and Interfaces
├── repository/
│   └── pg_menu_repo.go      # PostgreSQL database queries
├── usecase/
│   └── menu_usecase.go      # Business logic and validation
├── templates/
│   ├── menu.html            # Customer-facing digital menu UI
│   └── kitchen.html         # Live kitchen dashboard UI
├── go.mod
└── go.sum

## 💻 Local Setup & Testing

### 1. Start Local PostgreSQL Database (Docker)
Ensure Docker is running, then create and start the database:
```bash
docker run --name postgres-restaurant -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15-alpine
docker exec -it postgres-restaurant psql -U postgres -c "CREATE DATABASE restaurantdb;"