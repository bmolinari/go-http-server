package routes

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/bmolinari/go-http-server/internal/handlers"
	"github.com/bmolinari/go-http-server/internal/utils"
)

type Router struct {
	routes map[string]map[string]func(net.Conn, map[string]string, url.Values, []byte)
}

func NewRouter() *Router {
	r := &Router{routes: make(map[string]map[string]func(net.Conn, map[string]string, url.Values, []byte))}
	r.Register("GET", "/hello", handlers.HandleHello)
	r.Register("GET", "/goodbye", handlers.HandleGoodbye)
	r.Register("GET", "/search", handlers.HandleSearch)
	r.Register("POST", "/json", handlers.HandleJsonPost)
	return r
}

func (r *Router) Register(method, path string, handler func(net.Conn, map[string]string, url.Values, []byte)) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]func(net.Conn, map[string]string, url.Values, []byte))
	}
	r.routes[method][path] = handler
}

func (r *Router) HandleConnection(conn net.Conn) {
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
	queryParams := parsedUrl.Query()
	headers, body := utils.ReadHeadersAndBody(reader)

	if routes, ok := r.routes[method]; ok {
		if handler, ok := routes[path]; ok {
			handler(conn, headers, queryParams, body)
			return
		}
	}
	utils.HandleNotFound(conn, "404 Not Found")
}
