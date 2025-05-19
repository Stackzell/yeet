package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var logger = slog.Default().With("component", "messages")

type WebSocketEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
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
			// build url from a template
			t := template.Must(template.New("url").Parse(httpRequestMessage.URL))
			var buf bytes.Buffer
			err := t.Execute(&buf, httpRequestMessage.Variables)
			if err != nil {
				log.Println("Error executing template:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error executing template"))
				continue
			}

			// send the http request
			client := http.DefaultClient
			req, err := http.NewRequest(httpRequestMessage.Method, buf.String(), nil)
			if err != nil {
				log.Println("Error creating request:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error creating request"))
				continue
			}

			// set headers
			for key, value := range httpRequestMessage.Headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			//defer resp.Body.Close()
			if err != nil {
				log.Println("Error sending request:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error sending request"))
				continue
			}

			responseMessage := &HttpResponseMessage{
				Status: resp.StatusCode,
				Body:   "",
			}

			if contentType := resp.Header.Get("Content-Type"); contentType != "" {
				responseMessage.ContentType = contentType
			}

			if resp.Body != nil {
				body, err := io.ReadAll(resp.Body)
				responseMessage.Body = string(body)
				responseMessage.ContentType = resp.Header.Get("Content-Type")

				if err != nil {
					log.Println("Error reading response body:", err)
					conn.WriteMessage(websocket.TextMessage, []byte("Error reading response body"))
					continue
				}
			}

			responseData, _ := json.Marshal(responseMessage)
			conn.WriteJSON(&WebSocketEnvelope{
				Type: HttpResponseMessageType.String(),
				Data: responseData,
			})
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
