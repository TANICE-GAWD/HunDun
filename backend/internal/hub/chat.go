//ticker >> generate task >> send to hub >> hub >> all clients

package hub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)




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

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.now)

	
	
	for{
		var msg Message
		c.hub.broadcast <- msg
	}



}




func (c *Client) Close(){
	close(c.send)
}

