package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var logger = slog.Default().With("component", "messages")

type WebSocketEnvelope struct {
	CorrelationID string          `json:"correlationId"`
	Type          string          `json:"type"`
	Data          json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, or specify allowed origins here
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		var envelope WebSocketEnvelope
		if json.Unmarshal(p, &envelope) != nil {
			logger.Error("Error unmarshalling message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		message, err := UnmarshalMessage(MessageType(envelope.Type), envelope.Data)
		if err != nil {
			logger.Error("Error unmarshalling message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message type"))
			continue
		}

		switch message.EventType() {
		case HttpRequestMessageType:
			httpRequestMessage, _ := message.(*HttpRequestMessage)

			responseMessage, err := HandleHttpRequestMessage(httpRequestMessage)
			if err != nil {
				log.Println("Error handling request:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(
					fmt.Errorf("Error handling request: %w", err).Error(),
				))
				continue
			}

			data, err := json.Marshal(responseMessage)
			if err != nil {
				log.Println("Error marshalling response:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error marshalling response"))
				continue
			}

			envelope, err := json.Marshal(&WebSocketEnvelope{
				Type: HttpResponseMessageType.String(),
				Data: data,
			})

			conn.WriteMessage(websocket.TextMessage, envelope)
			break
		}
	}
}

func main() {
	r := chi.NewRouter()
	r.HandleFunc("/ws", handleWebSocket)

	// Define your routes here
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Start the server
	http.ListenAndServe(":8080", r)
}
