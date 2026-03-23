//ticker >> generate task >> send to hub >> hub >> all clients

package hub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


type Point struct{
	X float64
	Y float64
	Z float64
}
type Task struct{
	ID string
	Category


}

var upgrader = websocket.Upgrader{
	ReadBufferSize : 1024,
	WriteBufferSize : 1024,
	
}


type Client struct{
	ID string
	Conn *websocket.Conn
	send chan Message
	hub *Hub
}




func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID:id,Conn:conn, send: make(chan Message)}
}

func (c *Client) Read(){
	defer func(){
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	

	
	
	for{
		var msg Message
		c.hub.broadcast <- msg
	}



}




func (c *Client) Close(){
	close(c.send)
}

