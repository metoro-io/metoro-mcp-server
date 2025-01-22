package main

import (
	"fmt"
	"os"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/metoro-io/metoro-mcp-server/resources"
	"github.com/metoro-io/metoro-mcp-server/tools"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func main() {
	// Check if the appropriate environment variables are set
	if err := checkEnvVars(); err != nil {
		panic(err)
	}

	done := make(chan struct{})

	mcpServer := mcpgolang.NewServer(stdio.NewStdioServerTransport())

	// Add tools
	for _, tool := range tools.MetoroToolsList {
		err := mcpServer.RegisterTool(tool.Name, tool.Description, tool.Handler)
		if err != nil {
			panic(err)
		}
	}

	// Add resources
	for _, resource := range resources.MetoroResourcesList {
		err := mcpServer.RegisterResource(
			resource.Path,
			resource.Name,
			resource.Description,
			resource.ContentType,
			resource.Handler)
		if err != nil {
			panic(err)
		}
	}

	err := mcpServer.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}

func checkEnvVars() error {
	if os.Getenv(utils.METORO_API_URL_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", utils.METORO_API_URL_ENV_VAR)
	}
	if os.Getenv(utils.METORO_AUTH_TOKEN_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", utils.METORO_AUTH_TOKEN_ENV_VAR)
	}
	return nil
}
