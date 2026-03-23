//ticker >> generate task >> send to hub >> hub >> all clients

package chat

import (
	"fmt"
)



type Hub struct{
	register chan *Client
	unregister chan *Client
	clients map[string]map[*Client]bool //nested map
	broadcast chan Message
}

type Message struct{
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Content string `json:"content"`
	ID string `json:"id"`
	Type string `json:"type"`
}

func NewHub() *Hub{
	return &Hub{
		clients: make(map[string]map[*Client]bool)
		register: make(chan *Client)
		unregister: make(chan *Client)
		broadcast: make(chan Message)
	}
}

func Run(h *Hub){
	for{
		select{
		case client := <- h.register:
			h.RegisterClient(client)
		case client := <- h.unregister:
			h.UnRegisterClient(client)
		case message := <- h.message:
			h.MessageHandler(message)

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
	if message.Type == "message"{
		clients := h.clients[message.ID]
		for client := range clients{
			select{
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.ID],client)
			}
		}
	}
	else if message.Type == "notification"{
		clients:= h.clients[message.Receiver]
		for client := range clients{
			select{
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.Receiver],client)
			}
		}
	}
	


}