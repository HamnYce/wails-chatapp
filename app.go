package main

import (
	g "chatapp/global"
	"context"
	"io"
	"log"
	"net"
)

// App struct
type App struct {
	ctx      context.Context
	conn     net.Conn
	messages []g.Message
	buffer   []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.conn = nil
	a.messages = make([]g.Message, 0)
	a.buffer = make([]byte, 1024)
}

func (a *App) ConnectToChatroom() {
	if a.conn != nil {
		return
	}

  conn, err := net.Dial("tcp4", g.HOST+":"+g.PORT)
	if err != nil {
		conn = nil
		return
	}
	a.conn = conn

	log.Println("Successfully connected to server")
}

func (a *App) DisconnectFromChatroom() {
	if !a.IsConnected() {
		return
	}
	a.conn.Close()
	a.conn = nil
	log.Println("disconnected from server")
}

func (a *App) SendMsgToChatRoom(msg string) {
	if !a.IsConnected() {
		return
	}
	a.conn.Write([]byte(msg))
}

func (a *App) ListenForMessage() {
	if !a.IsConnected() {
		return
	}
	n, err := a.conn.Read(a.buffer)

	// _, err := a.conn.Read(a.buffer)
	if err == io.EOF {
		a.DisconnectFromChatroom()
		return
	} else if err != nil {
		a.DisconnectFromChatroom()
		return
	}
	log.Println("Received From Server:", string(a.buffer[:n]))
	a.messages = append(a.messages, g.NewMessageFromByteSlice(a.buffer[:n]))
}

func (a *App) GetMessages() []g.Message {
	return a.messages
}

func (a *App) IsConnected() bool {
	return a.conn != nil
}
