package mcp

import (
	"context"
	"dt2026/api/records"
	"dt2026/httpx"
	"dt2026/server/http"
	"dt2026/server/mcp/api/posts/post"
	"dt2026/server/mcp/api/probe/health"
	"dt2026/server/mcp/api/probe/postgresql"
	"dt2026/server/mcp/api/probe/qqbot"
	mcptime "dt2026/server/mcp/api/time"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
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
	s.addTool(McpTool[postgresql.Input, postgresql.Output]{
		Name:        postgresql.Name,
		Description: postgresql.Description,
		HandleFunc:  postgresql.ProbePostgreSQL,
	})
	s.addTool(McpTool[qqbot.Input, qqbot.Output]{
		Name:        qqbot.Name,
		Description: qqbot.Description,
		HandleFunc:  qqbot.ProbeQQBot,
	})
	s.addTool(McpTool[post.Input, post.Output]{
		Name:        post.Name,
		Description: post.Description,
		HandleFunc:  post.RecordsPost,
	})
}

func (s *Server) addTool[In McpInput, Out McpOutput](tool McpTool[In, Out]) {
	outSchema, err := jsonschema.For[Out](&jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeFor[records.JSONTime](): {Type: "string"},
		},
	})
	if err != nil {
		panic(err)
	}
	mcp_sdk.AddTool(
		s.mcpServer,
		&mcp_sdk.Tool{
			Name:         tool.Name,
			Description:  tool.Description,
			OutputSchema: outSchema,
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
