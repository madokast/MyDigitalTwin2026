package mcp

import (
	"context"
	"dt2026/api/records"
	"dt2026/httpx"
	"dt2026/lib"
	"dt2026/server/http"
	"dt2026/server/mcp/api/posts/get"
	"dt2026/server/mcp/api/posts/post"
	"dt2026/server/mcp/api/posts/query"
	tagsget "dt2026/server/mcp/api/posts/tags/get"
	"dt2026/server/mcp/api/posts/tags/get_all"
	"dt2026/server/mcp/api/probe/health"
	"dt2026/server/mcp/api/probe/postgresql"
	"dt2026/server/mcp/api/probe/qqbot"
	mcptime "dt2026/server/mcp/api/time"

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
	s.addTool(McpTool[get.Input, get.Output]{
		Name:        get.Name,
		Description: get.Description,
		HandleFunc:  get.RecordsGet,
	})
	s.addTool(McpTool[post.Input, post.Output]{
		Name:        post.Name,
		Description: post.Description,
		HandleFunc:  post.RecordsPost,
	})
	s.addTool(McpTool[query.Input, query.Output]{
		Name:        query.Name,
		Description: query.Description,
		HandleFunc:  query.RecordsQuery,
	})
	s.addTool(McpTool[tagsget.Input, tagsget.Output]{
		Name:        tagsget.Name,
		Description: tagsget.Description,
		HandleFunc:  tagsget.RecordsTagsGet,
	})
	s.addTool(McpTool[get_all.Input, get_all.Output]{
		Name:        get_all.Name,
		Description: get_all.Description,
		HandleFunc:  get_all.RecordsTagsGetAll,
	})
}

func (s *Server) addTool[In McpInput, Out McpOutput](tool McpTool[In, Out]) {
	successSchema, err := jsonschema.For[Out](&jsonschema.ForOptions{
		TypeSchemas: lib.MapsMerge(
			records.JSONSchemaTypes,
			httpx.TrueTypeSchemaTypes,
		),
	})
	if err != nil {
		panic(err)
	}
	mcp_sdk.AddTool(
		s.mcpServer,
		&mcp_sdk.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			OutputSchema: &jsonschema.Schema{
				OneOf: []*jsonschema.Schema{successSchema, httpx.ErrorSchema()},
			},
		},
		func(
			ctx context.Context,
			req *mcp_sdk.CallToolRequest,
			input In,
		) (*mcp_sdk.CallToolResult, any, error) {
			out, httpErr := tool.HandleFunc(s.httpServer, input)
			if httpErr != nil {
				return &mcp_sdk.CallToolResult{IsError: true}, httpErr, nil
			}
			return nil, out, nil
		},
	)
}
