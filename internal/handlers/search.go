package handlers

import (
	"fmt"
	"net"
	"net/url"

	"github.com/bmolinari/go-http-server/internal/utils"
	"github.com/bmolinari/go-http-server/pkg/httpstatus"
)

func HandleSearch(conn net.Conn, headers map[string]string, queryParams url.Values, body []byte) {
	requiredKeys := []string{"name", "city"}
	params, err := utils.ValidateQueryParams(queryParams, requiredKeys)
	if err != nil {
		utils.WriteResponse(conn, httpstatus.BadRequest, "text/plain", err.Error())
		return
	}
	result := fmt.Sprintf("Searched for %s from %s", params["name"], params["city"])
	utils.WriteResponse(conn, httpstatus.OK, "text/plain", result)
}
