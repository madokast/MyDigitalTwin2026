package mcp

import (
	"dt2026/httpx"
	"dt2026/server/http"
	http_sdk "net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpInput = any

type McpOutput = httpx.Response

type Server struct {
	mcpServer  *mcp_sdk.Server
	httpServer *http.Server
}

func NewServer(httpServer *http.Server) *Server {
	s := &Server{
		mcpServer: mcp_sdk.NewServer(
			&mcp_sdk.Implementation{
				Name:        "MyDigitalTwin2026",
				Title:       "mdk's personal digital twin",
				Description: ServerDescription,
				Version:     "0.0.1",
			},
			&mcp_sdk.ServerOptions{
				Instructions: ServerInstructions,
			},
		),
		httpServer: httpServer,
	}
	s.addAllTools()
	return s
}

func (s *Server) HttpHandler(stateless bool) *mcp_sdk.StreamableHTTPHandler {
	opts := &mcp_sdk.StreamableHTTPOptions{
		Stateless:    stateless,
		JSONResponse: true,
	}
	if !stateless {
		opts.SessionTimeout = 30 * time.Second
	}
	return mcp.NewStreamableHTTPHandler(func(req *http_sdk.Request) *mcp.Server {
		return s.mcpServer
	}, opts)
}
