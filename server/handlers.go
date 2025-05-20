package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

func HandleHttpRequestMessage(m *HttpRequestMessage) (*HttpResponseMessage, error) {
	// build url from a template
	t := template.Must(template.New("url").Parse(m.URL))
	var buf bytes.Buffer
	err := t.Execute(&buf, m.Variables)
	if err != nil {
		return nil, err
	}

	// send the http request
	client := http.DefaultClient
	req, err := http.NewRequest(m.Method, buf.String(), nil)
	if err != nil {
		return nil, err
	}

	// set headers
	for key, value := range m.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

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
			return nil, err
		}
	}

	return responseMessage, nil
}

func HandleTemplateRenderMessage(m *RenderTemplateRequestMessage) (*RenderTemplateResponseMessage, error) {
	tmpl, err := template.New(m.Template).Parse(m.Template)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, m.Variables)
	if err != nil {
		return nil, err
	}

	return &RenderTemplateResponseMessage{
		Render: buf.String(),
	}, nil
}
