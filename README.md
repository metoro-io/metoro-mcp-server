<div align="center">
<img src="./resources/Metoro_square.svg" height="300" alt="Metoro MCP Logo">
</div>
<br/>
<div align="center">

![GitHub stars](https://img.shields.io/github/stars/metoro-io/metoro-mcp-server?style=social)
![GitHub forks](https://img.shields.io/github/forks/metoro-io/metoro-mcp-server?style=social)
![GitHub issues](https://img.shields.io/github/issues/metoro-io/metoro-mcp-server)
![GitHub pull requests](https://img.shields.io/github/issues-pr/metoro-io/metoro-mcp-server)
![GitHub license](https://img.shields.io/github/license/metoro-io/metoro-mcp-server)
![GitHub contributors](https://img.shields.io/github/contributors/metoro-io/metoro-mcp-server)
![GitHub last commit](https://img.shields.io/github/last-commit/metoro-io/metoro-mcp-server)
[![GoDoc](https://pkg.go.dev/badge/github.com/metoro-io/metoro-mcp-server.svg)](https://pkg.go.dev/github.com/metoro-io/metoro-mcp-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/metoro-io/metoro-mcp-server)](https://goreportcard.com/report/github.com/metoro-io/metoro-mcp-server)
![Tests](https://github.com/metoro-io/metoro-mcp-server/actions/workflows/go-test.yml/badge.svg)

</div>

# metoro-mcp-server
This repository contains th Metoro MCP (Model Context Protocol) Server. This MCP Server allows you to interact with your Kubernetes cluster via the Claude Desktop App!

## What is MCP (Model Context Protocol)? 
You can read more about the Model Context Protocol here: https://modelcontextprotocol.io

But in a nutshell
> The Model Context Protocol (MCP) is an open protocol that enables seamless integration between LLM applications and external data sources and tools. Whether you’re building an AI-powered IDE, enhancing a chat interface, or creating custom AI workflows, MCP provides a standardized way to connect LLMs with the context they need.

## What is Metoro?
[Metoro](https://metoro.io/) is an observability platform designed for microservices running in Kubernetes and uses eBPF based instrumentation to generate deep telemetry without code changes.
The data that is generated by the eBPF agents is sent to Metoro's backend to be stored and in the Metoro frontend using our apis.

This MCP server exposes those APIs to an LLM so you can ask your AI questions about your Kubernetes cluster.

## How can I use Metoro MCP Server? 
1. Install the [Claude Desktop App](https://claude.ai/download).
2. Download the Metoro MCP Server from the latest release: https://github.com/metoro-io/metoro-mcp-server/releases

### If you already have a Metoro Account:
Copy your auth token from your Metoro account in [Settings](https://us-east.metoro.io/settings) -> Users Settings. 
Create a file in `~/Library/Application Support/Claude/claude_desktop_config.json` with the following contents:
```json
{
  "mcpServers": {
    "metoro-mcp-server": {
      "command": "<your path to Metoro MCP server go executable>/metoro-mcp-server",
      "args": [],
      "env": {
          "METORO_AUTH_TOKEN" : "<your auth token>",
          "METORO_API_URL": "https://us-east.metoro.io"
       }
    }
  }
}
```

### If you don't have a Metoro Account:
No worries, you can still play around using the [Live Demo Cluster](https://demo.us-east.metoro.io/).
The included token is a demo token, publicly available for anyone to use.
   Create a file in `~/Library/Application Support/Claude/claude_desktop_config.json` with the following contents:
```json
{
  "mcpServers": {
    "metoro-mcp-server": {
      "command": "<your path to Metoro MCP server go executable>/metoro-mcp-server",
      "args": [],
      "env": {
          "METORO_AUTH_TOKEN" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcklkIjoiOThlZDU1M2QtYzY4ZC00MDRhLWFhZjItNDM2ODllNWJiMGUzIiwiZW1haWwiOiJ0ZXN0QGNocmlzYmF0dGFyYmVlLmNvbSIsImV4cCI6MTgyMTI0NzIzN30.7G6alDpcZh_OThYj293Jce5rjeOBqAhOlANR_Fl5auw",
          "METORO_API_URL": "https://demo.us-east.metoro.io"
       }
    }
  }
}
```

4. Once you are done editing `claude_desktop_config.json` save the file and restart Claude Desktop app.
5. You should now see the Metoro MCP Server in the dropdown list of MCP Servers in the Claude Desktop App. You are ready to start using Metoro MCP Server with Claude Desktop App!


