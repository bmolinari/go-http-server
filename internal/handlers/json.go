package handlers

import (
	"fmt"
	"net"
	"net/url"

	"github.com/bmolinari/go-http-server/internal/utils"
	"github.com/bmolinari/go-http-server/pkg/httpstatus"
)

func HandleJsonPost(conn net.Conn, headers map[string]string, queryParams url.Values, body []byte) {
	fmt.Println("-JSON Body")
	fmt.Println(string(body))
	utils.WriteResponse(conn, httpstatus.OK, "text/plain", "JSON Data Receieved Sucessfully!")
}
