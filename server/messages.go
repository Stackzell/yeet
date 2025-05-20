package main

import (
	"encoding/json"
	"fmt"
)

type MessageType string

func (t MessageType) String() string {
	return string(t)
}

const (
	HttpRequestMessageType            MessageType = "http.request"
	HttpResponseMessageType           MessageType = "http.response"
	RenderTemplateRequestMessageType  MessageType = "render.template.request"
	RenderTemplateResponseMessageType MessageType = "render.template.response"
	ErrorMessageType                  MessageType = "error"
)

var messageRegistry = map[MessageType]func() Message{
	HttpRequestMessageType:            func() Message { return &HttpRequestMessage{} },
	HttpResponseMessageType:           func() Message { return &HttpResponseMessage{} },
	RenderTemplateRequestMessageType:  func() Message { return &RenderTemplateRequestMessage{} },
	RenderTemplateResponseMessageType: func() Message { return &RenderTemplateResponseMessage{} },
	ErrorMessageType:                  func() Message { return &ErrorMessage{} },
}

func UnmarshalMessage(messageType MessageType, data []byte) (Message, error) {
	constructor, ok := messageRegistry[messageType]
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", messageType)
	}

	msg := constructor()
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

type Message interface {
	EventType() MessageType
}

// HttpRequestMessage is a message type for sending HTTP requests.
type HttpRequestMessage struct {
	Method    string                 `json:"method"`
	URL       string                 `json:"url"`
	Headers   map[string]string      `json:"headers"`
	Body      string                 `json:"body,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (m *HttpRequestMessage) EventType() MessageType {
	return HttpRequestMessageType
}

// HttpResponseMessage is a message type for sending HTTP responses.
type HttpResponseMessage struct {
	Body        string `json:"body,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Status      int    `json:"status"`
}

func (m *HttpResponseMessage) EventType() MessageType {
	return HttpResponseMessageType
}

// RenderTemplateRequestMessage is a message type for rendering templates.
type RenderTemplateRequestMessage struct {
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (m *RenderTemplateRequestMessage) EventType() MessageType {
	return RenderTemplateRequestMessageType
}

// RenderTemplateResponseMessage is a message type for rendering template responses.
type RenderTemplateResponseMessage struct {
	Render string `json:"render"`
}

func (m *RenderTemplateResponseMessage) EventType() MessageType {
	return RenderTemplateResponseMessageType
}

// ErrorMessage is a message type for sending error messages.
type ErrorMessage struct {
	Message string `json:"message"`
}

func (m *ErrorMessage) EventType() MessageType {
	return ErrorMessageType
}
