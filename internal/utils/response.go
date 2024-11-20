package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/bmolinari/go-http-server/pkg/httpstatus"
)

func WriteResponse(conn net.Conn, statusCode, contentType, message string) {
	contentLength := len(message)
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, contentType, contentLength, message)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Failed to write response: ", err)
		return
	}
}

func WriteJSONResponse(conn net.Conn, statusCode, message string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		WriteResponse(conn, httpstatus.InternalServerError, "application/json", `{"error": "Internal Server Error"}`)
		return
	}

	contentLength := len(jsonData)
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", statusCode, contentLength, jsonData)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Failed to write JSON response: ", err)
	}

}

func WriteErrorResponse(conn net.Conn, statusCode string, errorMessage string) {
	WriteResponse(conn, statusCode, "application/json", fmt.Sprintf(`{"error": "%s"}`, errorMessage))
}

func WriteRedirectResponse(conn net.Conn, location string) {
	response := fmt.Sprintf("HTTP/1.1 302 Found\r\nLocation: %s\r\n\r\n", location)
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Failed to write redirect response: ", err)
	}
}

func HandleNotFound(conn net.Conn, message string) {
	WriteResponse(conn, httpstatus.NotFound, "text/plain", message)
}

func HandleMethodNotAllowed(conn net.Conn, message string) {
	WriteResponse(conn, httpstatus.NotFound, "text/plain", message)
}

func HandleBadRequest(conn net.Conn, message string) {
	WriteResponse(conn, httpstatus.BadRequest, "text/plain", message)
}

func HandleInternalServerError(conn net.Conn, message string) {
	WriteResponse(conn, httpstatus.InternalServerError, "text/plain", message)
}

func HandleUnsupportedMediaType(conn net.Conn, message string) {
	WriteResponse(conn, httpstatus.UnsupportedMediaType, "text/plain", message)
}

func ReadHeadersAndBody(reader *bufio.Reader) (map[string]string, []byte) {
	headers := make(map[string]string)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(headerLine) == "" {
			break
		}
		parts := strings.SplitN(headerLine, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	var body []byte
	if cl, ok := headers["Content-Length"]; ok {
		var contentLength int
		fmt.Sscanf(cl, "%d", &contentLength)
		body = make([]byte, contentLength)
		_, _ = reader.Read(body)
	}
	return headers, body
}
