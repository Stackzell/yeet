package main

import "encoding/json"

type JsonWriter interface {
	WriteJSON(interface{}) error
}

type ResponseWriter struct {
	correlationID string
	writer        JsonWriter
}

func NewResponseWriter(correlationID string, writer JsonWriter) *ResponseWriter {
	return &ResponseWriter{
		correlationID: correlationID,
		writer:        writer,
	}
}

func (w *ResponseWriter) WriteErrorResponse(err error) {
	w.WriteJsonResponse(ErrorMessageType, &ErrorMessage{
		Message: err.Error(),
	})
}

func (w *ResponseWriter) WriteJsonResponse(messageType MessageType, data interface{}) {
	dataBytes, _ := json.Marshal(data)
	reply := &WebSocketEnvelope{
		CorrelationID: w.correlationID,
		Type:          messageType.String(),
		Data:          dataBytes,
	}
	if err := w.writer.WriteJSON(reply); err != nil {
		logger.Error("Error writing message", "error", err)
	}
}
