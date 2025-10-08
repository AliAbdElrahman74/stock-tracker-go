package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go-stock-tracker/tracker"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []tracker.StockData)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	mu sync.Mutex
)

func main() {
	router := gin.Default()
	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	router.GET("/ws", handleWebSocket)

	go broadcaster()
	go fetchLoop([]string{"NVDA"})

	fmt.Println("🚀 Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	mu.Lock()
	clients[conn] = true
	mu.Unlock()
}

func broadcaster() {
	for {
		data := <-broadcast
		message, _ := json.Marshal(data)
		mu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func fetchLoop(symbols []string) {
	for {
		data, err := tracker.FetchStockData(symbols)
		if err != nil {
			log.Println("Fetch error:", err)
		} else {
			broadcast <- data
		}
		time.Sleep(20 * time.Second)
	}
}
