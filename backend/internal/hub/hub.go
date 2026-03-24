//ticker >> generate task >> send to hub >> hub >> all clients

package chat

import (
	"fmt"
	// "github.com/EliCDavis/vector/vector3"
	"github.com/google/uuid"
	// "github.com/orsinium-labs/enum"
	"math/rand"
)



type Hub struct{
	register chan *Client
	unregister chan *Client
	clients map[string]map[*Client]bool //nested map
	broadcast chan *Task
}

// type Message struct{
// 	Sender string `json:"sender"`
// 	Receiver string `json:"receiver"`
// 	Content string `json:"content"`
// 	ID string `json:"id"`
// 	Type string `json:"type"`
// }

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

ticker := time.Newticker(500 * time.Millisecond)
defer ticker.Stop()
for range ticker.C {
	task := GenerateTask()
	h.broadcast <- task
	
}


func GenerateTask() (t *Task){
	return &Task{
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




func NewHub() *Hub{
	return &Hub{
		clients: make(map[string]map[*Client]bool)
		register: make(chan *Client)
		unregister: make(chan *Client)
		broadcast: make(chan Task)
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
			h.TaskHandler(task)

		}
	}
}


func (h *Hub) RegisterClient(c *Client){
	conn := h.clients[c.ID]
	if conn == nil{
		conn := make(map[*Client]bool)
		h.clients[c.ID] = conn
	}
	h.clients[c.ID][c] = true

}

func (h *Hub) UnRegisterClient(c *Client){
	if _,ok:=h.clients[c.ID]; ok{
		delete(h.clients[c.ID], c)
		close(client.send)
	}

}

func (h *Hub) MessageHandler(message *Message){
	// if message.Type == "message"{
	// 	clients := h.clients[message.ID]
	// 	for client := range clients{
	// 		select{
	// 		case client.send <- message:
	// 		default:
	// 			close(client.send)
	// 			delete(h.clients[message.ID],client)
	// 		}
	// 	}
	// }
	// else if message.Type == "notification"{
	// 	clients:= h.clients[message.Receiver]
	// 	for client := range clients{
	// 		select{
	// 		case client.send <- message:
	// 		default:
	// 			close(client.send)
	// 			delete(h.clients[message.Receiver],client)
	// 		}
	// 	}
	// }
	
	clients := h.clients
	for i := range clients{
		select{
		case h.clients[i].send <- message:
		default:
			close(i.send)
			delete(h.clients[message.ID],h.clients[i])
		}
	}


}