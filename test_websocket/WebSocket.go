package WebSocket

import (
    "flag"
	"html/template"
    "net/http"
    "log"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrade:", err)
        return
    }
    // ... Use conn to send and receive messages.
    echo(conn)
}

// p is a []byte and messageType is an int with value websocket.BinaryMessage or websocket.TextMessage.

func echo(conn *Conn) {
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }    
}

/*
An application can also send and receive messages using the io.WriteCloser and io.Reader interfaces. 
To send a message, call the connection NextWriter method to get an io.WriteCloser, write the message 
to the writer and close the writer when done. 
To receive a message, call the connection NextReader method to get an io.Reader and read until io.EOF is returned. 
This snippet shows how to echo messages using the NextWriter and NextReader methods:
*/
func next() {
    for {
        messageType, r, err := conn.NextReader()
        if err != nil {
            return
        }
        w, err := conn.NextWriter(messageType)
        if err != nil {
            return err
        }
        if _, err := io.Copy(w, r); err != nil {
            return err
        }
        if err := w.Close(); err != nil {
            return err
        }
    }
}

/*
Control Messages
The WebSocket protocol defines three types of control messages: close, ping and pong. 
Call the connection WriteControl, WriteMessage or NextWriter methods to send a control message to the peer.
*/
/*
The application must read the connection to process close, ping and pong messages sent from the peer. 
If the application is not otherwise interested in messages from the peer, then the application should 
start a goroutine to read and discard messages from the peer. 
A simple example is:
*/
func readLoop(c *websocket.Conn) {
    for {
        if _, _, err := c.NextReader(); err != nil {
            c.Close()
            break
        }
    }
}
