package handlers

import (
	"net"
	"net/url"

	"github.com/bmolinari/go-http-server/internal/utils"
	"github.com/bmolinari/go-http-server/pkg/httpstatus"
)

func HandleGoodbye(conn net.Conn, headers map[string]string, queryParams url.Values, body []byte) {
	utils.WriteResponse(conn, httpstatus.OK, "text/plain", "Goodbye!")
}
