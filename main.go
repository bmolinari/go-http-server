package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/bmolinari/go-http-server/httpstatus"
)

var router = map[string]map[string]func(net.Conn, map[string]string, []byte){}

func main() {
	registerRoute("GET", "/hello", handleHello)
	registerRoute("GET", "/goodbye", handleGoodbye)
	registerRoute("POST", "/form-data", handleFormDataPost)
	registerRoute("POST", "/json", handleJsonPost)

	serverListen()
}

func registerRoute(method, path string, handler func(net.Conn, map[string]string, []byte)) {
	if router[method] == nil {
		router[method] = make(map[string]func(net.Conn, map[string]string, []byte))
	}
	router[method][path] = handler
}

func serverListen() {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {

		log.Println("Failed to read request line: ", err)
		return
	}
	fmt.Println("Request Line: ", requestLine)

	var method, path, httpVersion string
	_, err = fmt.Sscanf(requestLine, "%s %s %s", &method, &path, &httpVersion)
	if err != nil {
		log.Println("Failed to parse request line: ", err)
		return
	}

	parsedUrl, err := url.Parse(path)
	if err != nil {
		log.Println("Failed to parse URL: ", err)
		return
	}
	path = parsedUrl.Path
	fmt.Println("Path: ", path)

	headers := make(map[string]string)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading header: ", err)
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

	fmt.Println("-Headers-")
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
	}

	var body []byte
	if cl, ok := headers["Content-Length"]; ok {
		var contentLength int
		fmt.Sscanf(cl, "%d", &contentLength)
		body = make([]byte, contentLength)
		_, err := reader.Read(body)
		if err != nil {
			log.Println("Failed to read body: ", err)
		}
	}
	fmt.Println("Body: ", string(body))

	if routes, ok := router[method]; ok {
		if handler, ok := routes[path]; ok {
			handler(conn, headers, body)
			return
		}
	}
	handleNotFound(conn)
}

func handleHello(conn net.Conn, headers map[string]string, body []byte) {
	writeResponse(conn, httpstatus.OK, "text/plain", "Hello, World!")
}

func handleGoodbye(conn net.Conn, headers map[string]string, body []byte) {
	writeResponse(conn, httpstatus.OK, "text/plain", "Goodbye!")
}

func handleFormDataPost(conn net.Conn, headers map[string]string, body []byte) {
	formData, err := url.ParseQuery(string(body))
	if err != nil {
		log.Println("Failed to parse form data: ", err)
		handleBadRequest(conn)
		return
	}

	fmt.Println("-Form Data-")
	for key, values := range formData {
		fmt.Printf("%s: %s\n", key, strings.Join(values, ", "))
	}

	writeResponse(conn, httpstatus.OK, "text/plain", "Form Data Received Sucessfully!")
}

func handleJsonPost(conn net.Conn, headers map[string]string, body []byte) {
	fmt.Println("-JSON Body")
	fmt.Println(string(body))
	writeResponse(conn, httpstatus.OK, "text/plain", "JSON Data Receieved Sucessfully!")
}

func handleNotFound(conn net.Conn) {
	writeResponse(conn, httpstatus.NotFound, "text/plain", "404 Not Found")
}

func handleMethodNotAllowed(conn net.Conn) {
	writeResponse(conn, httpstatus.NotFound, "text/plain", "405 Method Not Allowed")
}

func handleBadRequest(conn net.Conn) {
	writeResponse(conn, httpstatus.BadRequest, "text/plain", "400 Bad Request")
}

func handleInternalServerError(conn net.Conn) {
	writeResponse(conn, httpstatus.InternalServerError, "text/plain", "Internal Server Error")
}

func handleUnsupportedMediaType(conn net.Conn) {
	writeResponse(conn, httpstatus.UnsupportedMediaType, "text/plain", "415 Unsupported Media Type")
}

func writeResponse(conn net.Conn, statusCode, contentType, message string) {
	contentLength := len(message)
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, contentType, contentLength, message)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Failed to write response: ", err)
		return
	}
}
