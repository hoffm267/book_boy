package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEClient struct {
	Channel chan []byte
}

type SSEManager struct {
	clients map[*SSEClient]bool
	mu      sync.RWMutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[*SSEClient]bool),
	}
}

func (m *SSEManager) AddClient(client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client] = true
}

func (m *SSEManager) RemoveClient(client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, client)
	close(client.Channel)
}

func (m *SSEManager) Broadcast(eventType string, data interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal SSE data: %v\n", err)
		return
	}

	message := []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, jsonData))

	fmt.Printf("Broadcasting SSE event '%s' to %d clients: %s\n", eventType, len(m.clients), string(jsonData))

	for client := range m.clients {
		select {
		case client.Channel <- message:
			fmt.Printf("Sent to client\n")
		default:
			fmt.Printf("Failed to send to client (channel full)\n")
		}
	}
}

func (m *SSEManager) ServeHTTP(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	client := &SSEClient{
		Channel: make(chan []byte, 10),
	}

	m.AddClient(client)
	defer m.RemoveClient(client)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return false
			}
			w.Write(msg)
			return true
		case <-ticker.C:
			w.Write([]byte(":keepalive\n\n"))
			return true
		}
	})
}
