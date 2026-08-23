package mcp

import (
	"context"
	"dt2026/httpx"
	"dt2026/server/http"
	"dt2026/server/mcp/api/probe/health"
	mcptime "dt2026/server/mcp/api/time"

	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpTool[In McpInput, Out McpOutput] struct {
	Name        string
	Description string
	HandleFunc  func(*http.Server, In) (Out, *httpx.Error)
}

func (s *Server) addAllTools() {
	s.addTool(McpTool[health.Input, health.Output]{
		Name:        health.Name,
		Description: health.Description,
		HandleFunc:  health.ProbeHealth,
	})
	s.addTool(McpTool[mcptime.Input, mcptime.Output]{
		Name:        mcptime.Name,
		Description: mcptime.Description,
		HandleFunc:  mcptime.Time,
	})
}

func (s *Server) addTool[In McpInput, Out McpOutput](tool McpTool[In, Out]) {
	mcp_sdk.AddTool(
		s.mcpServer,
		&mcp_sdk.Tool{
			Name:        tool.Name,
			Description: tool.Description,
		},
		func(
			ctx context.Context,
			req *mcp_sdk.CallToolRequest,
			input In,
		) (*mcp_sdk.CallToolResult, Out, error) {
			out, httpErr := tool.HandleFunc(s.httpServer, input)
			if httpErr != nil {
				var zero Out
				return nil, zero, httpErr
			}
			return nil, out, nil
		},
	)
}
