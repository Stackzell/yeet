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
	CorrelationID string          `json:"correlationId,omitempty"`
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
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		var envelope WebSocketEnvelope
		if json.Unmarshal(p, &envelope) != nil {
			logger.Error("Error unmarshalling message:", err)
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		message, err := UnmarshalMessage(MessageType(envelope.Type), envelope.Data)
		if err != nil {
			logger.Error("Error unmarshalling message:", err)
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Invalid message type"))
			continue
		}

		writer := NewResponseWriter(envelope.CorrelationID, conn)
		go func() {
			switch message.EventType() {
			case HttpRequestMessageType:
				httpRequestMessage, _ := message.(*HttpRequestMessage)

				responseMessage, err := HandleHttpRequestMessage(httpRequestMessage)
				if err != nil {
					log.Println("Error handling request:", err)
					writer.WriteErrorResponse(
						fmt.Errorf("error handling request: %w", err),
					)
					return
				}

				writer.WriteJsonResponse(HttpResponseMessageType, responseMessage)
				return
			case RenderTemplateRequestMessageType:
				renderTemplateRequestMessage, _ := message.(*RenderTemplateRequestMessage)
				responseMessage, err := HandleTemplateRenderMessage(renderTemplateRequestMessage)
				if err != nil {
					log.Println("Error handling template render request:", err)
					writer.WriteErrorResponse(
						fmt.Errorf("error handling template render request: %w", err),
					)
					return
				}

				writer.WriteJsonResponse(RenderTemplateResponseMessageType, responseMessage)
				return
			}
		}()
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
