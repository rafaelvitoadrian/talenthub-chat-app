package controllers

import (
	"chapter-d3/internal/service"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func SetupHandlers(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("../../web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/users", handleUsers)
	mux.HandleFunc("/ws", handleWebSocket)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("../../web/index.html")
	if err != nil {
		http.Error(w, "Could not open requested file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", content)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	numUsers := len(service.Connections)
	log.Printf("Ada berapa users: %v", numUsers)
	fmt.Fprintf(w, "%d", numUsers)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	currentConn := service.WebSocketConnection{Conn: currentGorillaConn, Username: username}
	service.Connections = append(service.Connections, &currentConn)
	// service.handleIO()
	go service.HandleIO(&currentConn, service.Connections)
}
