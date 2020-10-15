package BusinessModel

import "github.com/gorilla/websocket"

type Person struct {
	Id       int
	Email    string
	Password string
	GUID     string
}

type Connection struct {
	Id         int
	Connection *websocket.Conn
	RoomId     int
	GUID       string
}

type Room struct {
	Id          int
	NameRoom    string
	Created     string
	GUID        string
	Connections []Connection
}
