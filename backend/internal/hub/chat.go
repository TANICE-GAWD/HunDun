//ticker >> generate task >> send to hub >> hub >> all clients

package hub

import (
	"time"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)




var upgrader = websocket.Upgrader{
	ReadBufferSize : 1024,
	WriteBufferSize : 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
	
}


type Client struct{
	ID string
	Conn *websocket.Conn
	send chan Task
	hub *Hub
}




func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID: id, Conn: conn, send: make(chan Task), hub: hub}
}


func (c *Client) Read(){
	defer func(){
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	
	
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			break
		}
	}


}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for task := range c.send {
		if err := c.Conn.WriteJSON(task); err != nil {
			break
		}
	}
}


func (c *Client) Close(){
	close(c.send)
}

func ServeWs(h *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(uuid.New().String(), conn, h)

	h.register <- client

	go client.Read()
	go client.Write()
}