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

	if method == "GET" {
		switch {
		case path == "/hello":
			handleHello(conn)
		case path == "/goodbye":
			handleGoodbye(conn)
		default:
			handleNotFound(conn)
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
		"405 Not Found"
	writeResponse(conn, response)
}

func writeResponse(conn net.Conn, response string) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("Failed to write response: ", err)
		return
	}
}
