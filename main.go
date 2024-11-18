package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	serverListen()
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

	if method == "GET" {
		switch {
		case path == "/hello":
			handleHello(conn)
		case path == "/goodbye":
			handleGoodbye(conn)
		default:
			handleNotFound(conn)
		}
	} else if method == "POST" {
		contentType, ok := headers["Content-Type"]
		if ok && strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			handleFormDataPost(conn, body)
		} else {
			handleUnsupportedMediaType(conn)
		}
	} else {
		handleMethodNotAllowed(conn)
	}

}

func handleHello(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, World!"
	writeResponse(conn, response)
}

func handleGoodbye(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 8\r\n" +
		"\r\n" +
		"Goodbye!"
	writeResponse(conn, response)
}

func handleFormDataPost(conn net.Conn, body []byte) {
	formData, err := url.ParseQuery(string(body))
	if err != nil {
		log.Println("Failed to parse form data: ", err)
		handleBadRequest(conn, "Bad Request")
		return
	}

	fmt.Println("Parsed Form Data: ")
	for key, values := range formData {
		fmt.Printf("%s: %s\n", key, strings.Join(values, ", "))
	}

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: %d\r\n" +
		"\r\n" +
		"Form Data Received Sucessfully!"

	response = fmt.Sprintf(response, len("Form Data Receieved Sucessfully!"))
	writeResponse(conn, response)
}

func writeResponse(conn net.Conn, response string) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Failed to write response: ", err)
		return
	}
}

func handleNotFound(conn net.Conn) {
	response := "HTTP/1.1 404 Not Found\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 14\r\n" +
		"\r\n" +
		"404 Not Found"
	writeResponse(conn, response)
}

func handleMethodNotAllowed(conn net.Conn) {
	response := "HTTP/1.1 405 Method Not Allowed\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 24\r\n" +
		"\r\n" +
		"405 Method Not Allowed"
	writeResponse(conn, response)
}

func handleBadRequest(conn net.Conn, message string) {
	response := "HTTP/1.1 400 Bad Request\r\n" +
		"Content-Type: text/plain\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n", len(message)) +
		"\r\n" +
		message
	writeResponse(conn, response)
}

func handleInternalServerError(conn net.Conn) {
	response := "HTTP/1.1 500 Internal Server Error\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 21\r\n" +
		"\r\n" +
		"Internal Server Error"
	writeResponse(conn, response)
}

func handleUnsupportedMediaType(conn net.Conn) {
	response := "HTTP/1.1 415 Unsupported Media Type\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 25\r\n" +
		"\r\n" +
		"415 Unsupported Media Type"
	writeResponse(conn, response)
}
