package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/novalagung/gubrak/v2"
)

type SocketPayload struct {
	Type    string
	Message string
	Image   string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
	Image   string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

type M map[string]interface{}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_IMAGE = "Image"
const MESSAGE_LEAVE = "Leave"

var Connections = make([]*WebSocketConnection, 0)

func HandleIO(currentConn *WebSocketConnection, Connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	BroadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				BroadcastMessage(currentConn, MESSAGE_LEAVE, "")
				EjectConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		switch payload.Type {
		case MESSAGE_CHAT:
			log.Printf("Ini basedata payload: %v", payload.Message)
			BroadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
		case MESSAGE_IMAGE:
			// log.Printf("Ini basedata payload: %v", payload.Image)
			BroadcastMessageImage(currentConn, MESSAGE_IMAGE, payload.Image)
		}
	}
}

func EjectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(Connections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()
	Connections = filtered.([]*WebSocketConnection)
}

func BroadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	// log.Printf("Ini basedata payload: %v", message)
	for _, eachConn := range Connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}

func BroadcastMessageImage(currentConn *WebSocketConnection, kind, image string) {
	// log.Printf("Ini basedata payload: %v", image)
	for _, eachConn := range Connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:  currentConn.Username,
			Type:  kind,
			Image: image,
		})
	}
}

func testApakahWork() {

}
