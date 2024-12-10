# metoro-mcp-server
This repository contains the source code of the Metoro MCP (Model Context Protocol) Server. This MCP Server allows you to interact with your Kubernetes cluster via the Claude Desktop App (or soon Metoro MCP Client!).



## What is MCP (Model Context Protocol)? 
Definition on the MCP website: https://modelcontextprotocol.io
> The Model Context Protocol (MCP) is an open protocol that enables seamless integration between LLM applications and external data sources and tools. Whether you’re building an AI-powered IDE, enhancing a chat interface, or creating custom AI workflows, MCP provides a standardized way to connect LLMs with the context they need.


## What is Metoro?
[Metoro](https://metoro.io/) is an observability tool designed for microservices running in Kubernetes and uses eBPF based instrumentation to autogenerate telemetry for you.

## How can I use Metoro MCP Server? 

### If you already have a Metoro Account:
// TODO: 

### If you don't have a Metoro Account:
No worries, you can still play around using the [Live Demo Cluster](https://demo.us-east.metoro.io/). 

1. Install the [Claude Desktop App](https://claude.ai/download).
2. Install go on your machine:
   * Mac: `brew install go` 
   * Linux: `sudo apt-get install golang-go`
3. Clone this repository and build the Metoro MCP Server:
```bash
go build 
```
4. Create a file in `~/Library/Application Support/Claude/claude_desktop_config.json` with the following contents:
```json
{
  "mcpServers": {
    "metoro-mcp-server": {
      "command": "<your path to Metoro MCP server go executable>/metoro-mcp-server",
      "args": [],
      "env": {
          "METORO_AUTH_TOKEN" : "<demo metoro account auth token>",
          "METORO_API_URL": "https://demo.us-east.metoro.io"
       }
    }
  }
}
```
