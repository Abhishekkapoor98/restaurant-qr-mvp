package delivery

import (
	"log"
	"net/http"
	"restaurant-qr/domain"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader (HTTP request ko WebSocket connection me badalne ke liye)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // MVP ke liye sab allow kar rahe hain
	},
}

type MenuHandler struct {
	Usecase domain.MenuUsecase
	// Kitchen ke jitne bhi screens open honge, unka track rakhne ke liye
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
}

// Constructor
func NewMenuHandler(r *gin.Engine, us domain.MenuUsecase) {
	handler := &MenuHandler{
		Usecase: us,
		clients: make(map[*websocket.Conn]bool),
	}

	// Customer Routes
	r.GET("/", handler.RenderMenu)
	r.POST("/api/orders", handler.PlaceOrder)

	// Kitchen Routes
	r.GET("/kitchen", handler.RenderKitchen)
	r.GET("/ws", handler.HandleConnections) // WebSocket endpoint
	r.POST("/kitchen/complete/:id", handler.CompleteOrder)
}

// --- Customer Handlers ---

func (h *MenuHandler) RenderMenu(c *gin.Context) {
	menu, err := h.Usecase.FetchMenu(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Menu load karne me error aa gaya!")
		return
	}

	// URL se table number nikalo (e.g., scan karne par /?table=3 khulega)
	tableNum := c.Query("table")

	// Customer menu page render karo aur table number sath bhejo
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"menu":  menu,
		"table": tableNum,
	})
}

func (h *MenuHandler) PlaceOrder(c *gin.Context) {
	tableNum, _ := strconv.Atoi(c.PostForm("table_number"))
	totalAmount, _ := strconv.ParseFloat(c.PostForm("total_amount"), 64)

	order := &domain.Order{
		TableNumber: tableNum,
		Items:       c.PostForm("items"), // "2x Burger, 1x Coke"
		TotalAmount: totalAmount,
	}

	if err := h.Usecase.PlaceOrder(c.Request.Context(), order); err != nil {
		c.String(http.StatusBadRequest, "Order fail ho gaya. Details check karo.")
		return
	}

	// YAHAN MAGIC HOTA HAI! 🚀
	// Jaise hi order DB me gaya, kitchen ki saari screens par broadcast kar do
	h.broadcastToKitchen(order)

	// Order successful hone par wapas menu par bhej do
	c.Redirect(http.StatusFound, "/?success=1")
}

// --- Kitchen Handlers ---

func (h *MenuHandler) RenderKitchen(c *gin.Context) {
	orders, err := h.Usecase.FetchActiveOrders(c.Request.Context())
	if err != nil {
		orders = []domain.Order{}
	}
	c.HTML(http.StatusOK, "kitchen.html", gin.H{"orders": orders})
}

func (h *MenuHandler) CompleteOrder(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))
	h.Usecase.MarkOrderCompleted(c.Request.Context(), orderID)
	c.Redirect(http.StatusFound, "/kitchen")
}

// --- WebSocket Logic ---

// Kitchen screen WebSocket se connect hogi yahan se
func (h *MenuHandler) HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket error: %v", err)
		return
	}
	defer ws.Close()

	// Naye client (kitchen screen) ko map me add karo
	h.clientsMu.Lock()
	h.clients[ws] = true
	h.clientsMu.Unlock()

	// Infinite loop jab tak connection zinda hai
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// Agar kitchen tab close ho jaye toh client ko remove kar do
			h.clientsMu.Lock()
			delete(h.clients, ws)
			h.clientsMu.Unlock()
			break
		}
	}
}

// Har kitchen screen par JSON data bhejne wala function
func (h *MenuHandler) broadcastToKitchen(order *domain.Order) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	for client := range h.clients {
		err := client.WriteJSON(order)
		if err != nil {
			log.Printf("Kitchen screen ko bhejne me error: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
