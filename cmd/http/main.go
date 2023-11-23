package main

import (
	// b64 "encoding/base64"

	"chapter-d3/internal/controllers"
	"fmt"
	"net/http"
)

// var service.Connections = make([]*service.WebSocketConnection, 0)

func main() {
	// Serve static files (CSS and JS) from the "static" directory
	mux := http.NewServeMux()
	controllers.SetupHandlers(mux)

	port := 8080

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server is running on http://0.0.0.0%s\n", addr)

	fmt.Println("Server starting at :8080")
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(err)
	}
}
