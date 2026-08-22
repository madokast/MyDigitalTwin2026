package mcp

import (
	"dt2026/httpx"
	"dt2026/server/http"
	http_sdk "net/http"

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
				Title:       "mdk 个人的数字孪生项目",
				Description: ServerDescription,
				Version:     "0.0.1",
			},
			nil,
		),
		httpServer: httpServer,
	}
	s.addAllTools()
	return s
}

func (s *Server) HttpHandler() *mcp_sdk.StreamableHTTPHandler {
	return mcp.NewStreamableHTTPHandler(func(req *http_sdk.Request) *mcp.Server {
		return s.mcpServer
	}, &mcp_sdk.StreamableHTTPOptions{
		Stateless:    true,
		JSONResponse: true,
	})
}
