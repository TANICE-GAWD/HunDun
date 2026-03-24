//ticker >> generate task >> send to hub >> hub >> all clients

package hub

import (
	// "github.com/EliCDavis/vector/vector3"
	"github.com/google/uuid"
	// "github.com/orsinium-labs/enum"
	"math/rand"
	"time"
)



type Hub struct{
	register chan *Client
	unregister chan *Client
	clients map[*Client]bool //nested map
	broadcast chan Task
}



type CategoryEnum string

const (
	AI CategoryEnum = "ai"
	Fintech CategoryEnum = "fintech"
	Health CategoryEnum = "health"

)

var categories = []CategoryEnum{AI, Fintech, Health}

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}


type Task struct {
	ID          string       `json:"id"`
	Category    CategoryEnum `json:"category"`
	Coordinates Vector3      `json:"coordinates"`
}






func NewHub() *Hub{
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan Task),
	}
}

func Run(h *Hub){
	for{
		select{
		case client := <- h.register:
			h.RegisterClient(client)
		case client := <- h.unregister:
			h.UnRegisterClient(client)
		case task := <- h.broadcast:
			h.BroadcastTask(task)

		}
	}
}


func (h *Hub) RegisterClient(c *Client){
	h.clients[c] = true

}

func (h *Hub) UnRegisterClient(c *Client) {
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
}

func (h *Hub) BroadcastTask(task Task) {
	for client := range h.clients {
		select {
		case client.send <- task:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) StartTick() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			task := GenerateTask()
			h.broadcast <- task
		}
	}()
}



func GenerateTask() Task{
	return Task{
		ID: uuid.New().String(),
		Category : categories[rand.Intn(len(categories))],
		// Category: categories.Choice(0),
		Coordinates: Vector3{
			X: rand.Float64() * 100,
			Y: rand.Float64() * 100,
			Z: rand.Float64() * 100,
		},
	}
}